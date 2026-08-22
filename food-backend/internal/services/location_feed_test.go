package services

import (
	"context"
	"testing"

	"food-backend/internal/models"

	"github.com/google/uuid"
)

func TestResolveUsesCurrentLocationBeforeSavedLocation(t *testing.T) {
	svc := NewLocationService(nil, NewGeoService(), nil, 20)
	currentLat := 17.385044
	currentLng := 78.486671
	savedLat := 12.971599
	savedLng := 77.594566

	user := &models.User{
		CurrentLatitude:  &currentLat,
		CurrentLongitude: &currentLng,
		SavedLatitude:    &savedLat,
		SavedLongitude:   &savedLng,
	}

	coord, err := svc.Resolve(context.Background(), user, LocationInput{})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if coord == nil {
		t.Fatal("Resolve returned nil coord")
	}
	if coord.Latitude != currentLat || coord.Longitude != currentLng {
		t.Fatalf("expected current location, got lat=%v lng=%v", coord.Latitude, coord.Longitude)
	}
}

func TestSelectFeedCandidatesBalancesLocalAndBroader(t *testing.T) {
	svc := NewFeedService(nil, nil, nil, NewGeoService(), nil, nil, nil, 6, 60)
	candidates := []FeedCandidate{
		{Food: models.Food{ID: mustUUID("11111111-1111-1111-1111-111111111111")}, DistanceKM: 1},
		{Food: models.Food{ID: mustUUID("22222222-2222-2222-2222-222222222222")}, DistanceKM: 2},
		{Food: models.Food{ID: mustUUID("33333333-3333-3333-3333-333333333333")}, DistanceKM: 7},
		{Food: models.Food{ID: mustUUID("44444444-4444-4444-4444-444444444444")}, DistanceKM: 8},
		{Food: models.Food{ID: mustUUID("55555555-5555-5555-5555-555555555555")}, DistanceKM: 9},
	}

	selected := svc.selectFeedCandidates(candidates, 5, 6, 60)
	if len(selected) != 5 {
		t.Fatalf("expected 5 selected candidates, got %d", len(selected))
	}
	localCount := 0
	for _, candidate := range selected {
		if candidate.DistanceKM <= 6 {
			localCount++
		}
	}
	if localCount != 2 {
		t.Fatalf("expected 2 local candidates, got %d", localCount)
	}
}

func mustUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		panic(err)
	}
	return id
}
