package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"time"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
)

const (
	defaultFeedLimit = 20
	maxFeedLimit     = 50
)

type FeedService struct {
	foodPartnerRepo   *repository.FoodPartnerRepository
	foodRepo          *repository.FoodRepository
	saveRepo          *repository.SaveRepository
	geoService        *GeoService
	openService       *RestaurantOpenService
	deliveryService   *DeliveryService
	recommendationSvc *RecommendationService
	localRadiusKM     float64
	localFeedPercent  int
}

func NewFeedService(
	foodPartnerRepo *repository.FoodPartnerRepository,
	foodRepo *repository.FoodRepository,
	saveRepo *repository.SaveRepository,
	geoService *GeoService,
	openService *RestaurantOpenService,
	deliveryService *DeliveryService,
	recommendationSvc *RecommendationService,
	localRadiusKM float64,
	localFeedPercent int,
) *FeedService {
	if localRadiusKM <= 0 {
		localRadiusKM = DefaultFeedRadiusKM
	}

	if localFeedPercent <= 0 || localFeedPercent > 100 {
		localFeedPercent = 60
	}

	return &FeedService{
		foodPartnerRepo:   foodPartnerRepo,
		foodRepo:          foodRepo,
		saveRepo:          saveRepo,
		geoService:        geoService,
		openService:       openService,
		deliveryService:   deliveryService,
		recommendationSvc: recommendationSvc,
		localRadiusKM:     localRadiusKM,
		localFeedPercent:  localFeedPercent,
	}
}

func (s *FeedService) Generate(
	ctx context.Context,
	userID uuid.UUID,
	userLocation Coordinate,
	cursor string,
	limit int,
) (dto.FeedResponse, error) {
	if limit <= 0 {
		limit = defaultFeedLimit
	}

	if limit > maxFeedLimit {
		limit = maxFeedLimit
	}

	if err := s.geoService.ValidateCoordinate(
		userLocation.Latitude,
		userLocation.Longitude,
	); err != nil {
		return dto.FeedResponse{}, err
	}

	offset, err := decodeFeedCursor(cursor)
	if err != nil {
		return dto.FeedResponse{}, err
	}

	partners, err := s.foodPartnerRepo.FindAllLocated(ctx)
	if err != nil {
		return dto.FeedResponse{}, err
	}

	if len(partners) == 0 {
		return dto.FeedResponse{
			Items:   []dto.FeedItemResponse{},
			HasMore: false,
		}, nil
	}

	partnerDistances := make(map[uuid.UUID]float64, len(partners))
	allPartnerIDs := make([]uuid.UUID, 0, len(partners))

	for _, partner := range partners {
		partnerCoordinate, ok := PartnerCoordinate(partner)
		if !ok {
			continue
		}

		distance := s.geoService.DistanceKM(
			userLocation,
			*partnerCoordinate,
		)

		partnerDistances[partner.ID] = round1(distance)
		allPartnerIDs = append(allPartnerIDs, partner.ID)
	}

	if len(allPartnerIDs) == 0 {
		return dto.FeedResponse{
			Items:   []dto.FeedItemResponse{},
			HasMore: false,
		}, nil
	}

	foods, err := s.foodRepo.FindAvailableByPartnerIDsWithPartner(
		ctx,
		allPartnerIDs,
	)
	if err != nil {
		return dto.FeedResponse{}, err
	}

	if len(foods) == 0 {
		return dto.FeedResponse{
			Items:   []dto.FeedItemResponse{},
			HasMore: false,
		}, nil
	}

	savedFoodIDs, preferredCategories, err :=
		s.userFoodSignals(ctx, userID)
	if err != nil {
		return dto.FeedResponse{}, err
	}

	likedUUIDs, err := s.foodRepo.FindLikedFoodIDs(ctx, userID)
	if err != nil {
		return dto.FeedResponse{}, err
	}

	likedFoodIDs := make(map[string]bool, len(likedUUIDs))

	for foodID := range likedUUIDs {
		likedFoodIDs[foodID.String()] = true
	}

	recommendationContext := RecommendationContext{
		SavedFoodIDs:        savedFoodIDs,
		LikedFoodIDs:        likedFoodIDs,
		PreferredCategories: preferredCategories,
		Now:                 time.Now().UTC(),
	}

	localCandidates := make([]FeedCandidate, 0, len(foods))
	discoveryCandidates := make([]FeedCandidate, 0, len(foods))

	for _, food := range foods {
		distance, ok := partnerDistances[food.FoodPartnerID]
		if !ok {
			continue
		}

		partner := food.FoodPartner

		isOpen := s.openService.IsOpen(partner)
		canOrder := isOpen && distance <= partner.DeliveryRadiusKM

		candidate := FeedCandidate{
			Food:       food,
			DistanceKM: distance,
			IsOpen:     isOpen,
			CanOrder:   canOrder,
		}

		if distance <= s.localRadiusKM {
			localCandidates = append(localCandidates, candidate)
		} else {
			discoveryCandidates = append(discoveryCandidates, candidate)
		}
	}

	rankedLocal := s.recommendationSvc.Rank(
		localCandidates,
		recommendationContext,
	)

	rankedDiscovery := s.recommendationSvc.Rank(
		discoveryCandidates,
		recommendationContext,
	)

	mixed := s.mixFeedCandidates(
		rankedLocal,
		rankedDiscovery,
		limit,
	)

	if len(mixed) == 0 {
		return dto.FeedResponse{
			Items:   []dto.FeedItemResponse{},
			HasMore: false,
		}, nil
	}

	if offset > len(mixed) {
		offset = len(mixed)
	}

	end := offset + limit

	if end > len(mixed) {
		end = len(mixed)
	}

	items := make([]dto.FeedItemResponse, 0, end-offset)

	for _, candidate := range mixed[offset:end] {
		items = append(
			items,
			s.newFeedItemResponse(candidate),
		)
	}

	hasMore := end < len(mixed)

	nextCursor := ""

	if hasMore {
		nextCursor = encodeFeedCursor(end)
	}

	return dto.FeedResponse{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *FeedService) mixFeedCandidates(
	localCandidates []FeedCandidate,
	discoveryCandidates []FeedCandidate,
	limit int,
) []FeedCandidate {
	if limit <= 0 {
		limit = defaultFeedLimit
	}

	if len(localCandidates) == 0 {
		if len(discoveryCandidates) > limit {
			return discoveryCandidates[:limit]
		}

		return discoveryCandidates
	}

	if len(discoveryCandidates) == 0 {
		if len(localCandidates) > limit {
			return localCandidates[:limit]
		}

		return localCandidates
	}

	selected := make(
		[]FeedCandidate,
		0,
		min(limit, len(localCandidates)+len(discoveryCandidates)),
	)

	localIndex := 0
	discoveryIndex := 0

	pattern := []bool{
		true,
		true,
		true,
		false,
		false,
	}

	patternIndex := 0

	for len(selected) < limit {
		preferLocal := pattern[patternIndex]
		patternIndex++

		if patternIndex >= len(pattern) {
			patternIndex = 0
		}

		if preferLocal && localIndex < len(localCandidates) {
			selected = append(
				selected,
				localCandidates[localIndex],
			)
			localIndex++
			continue
		}

		if !preferLocal && discoveryIndex < len(discoveryCandidates) {
			selected = append(
				selected,
				discoveryCandidates[discoveryIndex],
			)
			discoveryIndex++
			continue
		}

		if localIndex < len(localCandidates) {
			selected = append(
				selected,
				localCandidates[localIndex],
			)
			localIndex++
			continue
		}

		if discoveryIndex < len(discoveryCandidates) {
			selected = append(
				selected,
				discoveryCandidates[discoveryIndex],
			)
			discoveryIndex++
			continue
		}

		break
	}

	return selected
}

func (s *FeedService) userFoodSignals(
	ctx context.Context,
	userID uuid.UUID,
) (map[string]bool, map[string]bool, error) {
	saves, err := s.saveRepo.FindByUserWithFoodAndPartner(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	savedFoodIDs := make(map[string]bool, len(saves))
	preferredCategories := make(map[string]bool)

	for _, save := range saves {
		savedFoodIDs[save.FoodID.String()] = true

		if save.Food.Category != nil &&
			*save.Food.Category != "" {
			preferredCategories[*save.Food.Category] = true
		}
	}

	return savedFoodIDs, preferredCategories, nil
}

func (s *FeedService) newFeedItemResponse(
	candidate FeedCandidate,
) dto.FeedItemResponse {
	food := candidate.Food
	partner := food.FoodPartner
	distance := round1(candidate.DistanceKM)

	return dto.FeedItemResponse{
		Reel: dto.NewFoodResponse(
			&food,
			false,
		),

		Restaurant: dto.FoodPartnerSummary{
			ID:               partner.ID.String(),
			Name:             partner.Name,
			Logo:             partner.ProfileImage,
			Rating:           partner.Rating,
			DeliveryRadiusKM: partner.DeliveryRadiusKM,
		},

		DistanceKM: distance,

		EstimatedDeliveryMinutes: s.deliveryService.EstimateDeliveryMinutes(distance),

		IsOpen: candidate.IsOpen,

		OpenStatus: s.openService.Status(partner),

		CanOrder: candidate.CanOrder,
	}
}

func encodeFeedCursor(offset int) string {
	return base64.RawURLEncoding.EncodeToString(
		[]byte(strconv.Itoa(offset)),
	)
}

func decodeFeedCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}

	decoded, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("invalid cursor")
	}

	offset, err := strconv.Atoi(string(decoded))
	if err != nil || offset < 0 {
		return 0, fmt.Errorf("invalid cursor")
	}

	return offset, nil
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func PartnerCoordinate(
	partner models.FoodPartner,
) (*Coordinate, bool) {
	if partner.Latitude == nil ||
		partner.Longitude == nil {
		return nil, false
	}

	return &Coordinate{
		Latitude:  *partner.Latitude,
		Longitude: *partner.Longitude,
	}, true
}
