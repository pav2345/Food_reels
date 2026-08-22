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
	FeedHandler        *handlers.FeedHandler
	AuthMiddleware     *middleware.AuthMiddleware
}

func Register(router *gin.Engine, deps Dependencies) {
	// Health check
	router.GET("/", func(c *gin.Context) {
		c.String(200, "FoodReels API is running 🚀")
	})

	// ============================================================
	// API V1
	// ============================================================

	api := router.Group("/api/v1")
	{
		// ========================================================
		// AUTHENTICATION
		// ========================================================

		auth := api.Group("/auth")
		{
			// -------------------------------
			// User Authentication
			// -------------------------------

			user := auth.Group("/user")
			{
				user.POST(
					"/register",
					deps.AuthHandler.RegisterUser,
				)

				user.POST(
					"/login",
					deps.AuthHandler.LoginUser,
				)

				user.GET(
					"/logout",
					deps.AuthHandler.LogoutUser,
				)
			}

			// -------------------------------
			// Food Partner Authentication
			// -------------------------------

			partner := auth.Group("/food-partner")
			{
				partner.POST(
					"/register",
					deps.AuthHandler.RegisterFoodPartner,
				)

				partner.POST(
					"/login",
					deps.AuthHandler.LoginFoodPartner,
				)

				partner.GET(
					"/logout",
					deps.AuthHandler.LogoutFoodPartner,
				)
			}
		}

		// ========================================================
		// FOOD
		// ========================================================

		food := api.Group("/food")
		{
			// Food Partner only
			food.POST(
				"/",
				deps.AuthMiddleware.AuthFoodPartner(),
				deps.FoodHandler.CreateFood,
			)

			// User only
			food.GET(
				"/",
				deps.AuthMiddleware.AuthUser(),
				deps.FoodHandler.GetFoodItems,
			)

			food.POST(
				"/like",
				deps.AuthMiddleware.AuthUser(),
				deps.FoodHandler.LikeFood,
			)

			food.POST(
				"/save",
				deps.AuthMiddleware.AuthUser(),
				deps.FoodHandler.SaveFood,
			)

			food.GET(
				"/save",
				deps.AuthMiddleware.AuthUser(),
				deps.FoodHandler.GetSaveFood,
			)

			// Food Partner only
			food.GET(
				"/partner/stats",
				deps.AuthMiddleware.AuthFoodPartner(),
				deps.FoodHandler.GetFoodPartnerStats,
			)
		}

		// ========================================================
		// USER FEED
		// ========================================================

		api.GET(
			"/feed",
			deps.AuthMiddleware.AuthUser(),
			deps.FeedHandler.GetFeed,
		)

		api.PATCH(
			"/user/location",
			deps.AuthMiddleware.AuthUser(),
			deps.FeedHandler.UpdateUserLocation,
		)

		// ========================================================
		// FOOD PARTNER
		// ========================================================

		foodPartner := api.Group("/food-partner")
		{
			// Public restaurant profile.
			// No authentication required.
			foodPartner.GET(
				"/:id",
				deps.FoodPartnerHandler.GetFoodPartnerByID,
			)

			// ====================================================
			// FOOD PARTNER DASHBOARD
			// ====================================================

			dashboard := foodPartner.Group(
				"/dashboard",
				deps.AuthMiddleware.AuthFoodPartner(),
			)
			{
				// Create reel
				dashboard.POST(
					"/reel",
					deps.FoodHandler.CreateFood,
				)

				// Upload food image
				dashboard.POST(
					"/food-image",
					deps.FoodPartnerHandler.UploadFoodImage,
				)

				// Upload menu image
				dashboard.POST(
					"/menu-image",
					deps.FoodPartnerHandler.UploadMenuImage,
				)

				// Upload menu PDF
				dashboard.POST(
					"/menu-pdf",
					deps.FoodPartnerHandler.UploadMenuPDF,
				)

				// Update delivery radius
				dashboard.PATCH(
					"/delivery-radius",
					deps.FoodPartnerHandler.UpdateDeliveryRadius,
				)

				// Update working hours
				dashboard.PATCH(
					"/working-hours",
					deps.FoodPartnerHandler.UpdateWorkingHours,
				)

				// Update profile
				dashboard.PATCH(
					"/profile",
					deps.FoodPartnerHandler.UpdateProfile,
				)

				// Analytics
				dashboard.GET(
					"/analytics",
					deps.FoodHandler.GetFoodPartnerStats,
				)
			}
		}
	}
} 