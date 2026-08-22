package services

import (
	"context"
	"errors"
	"net"
	"time"

	"food-backend/internal/models"
	"food-backend/internal/repository"

	"github.com/google/uuid"
)

var ErrLocationUnavailable = errors.New("location unavailable")

type LocationInput struct {
	Latitude  *float64
	Longitude *float64
	IPAddress string
}

type IPGeoProvider interface {
	Lookup(ctx context.Context, ipAddress string) (*Coordinate, error)
}

type NoopIPGeoProvider struct{}

func (NoopIPGeoProvider) Lookup(ctx context.Context, ipAddress string) (*Coordinate, error) {
	_ = ctx
	_ = ipAddress
	return nil, ErrLocationUnavailable
}

type LocationService struct {
	userRepo         *repository.UserRepository
	geoService       *GeoService
	ipGeoProvider    IPGeoProvider
	updateDistanceKM float64
}

func NewLocationService(
	userRepo *repository.UserRepository,
	geoService *GeoService,
	ipGeoProvider IPGeoProvider,
	updateDistanceKM float64,
) *LocationService {
	if ipGeoProvider == nil {
		ipGeoProvider = NoopIPGeoProvider{}
	}

	if updateDistanceKM <= 0 {
		updateDistanceKM = 20
	}
	return &LocationService{
		userRepo:         userRepo,
		geoService:       geoService,
		ipGeoProvider:    ipGeoProvider,
		updateDistanceKM: updateDistanceKM,
	}
}

func (s *LocationService) UpdateLoginLocation(ctx context.Context, user *models.User, input LocationInput) (*Coordinate, error) {
	coordinate, err := s.Resolve(ctx, user, input)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if user.SavedLatitude == nil || user.SavedLongitude == nil {
		if err := s.userRepo.InitializeLocation(ctx, user.ID, coordinate.Latitude, coordinate.Longitude, now); err != nil {
			return nil, err
		}
	} else {
		if err := s.userRepo.UpdateLocation(ctx, user.ID, coordinate.Latitude, coordinate.Longitude, now); err != nil {
			return nil, err
		}
	}

	user.Latitude = &coordinate.Latitude
	user.Longitude = &coordinate.Longitude
	user.CurrentLatitude = &coordinate.Latitude
	user.CurrentLongitude = &coordinate.Longitude
	user.SavedLatitude = &coordinate.Latitude
	user.SavedLongitude = &coordinate.Longitude
	user.LastLocationUpdated = &now
	user.CurrentLocationUpdatedAt = &now

	return coordinate, nil
}

func (s *LocationService) UpdateExplicitLocation(ctx context.Context, user *models.User, latitude, longitude float64) (*Coordinate, error) {
	coordinate, _, err := s.UpdateExplicitLocationWithStatus(ctx, user, latitude, longitude)
	return coordinate, err
}

func (s *LocationService) UpdateExplicitLocationWithStatus(ctx context.Context, user *models.User, latitude, longitude float64) (*Coordinate, bool, error) {
	if err := s.geoService.ValidateCoordinate(latitude, longitude); err != nil {
		return nil, false, err
	}

	coordinate := Coordinate{Latitude: latitude, Longitude: longitude}
	now := time.Now().UTC()
	if user.CurrentLatitude != nil && user.CurrentLongitude != nil {
		prev := Coordinate{Latitude: *user.CurrentLatitude, Longitude: *user.CurrentLongitude}
		if s.geoService.DistanceKM(prev, coordinate) >= s.updateDistanceKM {
			if err := s.userRepo.UpdateCurrentLocation(ctx, user.ID, latitude, longitude, now); err != nil {
				return nil, false, err
			}
			user.Latitude = &coordinate.Latitude
			user.Longitude = &coordinate.Longitude
			user.CurrentLatitude = &coordinate.Latitude
			user.CurrentLongitude = &coordinate.Longitude
			user.LastLocationUpdated = &now
			user.CurrentLocationUpdatedAt = &now
			return &coordinate, true, nil
		}

		user.Latitude = &coordinate.Latitude
		user.Longitude = &coordinate.Longitude
		user.LastLocationUpdated = &now
		return &coordinate, false, nil
	}

	if err := s.userRepo.InitializeLocation(ctx, user.ID, latitude, longitude, now); err != nil {
		return nil, false, err
	}

	user.Latitude = &coordinate.Latitude
	user.Longitude = &coordinate.Longitude
	user.CurrentLatitude = &coordinate.Latitude
	user.CurrentLongitude = &coordinate.Longitude
	user.SavedLatitude = &coordinate.Latitude
	user.SavedLongitude = &coordinate.Longitude
	user.LastLocationUpdated = &now
	user.CurrentLocationUpdatedAt = &now

	return &coordinate, true, nil
}

func (s *LocationService) Resolve(ctx context.Context, user *models.User, input LocationInput) (*Coordinate, error) {
	if input.Latitude != nil && input.Longitude != nil {
		if err := s.geoService.ValidateCoordinate(*input.Latitude, *input.Longitude); err != nil {
			return nil, err
		}
		return &Coordinate{Latitude: *input.Latitude, Longitude: *input.Longitude}, nil
	}

	if user.CurrentLatitude != nil && user.CurrentLongitude != nil {
		return &Coordinate{Latitude: *user.CurrentLatitude, Longitude: *user.CurrentLongitude}, nil
	}

	if user.SavedLatitude != nil && user.SavedLongitude != nil {
		return &Coordinate{Latitude: *user.SavedLatitude, Longitude: *user.SavedLongitude}, nil
	}

	if user.Latitude != nil && user.Longitude != nil {
		return &Coordinate{Latitude: *user.Latitude, Longitude: *user.Longitude}, nil
	}

	if input.IPAddress == "" || isPrivateIP(input.IPAddress) {
		return nil, ErrLocationUnavailable
	}

	coordinate, err := s.ipGeoProvider.Lookup(ctx, input.IPAddress)
	if err != nil {
		return nil, err
	}
	if coordinate == nil {
		return nil, ErrLocationUnavailable
	}
	if err := s.geoService.ValidateCoordinate(coordinate.Latitude, coordinate.Longitude); err != nil {
		return nil, err
	}

	return coordinate, nil
}

func isPrivateIP(raw string) bool {
	ip := net.ParseIP(raw)
	if ip == nil {
		return true
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified()
}

func (s *LocationService) GetUserLocation(ctx context.Context, userID uuid.UUID) (*Coordinate, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.CurrentLatitude == nil || user.CurrentLongitude == nil {
		if user == nil || user.SavedLatitude == nil || user.SavedLongitude == nil {
			if user == nil || user.Latitude == nil || user.Longitude == nil {
				return nil, ErrLocationUnavailable
			}
			return &Coordinate{Latitude: *user.Latitude, Longitude: *user.Longitude}, nil
		}
		return &Coordinate{Latitude: *user.SavedLatitude, Longitude: *user.SavedLongitude}, nil
	}
	return &Coordinate{Latitude: *user.CurrentLatitude, Longitude: *user.CurrentLongitude}, nil
}
