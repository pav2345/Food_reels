package services

import (
	"context"
	"errors"

	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrOutsideDeliveryRadius = errors.New("user is outside restaurant delivery radius")

type DeliveryService struct {
	geoService        *GeoService
	userRepo          *repository.UserRepository
	foodPartnerRepo   *repository.FoodPartnerRepository
	restaurantOpenSvc *RestaurantOpenService
}

func NewDeliveryService(
	geoService *GeoService,
	userRepo *repository.UserRepository,
	foodPartnerRepo *repository.FoodPartnerRepository,
	restaurantOpenSvc *RestaurantOpenService,
) *DeliveryService {
	return &DeliveryService{
		geoService:        geoService,
		userRepo:          userRepo,
		foodPartnerRepo:   foodPartnerRepo,
		restaurantOpenSvc: restaurantOpenSvc,
	}
}

func (s *DeliveryService) CanOrder(userLocation Coordinate, partner models.FoodPartner) bool {
	if !s.restaurantOpenSvc.IsOpen(partner) {
		return false
	}

	distance := s.DistanceToRestaurant(userLocation, partner)
	return distance >= 0 && distance <= partner.DeliveryRadiusKM
}

func (s *DeliveryService) DistanceToRestaurant(userLocation Coordinate, partner models.FoodPartner) float64 {
	if partner.Latitude == nil || partner.Longitude == nil {
		return -1
	}

	return s.geoService.DistanceKM(userLocation, Coordinate{
		Latitude:  *partner.Latitude,
		Longitude: *partner.Longitude,
	})
}

func (s *DeliveryService) EstimateDeliveryMinutes(distanceKM float64) string {
	switch {
	case distanceKM <= 2:
		return "20-25 mins"
	case distanceKM <= 4:
		return "25-35 mins"
	default:
		return "35-45 mins"
	}
}

func (s *DeliveryService) ValidateRestaurantDelivery(ctx context.Context, userID, restaurantID uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || user.Latitude == nil || user.Longitude == nil {
		return ErrLocationUnavailable
	}

	partner, err := s.foodPartnerRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return err
	}
	if partner == nil || partner.Latitude == nil || partner.Longitude == nil {
		return ErrLocationUnavailable
	}

	userLocation := Coordinate{Latitude: *user.Latitude, Longitude: *user.Longitude}
	if s.DistanceToRestaurant(userLocation, *partner) > partner.DeliveryRadiusKM {
		return ErrOutsideDeliveryRadius
	}

	return nil
}
