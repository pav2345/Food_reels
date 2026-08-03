package routes

import (
	"food-backend/internal/handlers"
	"food-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	AuthHandler        *handlers.AuthHandler
	FoodHandler        *handlers.FoodHandler
	FoodPartnerHandler *handlers.FoodPartnerHandler
	AuthMiddleware     *middleware.AuthMiddleware
}

func Register(router *gin.Engine, deps Dependencies) {
	router.GET("/", func(c *gin.Context) {
		c.String(200, "FoodReels API is running 🚀")
	})

	auth := router.Group("/api/auth")
	{
		auth.POST("/user/register", deps.AuthHandler.RegisterUser)
		auth.POST("/user/login", deps.AuthHandler.LoginUser)
		auth.GET("/user/logout", deps.AuthHandler.LogoutUser)

		auth.POST("/food-partner/register", deps.AuthHandler.RegisterFoodPartner)
		auth.POST("/food-partner/login", deps.AuthHandler.LoginFoodPartner)
		auth.GET("/food-partner/logout", deps.AuthHandler.LogoutFoodPartner)
	}

	food := router.Group("/api/food")
	{
		food.POST("/", deps.AuthMiddleware.AuthFoodPartner(), deps.FoodHandler.CreateFood)
		food.GET("/", deps.AuthMiddleware.AuthUser(), deps.FoodHandler.GetFoodItems)
		food.POST("/like", deps.AuthMiddleware.AuthUser(), deps.FoodHandler.LikeFood)
		food.POST("/save", deps.AuthMiddleware.AuthUser(), deps.FoodHandler.SaveFood)
		food.GET("/save", deps.AuthMiddleware.AuthUser(), deps.FoodHandler.GetSaveFood)
		food.GET("/partner/stats", deps.AuthMiddleware.AuthFoodPartner(), deps.FoodHandler.GetFoodPartnerStats)
	}

	foodPartner := router.Group("/api/food-partner")
	{
		foodPartner.GET("/:id", deps.AuthMiddleware.AuthUser(), deps.FoodPartnerHandler.GetFoodPartnerByID)
	}
}
