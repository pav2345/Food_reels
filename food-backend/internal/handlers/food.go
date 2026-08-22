package handlers

import (
	"log/slog"
	"net/http"

	"food-backend/internal/dto"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodService *services.FoodService
}

func NewFoodHandler(foodService *services.FoodService) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
	}
}

// CreateFood creates a new food reel for the authenticated food partner.
//
// Request:
// Content-Type: multipart/form-data
//
// Fields:
//   - video: File (required)
//   - name: Text (required)
//   - description: Text (optional)
func (h *FoodHandler) CreateFood(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	// Read uploaded video.
	file, err := readUploadedFile(c, "video")
	if err != nil || len(file) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Video file is required",
		})
		return
	}

	// Read form fields.
	name := c.PostForm("name")
	description := c.PostForm("description")

	// Validate food name.
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Food name is required",
		})
		return
	}

	// Create food/reel.
	food, err := h.foodService.CreateFood(
		c.Request.Context(),
		partner.ID,
		name,
		description,
		file,
	)
	if err != nil {
		slog.Error(
			"create food failed",
			"error", err,
			"partner_id", partner.ID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create food",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Food created successfully",
		"food":    dto.NewFoodResponse(food, false),
	})
}

// GetFoodItems returns all food items/reels.
func (h *FoodHandler) GetFoodItems(c *gin.Context) {
	foods, err := h.foodService.GetFoodItems(c.Request.Context())
	if err != nil {
		slog.Error(
			"get food items failed",
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch food items",
		})
		return
	}

	response := make([]dto.FoodResponse, 0, len(foods))

	for i := range foods {
		response = append(
			response,
			dto.NewFoodResponse(&foods[i], true),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"foods": response,
	})
}

// LikeFood toggles like/unlike for a food item.
func (h *FoodHandler) LikeFood(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	var req dto.FoodActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error(
			"like food invalid request",
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	foodID, err := parseFoodID(req.FoodID)
	if err != nil {
		slog.Error(
			"like food invalid food id",
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid food ID",
		})
		return
	}

	unliked, err := h.foodService.ToggleLike(
		c.Request.Context(),
		user.ID,
		foodID,
	)
	if err != nil {
		slog.Error(
			"like food failed",
			"error", err,
			"user_id", user.ID,
			"food_id", foodID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to like food",
		})
		return
	}

	if unliked {
		c.JSON(http.StatusOK, gin.H{
			"message": "Food unliked successfully",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Food liked successfully",
	})
}

// SaveFood toggles save/unsave for a food item.
func (h *FoodHandler) SaveFood(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	var req dto.FoodActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error(
			"save food invalid request",
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	foodID, err := parseFoodID(req.FoodID)
	if err != nil {
		slog.Error(
			"save food invalid food id",
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid food ID",
		})
		return
	}

	unsaved, err := h.foodService.ToggleSave(
		c.Request.Context(),
		user.ID,
		foodID,
	)
	if err != nil {
		slog.Error(
			"save food failed",
			"error", err,
			"user_id", user.ID,
			"food_id", foodID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save food",
		})
		return
	}

	if unsaved {
		c.JSON(http.StatusOK, gin.H{
			"message": "Food unsaved successfully",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Food saved successfully",
	})
}

// GetSaveFood returns the authenticated user's saved foods.
func (h *FoodHandler) GetSaveFood(c *gin.Context) {
	user, ok := getUserFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	foods, err := h.foodService.GetSavedFoods(
		c.Request.Context(),
		user.ID,
	)
	if err != nil {
		slog.Error(
			"get saved food failed",
			"error", err,
			"user_id", user.ID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch saved food",
		})
		return
	}

	response := make([]dto.FoodResponse, 0, len(foods))

	for i := range foods {
		response = append(
			response,
			dto.NewFoodResponse(&foods[i], true),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"foods": response,
	})
}

// GetFoodPartnerStats returns statistics for the authenticated food partner.
func (h *FoodHandler) GetFoodPartnerStats(c *gin.Context) {
	partner, ok := getFoodPartnerFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	totalReels, totalLikes, totalSaves, err :=
		h.foodService.GetFoodPartnerStats(
			c.Request.Context(),
			partner.ID,
		)

	if err != nil {
		slog.Error(
			"stats failed",
			"error", err,
			"partner_id", partner.ID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch stats",
		})
		return
	}

	c.JSON(http.StatusOK,
		dto.NewFoodPartnerStatsResponse(
			totalReels,
			totalLikes,
			totalSaves,
		),
	)
}
