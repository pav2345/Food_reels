package handlers

import (
	"net/http"

	"food-backend/internal/dto"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FoodPartnerHandler struct {
	foodPartnerService *services.FoodPartnerService
}

func NewFoodPartnerHandler(foodPartnerService *services.FoodPartnerService) *FoodPartnerHandler {
	return &FoodPartnerHandler{foodPartnerService: foodPartnerService}
}

func (h *FoodPartnerHandler) GetFoodPartnerByID(c *gin.Context) {
	foodPartnerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Food partner not found"})
		return
	}

	partner, foods, err := h.foodPartnerService.GetFoodPartnerByID(c.Request.Context(), foodPartnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	if partner == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Food partner not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Food partner retrieved successfully",
		"foodPartner": dto.NewFoodPartnerDetailResponse(partner, foods),
	})
}
