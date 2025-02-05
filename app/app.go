package app

import (
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/db"
	"github.com/dane4k/FinMarket/internal/handler"
	"github.com/dane4k/FinMarket/internal/middleware"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service"
)

type App struct {
	AuthHandler     *handler.AuthHandler
	CategoryHandler *handler.CategoryHandler
	HomeHandler     *handler.HomeHandler
	ProductHandler  *handler.ProductHandler
	UserHandler     *handler.UserHandler
	AuthMw          *middleware.AuthMiddleware
	ProductMw       *middleware.ProductMiddleware
	AuthRepository  pgdb.AuthRepository
	UserRepository  pgdb.UserRepository
}

func InitApp(cfg *config.Config) *App {
	dbObj := db.InitDB(cfg)

	authRepo := pgdb.NewAuthRepository(dbObj, cfg)
	categoryRepo := pgdb.NewCategoryRepository(dbObj, cfg)
	productRepo := pgdb.NewProductRepository(dbObj, cfg)
	userRepo := pgdb.NewUserRepository(dbObj, cfg)

	authService := service.NewAuthService(authRepo, userRepo, cfg)
	productService := service.NewProductService(productRepo, cfg)
	userService := service.NewUserService(userRepo, authRepo, cfg)

	_ = categoryRepo.GetAllCategories(&handler.Categories)

	return &App{
		AuthHandler:     handler.NewAuthHandler(authService),
		CategoryHandler: handler.NewCategoryHandler(),
		HomeHandler:     handler.NewHomeHandler(),
		ProductHandler:  handler.NewProductHandler(productService),
		UserHandler:     handler.NewUserHandler(userService),
		AuthMw:          middleware.NewAuthMiddleware(authService),
		ProductMw:       middleware.NewProductMiddleware(productService),
		AuthRepository:  authRepo,
		UserRepository:  userRepo,
	}
}
