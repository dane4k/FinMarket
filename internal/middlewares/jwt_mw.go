package middlewares

import "github.com/gin-gonic/gin"

func CheckAuth(c *gin.Context) {
	c.Next()
}
