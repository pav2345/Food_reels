package services

import (
	"context"

	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
)

type FoodPartnerService struct {
	foodPartnerRepo *repository.FoodPartnerRepository
	foodRepo        *repository.FoodRepository
}

func NewFoodPartnerService(
	foodPartnerRepo *repository.FoodPartnerRepository,
	foodRepo *repository.FoodRepository,
) *FoodPartnerService {
	return &FoodPartnerService{
		foodPartnerRepo: foodPartnerRepo,
		foodRepo:        foodRepo,
	}
}

func (s *FoodPartnerService) GetFoodPartnerByID(ctx context.Context, id uuid.UUID) (*models.FoodPartner, []models.Food, error) {
	partner, err := s.foodPartnerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	foods, err := s.foodRepo.FindByPartnerID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return partner, foods, nil
}
