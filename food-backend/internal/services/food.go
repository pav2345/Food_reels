package services

import (
	"context"
	"errors"

	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoodService struct {
	db         *gorm.DB
	foodRepo   *repository.FoodRepository
	likeRepo   *repository.LikeRepository
	saveRepo   *repository.SaveRepository
	storageSvc *StorageService
}

func NewFoodService(
	db *gorm.DB,
	foodRepo *repository.FoodRepository,
	likeRepo *repository.LikeRepository,
	saveRepo *repository.SaveRepository,
	storageSvc *StorageService,
) *FoodService {
	return &FoodService{
		db:         db,
		foodRepo:   foodRepo,
		likeRepo:   likeRepo,
		saveRepo:   saveRepo,
		storageSvc: storageSvc,
	}
}

func (s *FoodService) CreateFood(ctx context.Context, partnerID uuid.UUID, name, description string, video []byte) (*models.Food, error) {
	fileName := uuid.New().String() + ".mp4"
	uploadResult, err := s.storageSvc.UploadFile(ctx, video, fileName)
	if err != nil {
		return nil, err
	}

	var desc *string
	if description != "" {
		desc = &description
	}

	food := &models.Food{
		Name:          name,
		Description:   desc,
		VideoURL:      uploadResult.URL,
		FoodPartnerID: partnerID,
		Likes:         0,
		Saves:         0,
		Version:       0,
	}

	if err := s.foodRepo.Create(ctx, food); err != nil {
		return nil, err
	}

	return food, nil
}

func (s *FoodService) GetFoodItems(ctx context.Context) ([]models.Food, error) {
	return s.foodRepo.FindAllWithPartner(ctx)
}

func (s *FoodService) ToggleLike(ctx context.Context, userID, foodID uuid.UUID) (bool, error) {
	var unliked bool

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		likeRepo := repository.NewLikeRepository(tx)
		foodRepo := repository.NewFoodRepository(tx)

		existingLike, err := likeRepo.FindByUserAndFood(ctx, userID, foodID)
		if err != nil {
			return err
		}

		if existingLike != nil {
			if err := likeRepo.DeleteByID(ctx, existingLike.ID); err != nil {
				return err
			}
			if err := foodRepo.IncrementLikes(ctx, foodID, -1); err != nil {
				return err
			}
			unliked = true
			return nil
		}

		like := &models.Like{
			UserID: userID,
			FoodID: foodID,
		}
		if err := likeRepo.Create(ctx, like); err != nil {
			return err
		}
		if err := foodRepo.IncrementLikes(ctx, foodID, 1); err != nil {
			return err
		}
		unliked = false
		return nil
	})

	return unliked, err
}

func (s *FoodService) ToggleSave(ctx context.Context, userID, foodID uuid.UUID) (bool, error) {
	var unsaved bool

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		saveRepo := repository.NewSaveRepository(tx)
		foodRepo := repository.NewFoodRepository(tx)

		existingSave, err := saveRepo.FindByUserAndFood(ctx, userID, foodID)
		if err != nil {
			return err
		}

		if existingSave != nil {
			if err := saveRepo.DeleteByID(ctx, existingSave.ID); err != nil {
				return err
			}
			if err := foodRepo.IncrementSaves(ctx, foodID, -1); err != nil {
				return err
			}
			unsaved = true
			return nil
		}

		save := &models.Save{
			UserID: userID,
			FoodID: foodID,
		}
		if err := saveRepo.Create(ctx, save); err != nil {
			return err
		}
		if err := foodRepo.IncrementSaves(ctx, foodID, 1); err != nil {
			return err
		}
		unsaved = false
		return nil
	})

	return unsaved, err
}

func (s *FoodService) GetSavedFoods(ctx context.Context, userID uuid.UUID) ([]models.Food, error) {
	saves, err := s.saveRepo.FindByUserWithFoodAndPartner(ctx, userID)
	if err != nil {
		return nil, err
	}

	foods := make([]models.Food, 0, len(saves))
	for _, save := range saves {
		foods = append(foods, save.Food)
	}

	return foods, nil
}

func (s *FoodService) GetFoodPartnerStats(ctx context.Context, partnerID uuid.UUID) (totalReels, totalLikes, totalSaves int, err error) {
	foods, err := s.foodRepo.FindByPartnerID(ctx, partnerID)
	if err != nil {
		return 0, 0, 0, err
	}

	totalReels = len(foods)
	for _, food := range foods {
		totalLikes += food.Likes
		totalSaves += food.Saves
	}

	return totalReels, totalLikes, totalSaves, nil
}

var ErrVideoRequired = errors.New("video file is required")
