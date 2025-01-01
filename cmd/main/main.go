package main

import (
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/bot"
	"github.com/dane4k/FinMarket/internal/repository"
	"github.com/dane4k/FinMarket/internal/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load("./.env"); err != nil {
		log.Println(err)
	}

	logFile, err := os.OpenFile("./FinMarket.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create or open log file")
	}
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)

	db.InitDB()

	go bot.StartTelegramBot()
	repository.InitTGBot()
	router := gin.Default()
	router.Static("/static", "./internal/web/static")
	router.LoadHTMLGlob("internal/web/templates/*")
	routes.InitializeRoutes(router)

	err = router.Run(":8080")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to start the application")
		return
	}
	logrus.Info("Starting application")
}
