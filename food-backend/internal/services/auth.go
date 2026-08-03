package services

import (
	"context"
	"errors"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/repository"
	"food-backend/internal/utils"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	foodPartnerRepo  *repository.FoodPartnerRepository
	jwtSecret        string
}

func NewAuthService(
	userRepo *repository.UserRepository,
	foodPartnerRepo *repository.FoodPartnerRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		foodPartnerRepo: foodPartnerRepo,
		jwtSecret:       jwtSecret,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dto.RegisterUserRequest) (*models.User, string, error) {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", ErrUserAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := utils.SignToken(s.jwtSecret, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) LoginUser(ctx context.Context, req dto.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if user == nil || !utils.ComparePassword(user.Password, req.Password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.SignToken(s.jwtSecret, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) RegisterFoodPartner(ctx context.Context, req dto.RegisterFoodPartnerRequest) (*models.FoodPartner, string, error) {
	existingPartner, err := s.foodPartnerRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if existingPartner != nil {
		return nil, "", ErrFoodPartnerAlreadyExists
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	partner := &models.FoodPartner{
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashedPassword,
		Phone:       req.Phone,
		Address:     req.Address,
		ContactName: req.ContactName,
	}

	if err := s.foodPartnerRepo.Create(ctx, partner); err != nil {
		return nil, "", err
	}

	token, err := utils.SignToken(s.jwtSecret, partner.ID)
	if err != nil {
		return nil, "", err
	}

	return partner, token, nil
}

func (s *AuthService) LoginFoodPartner(ctx context.Context, req dto.LoginRequest) (*models.FoodPartner, string, error) {
	partner, err := s.foodPartnerRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", err
	}
	if partner == nil || !utils.ComparePassword(partner.Password, req.Password) {
		return nil, "", ErrInvalidCredentials
	}

	token, err := utils.SignToken(s.jwtSecret, partner.ID)
	if err != nil {
		return nil, "", err
	}

	return partner, token, nil
}

var (
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrFoodPartnerAlreadyExists = errors.New("food partner already exists")
	ErrInvalidCredentials       = errors.New("invalid email or password")
)
