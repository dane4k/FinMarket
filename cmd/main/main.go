package main

import (
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/logrs"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/route"
	"github.com/dane4k/FinMarket/internal/tg_bot"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	logrs.InitLogger()

	db.InitDB()

	go tg_bot.StartTelegramBot()
	pgdb.InitTGBot()

	router := gin.Default()
	router.Static("/static", "./internal/web/static")
	router.LoadHTMLGlob("internal/web/templates/*")
	route.InitializeRoutes(router)

	err := router.Run(":8080")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start the application")
		return
	}

	logrus.Info("Starting application")
}
