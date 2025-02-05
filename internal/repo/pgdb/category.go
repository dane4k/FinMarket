package pgdb

import (
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetAllCategories(categories *[]entity.Category) error
}

type categoryRepository struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewCategoryRepository(db *gorm.DB, cfg *config.Config) CategoryRepository {
	return &categoryRepository{db: db, cfg: cfg}
}

func (cr *categoryRepository) GetAllCategories(categories *[]entity.Category) error {
	if err := cr.db.Find(&categories).Error; err != nil {
		logrus.WithError(err).Fatal("Error getting all categories")
		return err
	}
	return nil
}
