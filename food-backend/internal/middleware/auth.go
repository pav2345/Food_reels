package middleware

import (
	"net/http"

	"food-backend/internal/repository"
	"food-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtSecret       string
	userRepo        *repository.UserRepository
	foodPartnerRepo *repository.FoodPartnerRepository
}

func NewAuthMiddleware(
	jwtSecret string,
	userRepo *repository.UserRepository,
	foodPartnerRepo *repository.FoodPartnerRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret:       jwtSecret,
		userRepo:        userRepo,
		foodPartnerRepo: foodPartnerRepo,
	}
}

func (m *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
			return
		}

		userID, err := utils.ParseToken(m.jwtSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		user, err := m.userRepo.FindByID(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func (m *AuthMiddleware) AuthFoodPartner() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first"})
			return
		}

		partnerID, err := utils.ParseToken(m.jwtSecret, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		partner, err := m.foodPartnerRepo.FindByID(c.Request.Context(), partnerID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}
		if partner == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Food partner not found"})
			return
		}

		c.Set("foodPartner", partner)
		c.Next()
	}
}
