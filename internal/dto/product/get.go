package product

import (
	"time"
)

type GetProductResponse struct {
	ID               int64     `json:"id"`
	SellerID         int64     `json:"seller_id"`
	CategoryID       int64     `json:"category_id"`
	SubwayID         int64     `json:"subway_id"`
	PhotosURLs       []string  `json:"photos_urls"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	NViews           int64     `json:"n_views"`
	ProductCondition string    `json:"product_condition"`
	DatePosted       time.Time `json:"date_posted"`
	Price            int64     `json:"price"`
	Active           bool      `json:"active"`
}

type GetProductsResponse struct {
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
	Products []ProdPreviewResponse `json:"products"`
}

type ProdPreviewResponse struct {
	ID       int64  `json:"id"`
	PhotoURI string `json:"photo_uri"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
}
