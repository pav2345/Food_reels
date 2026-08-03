package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const authCookieName = "token"
const authCookieMaxAge = 7 * 24 * 60 * 60

func SetAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(authCookieName, token, authCookieMaxAge, "/", "", true, true)
}

func ClearAuthCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(authCookieName, "", -1, "/", "", true, true)
}

func GetAuthCookie(c *gin.Context) (string, error) {
	return c.Cookie(authCookieName)
}
