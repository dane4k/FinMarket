package entity

import (
	"github.com/lib/pq"
	"time"
)

type Product struct {
	ID               int64          `gorm:"primaryKey;column:id"`
	SellerID         int64          `gorm:"column:seller_id"`
	CategoryId       int64          `gorm:"column:category_id"`
	SubwayID         int64          `gorm:"column:subway_id"`
	PhotosURLs       pq.StringArray `gorm:"type:text[];column:photos_urls"`
	Name             string         `gorm:"column:name"`
	Description      string         `gorm:"column:description"`
	NViews           int64          `gorm:"column:n_views"`
	ProductCondition string         `gorm:"column:product_condition"`
	DatePosted       time.Time      `gorm:"column:date_posted"`
	Price            int64          `gorm:"column:price"`
	Active           bool           `gorm:"column:active"`
}
