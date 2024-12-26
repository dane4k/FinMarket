package routes

import (
	"github.com/dane4k/FinMarket/internal/handlers"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(r *gin.Engine) {
	r.GET("/", handlers.LoadHome)

	r.GET("/auth", handlers.AuthUser)

	r.GET("/check-status/:token", handlers.CheckStatus)

	r.GET("/profile", handlers.ShowProfile)

	r.GET("/logout", handlers.Logout)
}
