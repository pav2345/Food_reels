package services

import (
	"strings"
	"time"

	"food-backend/internal/models"
)

const (
	RestaurantStatusOpen   = "OPEN"
	RestaurantStatusClosed = "CLOSED"
)

type RestaurantOpenService struct {
	now func() time.Time
}

func NewRestaurantOpenService() *RestaurantOpenService {
	return &RestaurantOpenService{now: time.Now}
}

func (s *RestaurantOpenService) Status(partner models.FoodPartner) string {
	if s.IsOpen(partner) {
		return RestaurantStatusOpen
	}
	return RestaurantStatusClosed
}

func (s *RestaurantOpenService) IsOpen(partner models.FoodPartner) bool {
	openingTime := strings.TrimSpace(partner.OpeningTime)
	closingTime := strings.TrimSpace(partner.ClosingTime)
	if openingTime == "" || closingTime == "" {
		return false
	}

	now := s.now()
	openAt, err := parseBusinessClock(now, openingTime)
	if err != nil {
		return false
	}
	closeAt, err := parseBusinessClock(now, closingTime)
	if err != nil {
		return false
	}

	if closeAt.Equal(openAt) {
		return true
	}
	if closeAt.Before(openAt) {
		closeAt = closeAt.Add(24 * time.Hour)
		if now.Before(openAt) {
			now = now.Add(24 * time.Hour)
		}
	}

	return !now.Before(openAt) && now.Before(closeAt)
}

func parseBusinessClock(base time.Time, raw string) (time.Time, error) {
	parsed, err := time.Parse("15:04", raw)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		base.Year(),
		base.Month(),
		base.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		base.Location(),
	), nil
}
