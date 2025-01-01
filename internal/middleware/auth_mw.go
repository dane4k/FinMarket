package middleware

import (
	"github.com/dane4k/FinMarket/internal/auth"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, valid := IsAuthed(c)
		if !valid {
			c.Redirect(http.StatusFound, "/auth")
			c.Abort()
			return
		}

		user, err := repository.GetUser(int(userID))
		if err != nil || user.TgID == 0 {
			c.JSON(http.StatusNotFound, gin.H{"err": "user not found"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func IsAuthed(c *gin.Context) (int64, bool) {
	token, err := c.Cookie("jwtToken")
	if err != nil {
		return 0, false
	}

	jti, err := auth.ExtractJTI(token)
	if err != nil {
		return 0, false
	}

	isValid, err := repository.IsJTIValid(jti)
	if err != nil || !isValid {
		return 0, false
	}

	userID, err := auth.ParseUIDFromJWT(token)
	if err != nil {
		return 0, false
	}

	return userID, true
}
