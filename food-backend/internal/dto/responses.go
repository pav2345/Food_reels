package dto

import (
	"encoding/json"

	"food-backend/internal/models"
	"food-backend/internal/utils"
)

// ============================================================
// USER RESPONSE
// ============================================================

type UserResponse struct {
	ID                  string   `json:"id"`
	FullName            string   `json:"fullName"`
	Latitude            *float64 `json:"latitude,omitempty"`
	Longitude           *float64 `json:"longitude,omitempty"`
	LastLocationUpdated string   `json:"lastLocationUpdated,omitempty"`
}

func NewUserResponse(user *models.User) UserResponse {
	locationUpdated := ""

	if user.LastLocationUpdated != nil {
		locationUpdated = utils.FormatMongoDate(*user.LastLocationUpdated)
	}

	return UserResponse{
		ID:                  user.ID.String(),
		FullName:            user.FullName,
		Latitude:            user.Latitude,
		Longitude:           user.Longitude,
		LastLocationUpdated: locationUpdated,
	}
}

// ============================================================
// FOOD PARTNER AUTH RESPONSE
// ============================================================

type FoodPartnerAuthResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewFoodPartnerAuthResponse(
	partner *models.FoodPartner,
) FoodPartnerAuthResponse {
	return FoodPartnerAuthResponse{
		ID:   partner.ID.String(),
		Name: partner.Name,
	}
}

// ============================================================
// FOOD PARTNER SUMMARY
// ============================================================

type FoodPartnerSummary struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Logo             string  `json:"logo,omitempty"`
	Rating           float64 `json:"rating,omitempty"`
	DeliveryRadiusKM float64 `json:"deliveryRadiusKm,omitempty"`
}

// ============================================================
// FOOD RESPONSE
// ============================================================

type FoodResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	VideoURL     string `json:"videoUrl"`
	ImageURL     string `json:"imageUrl,omitempty"`
	ThumbnailURL string `json:"thumbnailUrl,omitempty"`
	Category     string `json:"category,omitempty"`
	Available    bool   `json:"available"`
	Likes        int    `json:"likes"`
	Saves        int    `json:"saves"`
	FoodPartner  any    `json:"foodPartner,omitempty"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func NewFoodResponse(
	food *models.Food,
	populatePartner bool,
) FoodResponse {

	description := ""

	if food.Description != nil {
		description = *food.Description
	}

	resp := FoodResponse{
		ID:           food.ID.String(),
		Name:         food.Name,
		Description:  description,
		VideoURL:     food.VideoURL,
		ImageURL:     stringValue(food.ImageURL),
		ThumbnailURL: stringValue(food.ThumbnailURL),
		Category:     stringValue(food.Category),
		Available:    food.Available,
		Likes:        food.Likes,
		Saves:        food.Saves,
		CreatedAt:    utils.FormatMongoDate(food.CreatedAt),
		UpdatedAt:    utils.FormatMongoDate(food.UpdatedAt),
	}

	if populatePartner {
		resp.FoodPartner = FoodPartnerSummary{
			ID:               food.FoodPartner.ID.String(),
			Name:             food.FoodPartner.Name,
			Logo:             food.FoodPartner.ProfileImage,
			Rating:           food.FoodPartner.Rating,
			DeliveryRadiusKM: food.FoodPartner.DeliveryRadiusKM,
		}
	}

	return resp
}

// ============================================================
// FOOD PARTNER DETAIL RESPONSE
// ============================================================

type FoodPartnerDetailResponse struct {
	ID                string                 `json:"id"`
	Name              string                 `json:"name"`
	ContactName       string                 `json:"contactName"`
	Phone             string                 `json:"phone"`
	Address           string                 `json:"address"`
	Email             string                 `json:"email,omitempty"`
	Latitude          *float64               `json:"latitude,omitempty"`
	Longitude         *float64               `json:"longitude,omitempty"`
	DeliveryRadiusKM  float64                `json:"deliveryRadiusKm"`
	OpeningTime       string                 `json:"openingTime"`
	ClosingTime       string                 `json:"closingTime"`
	WorkingHours      WorkingHoursResponse   `json:"workingHours"`
	CurrentOpenStatus string                 `json:"currentOpenStatus"`
	Rating            float64                `json:"rating"`
	ProfileImage      string                 `json:"profileImage"`
	CoverImage        string                 `json:"coverImage"`
	FoodImages        []string               `json:"foodImages"`
	FoodReels         []FoodResponse         `json:"foodReels"`
	MenuImages        []string               `json:"menuImages"`
	MenuPDFURL        string                 `json:"menuPdfUrl"`
	ReviewsSummary    ReviewsSummaryResponse `json:"reviewsSummary"`
	FoodItems         []FoodResponse         `json:"foodItems"`
}

type WorkingHoursResponse struct {
	OpeningTime string `json:"openingTime"`
	ClosingTime string `json:"closingTime"`
}

type ReviewsSummaryResponse struct {
	Rating      float64 `json:"rating"`
	ReviewCount int     `json:"reviewCount"`
}

func NewFoodPartnerDetailResponse(
	partner *models.FoodPartner,
	foods []models.Food,
	openStatus string,
) FoodPartnerDetailResponse {

	foodItems := make([]FoodResponse, 0, len(foods))

	for i := range foods {
		foodItems = append(
			foodItems,
			NewFoodResponse(&foods[i], false),
		)
	}

	return FoodPartnerDetailResponse{
		ID:               partner.ID.String(),
		Name:             partner.Name,
		ContactName:      partner.ContactName,
		Phone:            partner.Phone,
		Address:          partner.Address,
		Email:            partner.Email,
		Latitude:         partner.Latitude,
		Longitude:        partner.Longitude,
		DeliveryRadiusKM: partner.DeliveryRadiusKM,
		OpeningTime:      partner.OpeningTime,
		ClosingTime:      partner.ClosingTime,

		WorkingHours: WorkingHoursResponse{
			OpeningTime: partner.OpeningTime,
			ClosingTime: partner.ClosingTime,
		},

		CurrentOpenStatus: openStatus,
		Rating:            partner.Rating,
		ProfileImage:      partner.ProfileImage,
		CoverImage:        partner.CoverImage,

		FoodImages: parseStringList(partner.FoodImages),

		FoodReels: foodItems,

		MenuImages: parseStringList(partner.MenuImages),

		MenuPDFURL: partner.MenuPDF,

		ReviewsSummary: ReviewsSummaryResponse{
			Rating:      partner.Rating,
			ReviewCount: 0,
		},

		FoodItems: foodItems,
	}
}

// ============================================================
// FEED RESPONSE
// ============================================================

type FeedItemResponse struct {
	Reel                     FoodResponse       `json:"reel"`
	Restaurant               FoodPartnerSummary `json:"restaurant"`
	DistanceKM               float64            `json:"distanceKm"`
	EstimatedDeliveryMinutes string             `json:"estimatedDeliveryMinutes"`
	IsOpen                   bool               `json:"isOpen"`
	OpenStatus               string             `json:"openStatus"`
	CanOrder                 bool               `json:"canOrder"`
}

type FeedResponse struct {
	Items      []FeedItemResponse `json:"items"`
	NextCursor string             `json:"nextCursor,omitempty"`
	HasMore    bool               `json:"hasMore"`
}

// ============================================================
// FOOD PARTNER STATS
// ============================================================

type FoodPartnerStatsResponse struct {
	TotalReels    int     `json:"totalReels"`
	TotalLikes    int     `json:"totalLikes"`
	TotalSaves    int     `json:"totalSaves"`
	AvgEngagement float64 `json:"avgEngagement"`
}

func NewFoodPartnerStatsResponse(
	totalReels,
	totalLikes,
	totalSaves int,
) FoodPartnerStatsResponse {

	var avgEngagement float64

	if totalReels > 0 {
		avgEngagement =
			float64(totalLikes+totalSaves) /
				float64(totalReels)
	}

	return FoodPartnerStatsResponse{
		TotalReels:    totalReels,
		TotalLikes:    totalLikes,
		TotalSaves:    totalSaves,
		AvgEngagement: avgEngagement,
	}
}

// ============================================================
// HELPERS
// ============================================================

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func parseStringList(raw string) []string {
	if raw == "" {
		return []string{}
	}

	var values []string

	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		return []string{}
	}

	return values
}
