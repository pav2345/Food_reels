package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	user, token, err := h.authService.RegisterUser(c.Request.Context(), req)
	if errors.Is(err, services.ErrUserAlreadyExists) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User already exists"})
		return
	}
	if err != nil {
		slog.Error("register user failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    dto.NewUserResponse(user),
	})
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	user, token, err := h.authService.LoginUser(c.Request.Context(), req)
	if errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}
	if err != nil {
		slog.Error("login user failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{
		"message": "User logged in successfully",
		"user":    dto.NewUserResponse(user),
	})
}

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	clearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "User logged out successfully"})
}

func (h *AuthHandler) RegisterFoodPartner(c *gin.Context) {
	var req dto.RegisterFoodPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	partner, token, err := h.authService.RegisterFoodPartner(c.Request.Context(), req)
	if errors.Is(err, services.ErrFoodPartnerAlreadyExists) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Food partner already exists"})
		return
	}
	if err != nil {
		slog.Error("register food partner failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusCreated, gin.H{
		"message":     "Food partner registered successfully",
		"foodPartner": dto.NewFoodPartnerAuthResponse(partner),
	})
}

func (h *AuthHandler) LoginFoodPartner(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	partner, token, err := h.authService.LoginFoodPartner(c.Request.Context(), req)
	if errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}
	if err != nil {
		slog.Error("login food partner failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
			"error":   gin.H{},
		})
		return
	}

	setAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Food partner logged in successfully",
		"foodPartner": dto.NewFoodPartnerAuthResponse(partner),
	})
}

func (h *AuthHandler) LogoutFoodPartner(c *gin.Context) {
	clearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Food partner logged out successfully"})
}

func setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", token, 7*24*60*60, "/", "", true, true)
}

func clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", "", -1, "/", "", true, true)
}

func getUserFromContext(c *gin.Context) (*models.User, bool) {
	value, ok := c.Get("user")
	if !ok {
		return nil, false
	}
	user, ok := value.(*models.User)
	return user, ok
}

func getFoodPartnerFromContext(c *gin.Context) (*models.FoodPartner, bool) {
	value, ok := c.Get("foodPartner")
	if !ok {
		return nil, false
	}
	partner, ok := value.(*models.FoodPartner)
	return partner, ok
}

func parseFoodID(raw string) (uuid.UUID, error) {
	return uuid.Parse(raw)
}

func readUploadedFile(c *gin.Context, fieldName string) ([]byte, error) {
	fileHeader, err := c.FormFile(fieldName)
	if err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
