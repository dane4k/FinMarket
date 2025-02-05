package main

import (
	"github.com/dane4k/FinMarket/app"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/logger"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/route"
	"github.com/dane4k/FinMarket/internal/service/validation"
	"github.com/dane4k/FinMarket/internal/tg_bot"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		logrus.Fatal(err)
	}
	log.Println(cfg.Database.Name)

	logger.InitLogger()

	app := app.InitApp(cfg)

	bot, err := tg_bot.NewTGBot(cfg, app.AuthRepository, app.UserRepository)
	if err != nil {
		log.Fatal(err)
	}

	go bot.Start()

	pgdb.InitTGBot(cfg)

	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("condition", validation.ConditionValidator)
		v.RegisterValidation("price", validation.IsPriceValid)
	}

	router.Static("/static", "./internal/web/static")
	router.LoadHTMLGlob("internal/web/templates/*")
	route.InitializeRoutes(router, app)

	if err = router.Run(":8080"); err != nil {
		logrus.WithError(err).Error("Failed to start the server")
	} else {
		logrus.Info("Server started on port 8080")
	}
}
