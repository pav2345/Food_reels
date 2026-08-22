package middleware

import (
	"net/http"
	"strings"

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

// ============================================================
// TOKEN EXTRACTION
// ============================================================
//
// Priority:
// 1. Authorization: Bearer <token>
// 2. Cookie: token=<token>
//

func getAuthToken(c *gin.Context) string {
	// 1. Try Authorization header.
	authHeader := strings.TrimSpace(
		c.GetHeader("Authorization"),
	)

	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimSpace(
			strings.TrimPrefix(authHeader, "Bearer "),
		)

		if token != "" {
			return token
		}
	}

	// 2. Fall back to cookie.
	cookieToken, err := c.Cookie("token")
	if err == nil && cookieToken != "" {
		return cookieToken
	}

	return ""
}

// ============================================================
// USER AUTHENTICATION
// ============================================================

func (m *AuthMiddleware) AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAuthToken(c)

		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Please login first",
				},
			)
			return
		}

		userID, err := utils.ParseToken(
			m.jwtSecret,
			token,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Invalid or expired token",
				},
			)
			return
		}

		user, err := m.userRepo.FindByID(
			c.Request.Context(),
			userID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Invalid or expired token",
				},
			)
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "User not found",
				},
			)
			return
		}

		c.Set("user", user)

		c.Next()
	}
}

// ============================================================
// FOOD PARTNER AUTHENTICATION
// ============================================================

func (m *AuthMiddleware) AuthFoodPartner() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getAuthToken(c)

		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Please login first",
				},
			)
			return
		}

		partnerID, err := utils.ParseToken(
			m.jwtSecret,
			token,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Invalid or expired token",
				},
			)
			return
		}

		partner, err := m.foodPartnerRepo.FindByID(
			c.Request.Context(),
			partnerID,
		)

		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Invalid or expired token",
				},
			)
			return
		}

		if partner == nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"message": "Food partner not found",
				},
			)
			return
		}

		c.Set("foodPartner", partner)

		c.Next()
	}
}