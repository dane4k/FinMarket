package middlewares

import (
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/auth"
	"github.com/dane4k/FinMarket/internal/models"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwtToken")
		if err != nil {
			c.Redirect(http.StatusFound, "/auth")
			c.Abort()
			return
		}

		jti, err := auth.ExtractJTI(token)
		if err != nil || !repository.IsTokenValid(jti) {
			c.Redirect(http.StatusFound, "/auth")
			c.Abort()
			return
		}

		userID, err := auth.ParseUIDFromJWT(token)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth")
			c.Abort()
			return
		}

		var user models.User
		if err := db.DB.Where("tg_id = ?", userID).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
			c.Abort()
			return
		}
		c.Set("user", user)

		c.Next()
	}
}
