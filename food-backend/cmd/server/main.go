package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"food-backend/internal/config"
	"food-backend/internal/database"
	"food-backend/internal/handlers"
	"food-backend/internal/middleware"
	"food-backend/internal/repository"
	"food-backend/internal/routes"
	"food-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.Database.URL)
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	foodPartnerRepo := repository.NewFoodPartnerRepository(db)
	foodRepo := repository.NewFoodRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	saveRepo := repository.NewSaveRepository(db)

	authService := services.NewAuthService(userRepo, foodPartnerRepo, cfg.JWT.Secret)
	storageService := services.NewStorageService(cfg.Storage)
	foodService := services.NewFoodService(db, foodRepo, likeRepo, saveRepo, storageService)
	foodPartnerService := services.NewFoodPartnerService(foodPartnerRepo, foodRepo)

	authHandler := handlers.NewAuthHandler(authService)
	foodHandler := handlers.NewFoodHandler(foodService)
	foodPartnerHandler := handlers.NewFoodPartnerHandler(foodPartnerService)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, userRepo, foodPartnerRepo)

	router := gin.New()
	router.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return true
			}
			for _, allowed := range cfg.CORS.AllowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	routes.Register(router, routes.Dependencies{
		AuthHandler:        authHandler,
		FoodHandler:        foodHandler,
		FoodPartnerHandler: foodPartnerHandler,
		AuthMiddleware:     authMiddleware,
	})

	server := &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: router,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("server shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}

	if err := database.Close(db); err != nil {
		slog.Error("database close failed", "error", err)
	}

	slog.Info("server stopped")
}
