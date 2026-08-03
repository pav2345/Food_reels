package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	JWT      JWTConfig
	Storage  StorageConfig
	CORS     CORSConfig
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret string
}

type StorageConfig struct {
	ImageKitPublicKey  string
	ImageKitPrivateKey string
	ImageKitURLEndpoint string
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

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

	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{"http://localhost:5174"}
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	return &Config{
		Port: port,
		Database: DatabaseConfig{
			URL: databaseURL,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
		},
		Storage: StorageConfig{
			ImageKitPublicKey:   os.Getenv("IMAGEKIT_PUBLIC_KEY"),
			ImageKitPrivateKey:  os.Getenv("IMAGEKIT_PRIVATE_KEY"),
			ImageKitURLEndpoint: os.Getenv("IMAGEKIT_URL_ENDPOINT"),
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
	}, nil
}

func (c *Config) ServerAddress() string {
	if _, err := strconv.Atoi(c.Port); err != nil {
		return ":" + c.Port
	}
	return ":" + c.Port
}
