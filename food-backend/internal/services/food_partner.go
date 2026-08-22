package services

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
)

type FoodPartnerService struct {
	foodPartnerRepo *repository.FoodPartnerRepository
	foodRepo        *repository.FoodRepository
	storageSvc      *StorageService
	geoService      *GeoService
}

func NewFoodPartnerService(
	foodPartnerRepo *repository.FoodPartnerRepository,
	foodRepo *repository.FoodRepository,
	storageSvc *StorageService,
	geoService *GeoService,
) *FoodPartnerService {
	return &FoodPartnerService{
		foodPartnerRepo: foodPartnerRepo,
		foodRepo:        foodRepo,
		storageSvc:      storageSvc,
		geoService:      geoService,
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

func (s *FoodPartnerService) UpdateDeliveryRadius(ctx context.Context, partnerID uuid.UUID, radiusKM float64) (*models.FoodPartner, error) {
	if radiusKM <= 0 {
		return nil, fmt.Errorf("delivery radius must be greater than zero")
	}

	if err := s.foodPartnerRepo.UpdateProfile(ctx, partnerID, map[string]interface{}{
		"delivery_radius_km": radiusKM,
	}); err != nil {
		return nil, err
	}

	return s.foodPartnerRepo.FindByID(ctx, partnerID)
}

func (s *FoodPartnerService) UpdateWorkingHours(ctx context.Context, partnerID uuid.UUID, openingTime, closingTime string) (*models.FoodPartner, error) {
	openingTime = strings.TrimSpace(openingTime)
	closingTime = strings.TrimSpace(closingTime)
	if _, err := parseBusinessClock(nowUTC(), openingTime); err != nil {
		return nil, fmt.Errorf("invalid opening time")
	}
	if _, err := parseBusinessClock(nowUTC(), closingTime); err != nil {
		return nil, fmt.Errorf("invalid closing time")
	}

	if err := s.foodPartnerRepo.UpdateProfile(ctx, partnerID, map[string]interface{}{
		"opening_time": openingTime,
		"closing_time": closingTime,
	}); err != nil {
		return nil, err
	}

	return s.foodPartnerRepo.FindByID(ctx, partnerID)
}

func (s *FoodPartnerService) UpdateProfile(ctx context.Context, partnerID uuid.UUID, req dto.UpdateFoodPartnerProfileRequest) (*models.FoodPartner, error) {
	updates := make(map[string]interface{})
	setStringUpdate(updates, "name", req.Name)
	setStringUpdate(updates, "contact_name", req.ContactName)
	setStringUpdate(updates, "phone", req.Phone)
	setStringUpdate(updates, "address", req.Address)
	setStringUpdate(updates, "profile_image", req.ProfileImage)
	setStringUpdate(updates, "cover_image", req.CoverImage)

	if req.Latitude != nil {
		if err := s.geoService.ValidateCoordinate(*req.Latitude, longitudeOrZero(req.Longitude)); err != nil {
			return nil, err
		}
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		if err := s.geoService.ValidateCoordinate(latitudeOrZero(req.Latitude), *req.Longitude); err != nil {
			return nil, err
		}
		updates["longitude"] = *req.Longitude
	}

	if err := s.foodPartnerRepo.UpdateProfile(ctx, partnerID, updates); err != nil {
		return nil, err
	}

	return s.foodPartnerRepo.FindByID(ctx, partnerID)
}

func (s *FoodPartnerService) UploadFoodImage(ctx context.Context, partner models.FoodPartner, file []byte, fileName, mimeType string) (string, error) {
	url, err := s.uploadPartnerAsset(ctx, file, fileName, mimeType)
	if err != nil {
		return "", err
	}

	images := appendStringList(partner.FoodImages, url)
	if err := s.foodPartnerRepo.UpdateProfile(ctx, partner.ID, map[string]interface{}{"food_images": images}); err != nil {
		return "", err
	}

	return url, nil
}

func (s *FoodPartnerService) UploadMenuImage(ctx context.Context, partner models.FoodPartner, file []byte, fileName, mimeType string) (string, error) {
	url, err := s.uploadPartnerAsset(ctx, file, fileName, mimeType)
	if err != nil {
		return "", err
	}

	images := appendStringList(partner.MenuImages, url)
	if err := s.foodPartnerRepo.UpdateProfile(ctx, partner.ID, map[string]interface{}{"menu_images": images}); err != nil {
		return "", err
	}

	return url, nil
}

func (s *FoodPartnerService) UploadMenuPDF(ctx context.Context, partnerID uuid.UUID, file []byte, fileName, mimeType string) (string, error) {
	url, err := s.uploadPartnerAsset(ctx, file, fileName, mimeType)
	if err != nil {
		return "", err
	}

	if err := s.foodPartnerRepo.UpdateProfile(ctx, partnerID, map[string]interface{}{"menu_pdf": url}); err != nil {
		return "", err
	}

	return url, nil
}

func (s *FoodPartnerService) uploadPartnerAsset(ctx context.Context, file []byte, fileName, mimeType string) (string, error) {
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".bin"
	}

	uploadResult, err := s.storageSvc.UploadFileWithMime(ctx, file, uuid.New().String()+ext, mimeType)
	if err != nil {
		return "", err
	}

	return uploadResult.URL, nil
}

func setStringUpdate(updates map[string]interface{}, column string, value *string) {
	if value != nil {
		updates[column] = strings.TrimSpace(*value)
	}
}

func appendStringList(raw, value string) string {
	var values []string
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &values)
	}
	values = append(values, value)

	encoded, err := json.Marshal(values)
	if err != nil {
		return raw
	}
	return string(encoded)
}

func latitudeOrZero(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func longitudeOrZero(value *float64) float64 {
	if value == nil {
		return 0
	}
	return *value
}

func nowUTC() time.Time {
	return time.Now().UTC()
}
