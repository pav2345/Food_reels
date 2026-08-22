package services

import (
	"errors"
	"math"
)

const DefaultFeedRadiusKM = 6.0

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type GeoService struct{}

func NewGeoService() *GeoService {
	return &GeoService{}
}

func (s *GeoService) ValidateCoordinate(
	latitude,
	longitude float64,
) error {

	if latitude < -90 || latitude > 90 {
		return errors.New(
			"latitude must be between -90 and 90",
		)
	}

	if longitude < -180 || longitude > 180 {
		return errors.New(
			"longitude must be between -180 and 180",
		)
	}

	return nil
}

func (s *GeoService) DistanceKM(
	from,
	to Coordinate,
) float64 {

	const earthRadiusKM = 6371.0

	lat1 := degreesToRadians(from.Latitude)
	lat2 := degreesToRadians(to.Latitude)

	deltaLat := degreesToRadians(
		to.Latitude - from.Latitude,
	)

	deltaLng := degreesToRadians(
		to.Longitude - from.Longitude,
	)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*
			math.Cos(lat2)*
			math.Sin(deltaLng/2)*
			math.Sin(deltaLng/2)

	c := 2 * math.Atan2(
		math.Sqrt(a),
		math.Sqrt(1-a),
	)

	return earthRadiusKM * c
}

func degreesToRadians(value float64) float64 {
	return value * math.Pi / 180
}
