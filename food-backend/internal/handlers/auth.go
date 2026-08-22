package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService     *services.AuthService
	locationService *services.LocationService
	feedService     *services.FeedService
}

func NewAuthHandler(
	authService *services.AuthService,
	locationService *services.LocationService,
	feedService *services.FeedService,
) *AuthHandler {
	return &AuthHandler{
		authService:     authService,
		locationService: locationService,
		feedService:     feedService,
	}
}

// ============================================================
// USER REGISTRATION
// ============================================================

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	_, token, err := h.authService.RegisterUser(
		c.Request.Context(),
		req,
	)

	if errors.Is(err, services.ErrUserAlreadyExists) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}

	if err != nil {
		slog.Error("register user failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
		})
		return
	}

	// Set authentication cookie.
	setAuthCookie(c, token)

	// Do not expose user details or feed data.
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// ============================================================
// USER LOGIN
// ============================================================

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	_, token, err := h.authService.LoginUser(
		c.Request.Context(),
		req,
	)

	if errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	if err != nil {
		slog.Error("login user failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
		})
		return
	}

	// Set authentication cookie.
	setAuthCookie(c, token)

	// Return only the token.
	// No email, user ID, full name, or feed.
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// ============================================================
// USER LOGOUT
// ============================================================

func (h *AuthHandler) LogoutUser(c *gin.Context) {
	clearAuthCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "User logged out successfully",
	})
}

// ============================================================
// FOOD PARTNER REGISTRATION
// ============================================================

func (h *AuthHandler) RegisterFoodPartner(c *gin.Context) {
	var req dto.RegisterFoodPartnerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	_, token, err := h.authService.RegisterFoodPartner(
		c.Request.Context(),
		req,
	)

	if errors.Is(err, services.ErrFoodPartnerAlreadyExists) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Food partner already exists",
		})
		return
	}

	if err != nil {
		slog.Error("register food partner failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
		})
		return
	}

	// Set authentication cookie.
	setAuthCookie(c, token)

	// Do not expose partner details.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Food partner registered successfully",
	})
}

// ============================================================
// FOOD PARTNER LOGIN
// ============================================================

func (h *AuthHandler) LoginFoodPartner(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	_, token, err := h.authService.LoginFoodPartner(
		c.Request.Context(),
		req,
	)

	if errors.Is(err, services.ErrInvalidCredentials) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid email or password",
		})
		return
	}

	if err != nil {
		slog.Error("login food partner failed", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Server error",
		})
		return
	}

	// Set authentication cookie.
	setAuthCookie(c, token)

	// Return only the token.
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// ============================================================
// FOOD PARTNER LOGOUT
// ============================================================

func (h *AuthHandler) LogoutFoodPartner(c *gin.Context) {
	clearAuthCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Food partner logged out successfully",
	})
}

// ============================================================
// AUTH COOKIE
// ============================================================

func setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		"token",
		token,
		7*24*60*60,
		"/",
		"",
		true,
		true,
	)
}

func clearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)
}

// ============================================================
// AUTH FEED GENERATION
// ============================================================
//
// This function is kept because it may be used by other parts
// of the application. It is no longer called during registration
// or login because authentication responses should stay minimal.
//

func (h *AuthHandler) generateAuthFeed(
	c *gin.Context,
	user *models.User,
	latitude,
	longitude *float64,
) (dto.FeedResponse, error) {

	if headerLat, headerLng := optionalHeaderCoordinate(c); latitude == nil && longitude == nil {
		latitude = headerLat
		longitude = headerLng
	}

	coordinate, err := h.locationService.UpdateLoginLocation(
		c.Request.Context(),
		user,
		services.LocationInput{
			Latitude:  latitude,
			Longitude: longitude,
			IPAddress: c.ClientIP(),
		},
	)

	if errors.Is(err, services.ErrLocationUnavailable) {
		return dto.FeedResponse{
			Items: []dto.FeedItemResponse{},
		}, nil
	}

	if err != nil {
		return dto.FeedResponse{}, err
	}

	return h.feedService.Generate(
		c.Request.Context(),
		user.ID,
		*coordinate,
		"",
		20,
	)
}

// ============================================================
// LOCATION HEADER HELPERS
// ============================================================

func optionalHeaderCoordinate(c *gin.Context) (*float64, *float64) {
	latRaw := firstNonEmpty(
		c.GetHeader("X-Latitude"),
		c.GetHeader("X-Geo-Latitude"),
	)

	lngRaw := firstNonEmpty(
		c.GetHeader("X-Longitude"),
		c.GetHeader("X-Geo-Longitude"),
	)

	if latRaw == "" || lngRaw == "" {
		return nil, nil
	}

	lat, err := strconv.ParseFloat(latRaw, 64)
	if err != nil {
		return nil, nil
	}

	lng, err := strconv.ParseFloat(lngRaw, 64)
	if err != nil {
		return nil, nil
	}

	return &lat, &lng
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

// ============================================================
// CONTEXT HELPERS
// ============================================================

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

// ============================================================
// FOOD ID
// ============================================================

func parseFoodID(raw string) (uuid.UUID, error) {
	return uuid.Parse(raw)
}

// ============================================================
// FILE UPLOAD HELPERS
// ============================================================

func readUploadedFile(
	c *gin.Context,
	fieldName string,
) ([]byte, error) {

	file, _, _, err := readUploadedFileWithInfo(
		c,
		fieldName,
	)

	return file, err
}

func readUploadedFileWithInfo(
	c *gin.Context,
	fieldName string,
) ([]byte, string, string, error) {

	fileHeader, err := c.FormFile(fieldName)

	if err != nil {
		return nil, "", "", err
	}

	file, err := fileHeader.Open()

	if err != nil {
		return nil, "", "", err
	}

	defer file.Close()

	contents, err := io.ReadAll(file)

	if err != nil {
		return nil, "", "", err
	}

	return contents,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		nil
}