package route

import (
	"github.com/dane4k/FinMarket/app"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, app *app.App) {
	router.GET("/", app.AuthMw.Handle(), app.HomeHandler.LoadHomeHandler)

	router.GET("/auth", app.AuthHandler.Authorize)

	router.GET("/check-status/:token", app.AuthHandler.CheckStatus)

	router.GET("/profile", app.AuthMw.Handle(), app.UserHandler.ShowProfileHandler)

	router.GET("/logout", app.UserHandler.LogoutHandler)

	router.POST("/api/user/:userID/update-avatar", app.AuthMw.Handle(), app.UserHandler.UpdateAvatarHandler)

	router.GET("/add", app.AuthMw.Handle(), app.ProductHandler.InputProductHandler)

	router.POST("/add/product", app.AuthMw.Handle(), app.ProductHandler.CreateProductHandler)

	router.PUT("/update/product/:id", app.ProductHandler.EditProductHandler)

	router.GET("/product/:id", app.ProductHandler.GetProductHandler)

	router.GET("/api/products", app.ProductHandler.GetProductsHandler)

	router.GET("/products/:category_id", app.AuthMw.Handle(), app.CategoryHandler.GetProductsByCategory)
}
