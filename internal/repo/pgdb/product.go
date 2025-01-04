package pgdb

import (
	"errors"
	"fmt"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/entity"
)

func AddProduct(product *entity.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	if err := db.DB.Create(product).Error; err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}
