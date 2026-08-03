package dto

import (
	"food-backend/internal/models"
	"food-backend/internal/utils"
)

type UserResponse struct {
	ID       string `json:"_id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
	}
}

type FoodPartnerAuthResponse struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewFoodPartnerAuthResponse(partner *models.FoodPartner) FoodPartnerAuthResponse {
	return FoodPartnerAuthResponse{
		ID:    partner.ID.String(),
		Name:  partner.Name,
		Email: partner.Email,
	}
}

type FoodPartnerSummary struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type FoodResponse struct {
	ID          string      `json:"_id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	VideoURL    string      `json:"videoUrl"`
	Likes       int         `json:"likes"`
	Saves       int         `json:"saves"`
	FoodPartner interface{} `json:"foodPartner"`
	Version     int         `json:"__v"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}

func NewFoodResponse(food *models.Food, populatePartner bool) FoodResponse {
	description := ""
	if food.Description != nil {
		description = *food.Description
	}

	resp := FoodResponse{
		ID:          food.ID.String(),
		Name:        food.Name,
		Description: description,
		VideoURL:    food.VideoURL,
		Likes:       food.Likes,
		Saves:       food.Saves,
		Version:     food.Version,
		CreatedAt:   utils.FormatMongoDate(food.CreatedAt),
		UpdatedAt:   utils.FormatMongoDate(food.UpdatedAt),
	}

	if populatePartner {
		resp.FoodPartner = FoodPartnerSummary{
			ID:   food.FoodPartner.ID.String(),
			Name: food.FoodPartner.Name,
		}
	} else {
		resp.FoodPartner = food.FoodPartnerID.String()
	}

	return resp
}

type FoodPartnerDetailResponse struct {
	ID          string         `json:"_id"`
	Name        string         `json:"name"`
	ContactName string         `json:"contactName"`
	Phone       string         `json:"phone"`
	Address     string         `json:"address"`
	Email       string         `json:"email"`
	Password    string         `json:"password"`
	Version     int            `json:"__v"`
	FoodItems   []FoodResponse `json:"foodItems"`
}

func NewFoodPartnerDetailResponse(partner *models.FoodPartner, foods []models.Food) FoodPartnerDetailResponse {
	foodItems := make([]FoodResponse, 0, len(foods))
	for i := range foods {
		foodItems = append(foodItems, NewFoodResponse(&foods[i], false))
	}

	return FoodPartnerDetailResponse{
		ID:          partner.ID.String(),
		Name:        partner.Name,
		ContactName: partner.ContactName,
		Phone:       partner.Phone,
		Address:     partner.Address,
		Email:       partner.Email,
		Password:    partner.Password,
		Version:     0,
		FoodItems:   foodItems,
	}
}

type FoodPartnerStatsResponse struct {
	TotalReels     int         `json:"totalReels"`
	TotalLikes     int         `json:"totalLikes"`
	TotalSaves     int         `json:"totalSaves"`
	AvgEngagement  interface{} `json:"avgEngagement"`
}

func NewFoodPartnerStatsResponse(totalReels, totalLikes, totalSaves int) FoodPartnerStatsResponse {
	var avgEngagement interface{} = 0
	if totalReels > 0 {
		avg := float64(totalLikes+totalSaves) / float64(totalReels)
		avgEngagement = formatAvgEngagement(avg)
	}

	return FoodPartnerStatsResponse{
		TotalReels:    totalReels,
		TotalLikes:    totalLikes,
		TotalSaves:    totalSaves,
		AvgEngagement: avgEngagement,
	}
}

func formatAvgEngagement(value float64) string {
	return utils.FormatFixed1(value)
}
