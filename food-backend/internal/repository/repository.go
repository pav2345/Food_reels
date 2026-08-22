package repository

import (
	"context"
	"errors"
	"time"

	"food-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) UpdateLocation(ctx context.Context, userID uuid.UUID, latitude, longitude float64, updatedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"latitude":                    latitude,
			"longitude":                   longitude,
			"current_latitude":            latitude,
			"current_longitude":           longitude,
			"last_location_updated":       updatedAt,
			"current_location_updated_at": updatedAt,
			"updated_at":                  updatedAt,
		}).Error
}

func (r *UserRepository) SetSavedLocation(ctx context.Context, userID uuid.UUID, latitude, longitude float64, updatedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"saved_latitude":  latitude,
			"saved_longitude": longitude,
			"updated_at":      updatedAt,
		}).Error
}

func (r *UserRepository) InitializeLocation(ctx context.Context, userID uuid.UUID, latitude, longitude float64, updatedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"latitude":                    latitude,
			"longitude":                   longitude,
			"saved_latitude":              latitude,
			"saved_longitude":             longitude,
			"current_latitude":            latitude,
			"current_longitude":           longitude,
			"last_location_updated":       updatedAt,
			"current_location_updated_at": updatedAt,
			"updated_at":                  updatedAt,
		}).Error
}

func (r *UserRepository) UpdateCurrentLocation(ctx context.Context, userID uuid.UUID, latitude, longitude float64, updatedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"latitude":                    latitude,
			"longitude":                   longitude,
			"current_latitude":            latitude,
			"current_longitude":           longitude,
			"last_location_updated":       updatedAt,
			"current_location_updated_at": updatedAt,
			"updated_at":                  updatedAt,
		}).Error
}

type FoodPartnerRepository struct {
	db *gorm.DB
}

func NewFoodPartnerRepository(db *gorm.DB) *FoodPartnerRepository {
	return &FoodPartnerRepository{db: db}
}

func (r *FoodPartnerRepository) FindByEmail(ctx context.Context, email string) (*models.FoodPartner, error) {
	var partner models.FoodPartner
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&partner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

func (r *FoodPartnerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.FoodPartner, error) {
	var partner models.FoodPartner
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&partner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &partner, nil
}

func (r *FoodPartnerRepository) Create(ctx context.Context, partner *models.FoodPartner) error {
	return r.db.WithContext(ctx).Create(partner).Error
}

func (r *FoodPartnerRepository) FindAllLocated(ctx context.Context) ([]models.FoodPartner, error) {
	var partners []models.FoodPartner
	err := r.db.WithContext(ctx).
		Where("latitude IS NOT NULL AND longitude IS NOT NULL").
		Find(&partners).Error
	return partners, err
}

func (r *FoodPartnerRepository) UpdateProfile(ctx context.Context, partnerID uuid.UUID, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Model(&models.FoodPartner{}).
		Where("id = ?", partnerID).
		Updates(updates).Error
}

type FoodRepository struct {
	db *gorm.DB
}

func NewFoodRepository(db *gorm.DB) *FoodRepository {
	return &FoodRepository{db: db}
}

func (r *FoodRepository) Create(ctx context.Context, food *models.Food) error {
	return r.db.WithContext(ctx).Create(food).Error
}

func (r *FoodRepository) FindAllWithPartner(ctx context.Context) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.WithContext(ctx).
		Preload("FoodPartner").
		Find(&foods).Error
	return foods, err
}

func (r *FoodRepository) FindByPartnerID(ctx context.Context, partnerID uuid.UUID) ([]models.Food, error) {
	var foods []models.Food
	err := r.db.WithContext(ctx).
		Where("food_partner_id = ?", partnerID).
		Find(&foods).Error
	return foods, err
}

func (r *FoodRepository) FindAvailableByPartnerIDsWithPartner(ctx context.Context, partnerIDs []uuid.UUID) ([]models.Food, error) {
	if len(partnerIDs) == 0 {
		return []models.Food{}, nil
	}

	var foods []models.Food
	err := r.db.WithContext(ctx).
		Preload("FoodPartner").
		Where("available = ? AND food_partner_id IN ?", true, partnerIDs).
		Find(&foods).Error
	return foods, err
}

func (r *FoodRepository) FindLikedFoodIDs(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	var likes []models.Like
	err := r.db.WithContext(ctx).
		Select("food_id").
		Where("user_id = ?", userID).
		Find(&likes).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]bool, len(likes))
	for _, like := range likes {
		result[like.FoodID] = true
	}
	return result, nil
}

func (r *FoodRepository) IncrementLikes(ctx context.Context, foodID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Food{}).
		Where("id = ?", foodID).
		UpdateColumn("likes", gorm.Expr("likes + ?", delta)).Error
}

func (r *FoodRepository) IncrementSaves(ctx context.Context, foodID uuid.UUID, delta int) error {
	return r.db.WithContext(ctx).
		Model(&models.Food{}).
		Where("id = ?", foodID).
		UpdateColumn("saves", gorm.Expr("saves + ?", delta)).Error
}

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) FindByUserAndFood(ctx context.Context, userID, foodID uuid.UUID) (*models.Like, error) {
	var like models.Like
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND food_id = ?", userID, foodID).
		First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *LikeRepository) Create(ctx context.Context, like *models.Like) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *LikeRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Like{}, "id = ?", id).Error
}

type SaveRepository struct {
	db *gorm.DB
}

func NewSaveRepository(db *gorm.DB) *SaveRepository {
	return &SaveRepository{db: db}
}

func (r *SaveRepository) FindByUserAndFood(ctx context.Context, userID, foodID uuid.UUID) (*models.Save, error) {
	var save models.Save
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND food_id = ?", userID, foodID).
		First(&save).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &save, nil
}

func (r *SaveRepository) FindByUserWithFoodAndPartner(ctx context.Context, userID uuid.UUID) ([]models.Save, error) {
	var saves []models.Save
	err := r.db.WithContext(ctx).
		Preload("Food.FoodPartner").
		Where("user_id = ?", userID).
		Find(&saves).Error
	return saves, err
}

func (r *SaveRepository) Create(ctx context.Context, save *models.Save) error {
	return r.db.WithContext(ctx).Create(save).Error
}

func (r *SaveRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Save{}, "id = ?", id).Error
}
