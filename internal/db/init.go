package db

import (
	"fmt"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) *gorm.DB {

	connectionStr := fmt.Sprintf(
		"host=%s dbname=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(connectionStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		logrus.Fatal(err)
	}

	return db
}
