package routes

import (
	"github.com/dane4k/FinMarket/internal/handlers"
	"github.com/dane4k/FinMarket/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.GET("/", handlers.LoadHome)

	router.GET("/auth", handlers.AuthUser)

	router.GET("/check-status/:token", handlers.CheckStatus)

	router.GET("/profile", middlewares.AuthMiddleware(), handlers.ShowProfile)

	router.GET("/logout", handlers.Logout)

	router.POST("/api/user/:userID/update-avatar", handlers.UpdateAvatarHandler)
}
