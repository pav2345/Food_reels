package handlers

import (
	"log/slog"
	"net/http"

	"food-backend/internal/dto"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FoodPartnerHandler struct {
	foodPartnerService *services.FoodPartnerService
	openService        *services.RestaurantOpenService
}

func NewFoodPartnerHandler(
	foodPartnerService *services.FoodPartnerService,
	openService *services.RestaurantOpenService,
) *FoodPartnerHandler {
	return &FoodPartnerHandler{
		foodPartnerService: foodPartnerService,
		openService:        openService,
	}
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
		"foodPartner": dto.NewFoodPartnerDetailResponse(partner, foods, h.openService.Status(*partner)),
	})
}

func (h *FoodPartnerHandler) UploadFoodImage(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	file, fileName, mimeType, err := readUploadedFileWithInfo(c, "image")
	if err != nil || len(file) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Food image is required"})
		return
	}

	url, err := h.foodPartnerService.UploadFoodImage(c.Request.Context(), *partner, file, fileName, mimeType)
	if err != nil {
		slog.Error("upload food image failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload food image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Food image uploaded successfully", "url": url})
}

func (h *FoodPartnerHandler) UploadMenuImage(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	file, fileName, mimeType, err := readUploadedFileWithInfo(c, "image")
	if err != nil || len(file) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Menu image is required"})
		return
	}

	url, err := h.foodPartnerService.UploadMenuImage(c.Request.Context(), *partner, file, fileName, mimeType)
	if err != nil {
		slog.Error("upload menu image failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload menu image"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu image uploaded successfully", "url": url})
}

func (h *FoodPartnerHandler) UploadMenuPDF(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	file, fileName, mimeType, err := readUploadedFileWithInfo(c, "pdf")
	if err != nil || len(file) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Menu PDF is required"})
		return
	}

	url, err := h.foodPartnerService.UploadMenuPDF(c.Request.Context(), partner.ID, file, fileName, mimeType)
	if err != nil {
		slog.Error("upload menu pdf failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload menu PDF"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Menu PDF uploaded successfully", "url": url})
}

func (h *FoodPartnerHandler) UpdateDeliveryRadius(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	var req dto.UpdateDeliveryRadiusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	updated, err := h.foodPartnerService.UpdateDeliveryRadius(c.Request.Context(), partner.ID, req.DeliveryRadiusKM)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery radius updated successfully", "foodPartner": dto.NewFoodPartnerDetailResponse(updated, nil, h.openService.Status(*updated))})
}

func (h *FoodPartnerHandler) UpdateWorkingHours(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	var req dto.UpdateWorkingHoursRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	updated, err := h.foodPartnerService.UpdateWorkingHours(c.Request.Context(), partner.ID, req.OpeningTime, req.ClosingTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Working hours updated successfully", "foodPartner": dto.NewFoodPartnerDetailResponse(updated, nil, h.openService.Status(*updated))})
}

func (h *FoodPartnerHandler) UpdateProfile(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
		return
	}

	var req dto.UpdateFoodPartnerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	updated, err := h.foodPartnerService.UpdateProfile(c.Request.Context(), partner.ID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restaurant profile updated successfully", "foodPartner": dto.NewFoodPartnerDetailResponse(updated, nil, h.openService.Status(*updated))})
}
