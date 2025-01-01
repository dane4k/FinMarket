package route

import (
	"github.com/dane4k/FinMarket/internal/handler"
	"github.com/dane4k/FinMarket/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {
	router.GET("/", handler.LoadHomeHandler)

	router.GET("/auth", handler.AuthHandler)

	router.GET("/check-status/:token", handler.CheckStatusHandler)

	router.GET("/profile", middleware.AuthMiddleware(), handler.ShowProfileHandler)

	router.GET("/logout", handler.LogoutHandler)

	router.POST("/api/user/:userID/update-avatar", handler.UpdateAvatarHandler)
}
