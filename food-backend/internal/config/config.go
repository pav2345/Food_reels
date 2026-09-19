package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	JWT      JWTConfig
	Storage  StorageConfig
	CORS     CORSConfig
	Location LocationConfig
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
}

type StorageConfig struct {
	ImageKitPublicKey   string
	ImageKitPrivateKey  string
	ImageKitURLEndpoint string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type LocationConfig struct {
	UpdateDistanceKM  float64
	FeedLocalRadiusKM float64
	LocalFeedPercent  int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(
			"could not load .env file",
			"error", err,
		)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	imageKitPublicKey := os.Getenv("IMAGEKIT_PUBLIC_KEY")
	imageKitPrivateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	imageKitURLEndpoint := os.Getenv("IMAGEKIT_URL_ENDPOINT")

	slog.Info(
		"ImageKit configuration",
		"public_key_configured", imageKitPublicKey != "",
		"private_key_configured", imageKitPrivateKey != "",
		"url_endpoint_configured", imageKitURLEndpoint != "",
	)

	if imageKitPrivateKey == "" {
		return nil, fmt.Errorf("IMAGEKIT_PRIVATE_KEY is required")
	}

	if imageKitPublicKey == "" {
		return nil, fmt.Errorf("IMAGEKIT_PUBLIC_KEY is required")
	}

	if imageKitURLEndpoint == "" {
		return nil, fmt.Errorf("IMAGEKIT_URL_ENDPOINT is required")
	}

	locationUpdateDistanceKM, err := strconv.ParseFloat(
		os.Getenv("LOCATION_UPDATE_DISTANCE_KM"),
		64,
	)

	if err != nil || locationUpdateDistanceKM <= 0 {
		locationUpdateDistanceKM = 20
	}

	feedLocalRadiusKM, err := strconv.ParseFloat(
		os.Getenv("FEED_LOCAL_RADIUS_KM"),
		64,
	)

	if err != nil || feedLocalRadiusKM <= 0 {
		feedLocalRadiusKM = 6
	}

	localFeedPercent, err := strconv.Atoi(
		os.Getenv("LOCAL_FEED_PERCENTAGE"),
	)

	if err != nil ||
		localFeedPercent <= 0 ||
		localFeedPercent > 100 {
		localFeedPercent = 60
	}

	allowedOrigins := buildAllowedOrigins(
		os.Getenv("FRONTEND_URL"),
	)

	return &Config{
		Port: port,

		Database: DatabaseConfig{
			URL: databaseURL,
		},

		JWT: JWTConfig{
			Secret: jwtSecret,
		},

		Storage: StorageConfig{
			ImageKitPublicKey:   imageKitPublicKey,
			ImageKitPrivateKey:  imageKitPrivateKey,
			ImageKitURLEndpoint: imageKitURLEndpoint,
		},

		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},

		Location: LocationConfig{
			UpdateDistanceKM:  locationUpdateDistanceKM,
			FeedLocalRadiusKM: feedLocalRadiusKM,
			LocalFeedPercent:  localFeedPercent,
		},
	}, nil
}

func (c *Config) ServerAddress() string {
	return ":" + c.Port
}

func buildAllowedOrigins(frontendURL string) []string {
	origins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:5174",
		"http://127.0.0.1:5174",
	}

	for _, origin := range parseAllowedOrigins(frontendURL) {
		if !containsString(origins, origin) {
			origins = append(origins, origin)
		}
	}

	return origins
}

func parseAllowedOrigins(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")

	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		origin := strings.TrimSpace(part)

		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}

	return false
}
