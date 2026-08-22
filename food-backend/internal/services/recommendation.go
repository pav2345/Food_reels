package services

import (
	"math"
	"sort"
	"time"

	"food-backend/internal/models"
)

type FeedCandidate struct {
	Food       models.Food
	DistanceKM float64
	IsOpen     bool
	CanOrder   bool
}

type RecommendationContext struct {
	SavedFoodIDs        map[string]bool
	LikedFoodIDs        map[string]bool
	PreferredCategories map[string]bool
	Now                 time.Time
}

type RecommendationService struct{}

func NewRecommendationService() *RecommendationService {
	return &RecommendationService{}
}

func (s *RecommendationService) Rank(candidates []FeedCandidate, ctx RecommendationContext) []FeedCandidate {
	if ctx.Now.IsZero() {
		ctx.Now = time.Now().UTC()
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return s.score(candidates[i], ctx) > s.score(candidates[j], ctx)
	})

	return candidates
}

func (s *RecommendationService) score(candidate FeedCandidate, ctx RecommendationContext) float64 {
	partner := candidate.Food.FoodPartner
	score := 0.0

	score += math.Max(0, DefaultFeedRadiusKM-candidate.DistanceKM) * 8
	if candidate.IsOpen {
		score += 25
	}
	if candidate.CanOrder {
		score += 10
	}
	score += partner.Rating * 6
	score += float64(candidate.Food.Likes)*0.7 + float64(candidate.Food.Saves)
	score += recencyScore(candidate.Food.CreatedAt, ctx.Now)

	foodID := candidate.Food.ID.String()
	if ctx.LikedFoodIDs[foodID] {
		score += 4
	}
	if ctx.SavedFoodIDs[foodID] {
		score += 6
	}
	if category := candidate.Food.Category; category != nil && ctx.PreferredCategories[*category] {
		score += 8
	}

	return score
}

func recencyScore(createdAt, now time.Time) float64 {
	if createdAt.IsZero() {
		return 0
	}

	ageHours := now.Sub(createdAt).Hours()
	if ageHours <= 0 {
		return 10
	}
	if ageHours >= 168 {
		return 0
	}
	return 10 * (1 - ageHours/168)
}
