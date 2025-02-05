package pgdb

import (
	"errors"
	"fmt"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductRepository interface {
	AddProduct(product *entity.Product) error
	GetProduct(id int64) (*entity.Product, error)
	UpdateProduct(product *entity.Product) error
	GetProducts(limit int, page int, categoryID int) ([]entity.Product, error)
	GetProductSeller(id int64) (int64, error)
}

type productRepository struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewProductRepository(db *gorm.DB, cfg *config.Config) ProductRepository {
	return &productRepository{db: db, cfg: cfg}
}

func (pr *productRepository) AddProduct(product *entity.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	if err := pr.db.Create(product).Error; err != nil {
		logrus.WithError(err).Error("Failed to add product")
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (pr *productRepository) GetProduct(id int64) (*entity.Product, error) {
	product := &entity.Product{}
	if err := pr.db.Where("id = ?", id).First(product).Error; err != nil {
		logrus.WithError(err).Error("Failed to get product")
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (pr *productRepository) UpdateProduct(product *entity.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}
	if err := pr.db.Save(product).Error; err != nil {
		logrus.WithError(err).Error("Failed to update product")
		return err
	}
	return nil
}

func (pr *productRepository) GetProducts(limit int, page int, categoryID int) ([]entity.Product, error) {
	var products []entity.Product
	query := pr.db.
		Limit(limit).
		Offset((page - 1) * limit).
		Order("date_posted desc")
	if categoryID != -1 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Find(&products).Error; err != nil {
		logrus.WithError(err).Error("Failed to get products")
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

func (pr *productRepository) GetProductSeller(id int64) (int64, error) {
	product, err := pr.GetProduct(id)
	if err != nil {
		return 0, err
	}
	return product.ID, err
}
