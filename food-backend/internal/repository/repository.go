package repository

import (
	"context"
	"errors"

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
