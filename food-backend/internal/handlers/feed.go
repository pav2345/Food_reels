package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"food-backend/internal/dto"
	"food-backend/internal/models"
	"food-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	locationService *services.LocationService
	feedService     *services.FeedService
}

func NewFeedHandler(
	locationService *services.LocationService,
	feedService *services.FeedService,
) *FeedHandler {
	return &FeedHandler{
		locationService: locationService,
		feedService:     feedService,
	}
}

// ============================================================
// GET FEED
// ============================================================

func (h *FeedHandler) GetFeed(c *gin.Context) {
	user, ok := getUserFromContext(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	limit := parseIntDefault(c.Query("limit"), 20)

	// Keep the feed size within a reasonable range.
	if limit < 1 {
		limit = 20
	}

	if limit > 50 {
		limit = 50
	}

	coordinate, err := h.resolveFeedLocation(c, user)

	if errors.Is(err, services.ErrLocationUnavailable) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User location is required to generate the feed",
		})
		return
	}

	if err != nil {
		slog.Error(
			"resolve feed location failed",
			"error", err,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to determine user location",
		})
		return
	}

	feed, err := h.feedService.Generate(
		c.Request.Context(),
		user.ID,
		*coordinate,
		c.Query("cursor"),
		limit,
	)

	if err != nil {
		slog.Error(
			"generate feed failed",
			"error", err,
			"user_id", user.ID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate feed",
		})
		return
	}

	c.JSON(http.StatusOK, feed)
}

// ============================================================
// LOCATION RESOLUTION
// ============================================================

func (h *FeedHandler) resolveFeedLocation(
	c *gin.Context,
	currentUser *models.User,
) (*services.Coordinate, error) {

	// Explicit coordinates supplied by the client.
	latitude, longitude := optionalQueryCoordinate(c)

	if latitude != nil && longitude != nil {
		return h.locationService.UpdateExplicitLocation(
			c.Request.Context(),
			currentUser,
			*latitude,
			*longitude,
		)
	}

	// Otherwise resolve from stored user location / IP.
	return h.locationService.Resolve(
		c.Request.Context(),
		currentUser,
		services.LocationInput{
			IPAddress: c.ClientIP(),
		},
	)
}

// ============================================================
// UPDATE USER LOCATION
// ============================================================

func (h *FeedHandler) UpdateUserLocation(c *gin.Context) {
	user, ok := getUserFromContext(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Please login first",
		})
		return
	}

	var req dto.UpdateLocationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	coordinate, changed, err :=
		h.locationService.UpdateExplicitLocationWithStatus(
			c.Request.Context(),
			user,
			req.Latitude,
			req.Longitude,
		)

	if err != nil {
		slog.Error(
			"update user location failed",
			"error", err,
			"user_id", user.ID,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	message := "Location unchanged"
	feedRefreshRequired := false

	if changed {
		message = "Location updated"
		feedRefreshRequired = true
	}

	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"location": gin.H{
			"latitude":  coordinate.Latitude,
			"longitude": coordinate.Longitude,
		},
		"feedRefreshRequired": feedRefreshRequired,
	})
}

// ============================================================
// QUERY COORDINATES
// ============================================================

func optionalQueryCoordinate(
	c *gin.Context,
) (*float64, *float64) {

	latRaw := firstNonEmpty(
		c.Query("latitude"),
		c.Query("lat"),
		c.GetHeader("X-Latitude"),
	)

	lngRaw := firstNonEmpty(
		c.Query("longitude"),
		c.Query("lng"),
		c.Query("lon"),
		c.GetHeader("X-Longitude"),
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

// ============================================================
// INTEGER QUERY PARAMETER
// ============================================================

func parseIntDefault(
	raw string,
	fallback int,
) int {

	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)

	if err != nil {
		return fallback
	}

	return value
}
