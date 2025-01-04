package dto

type CreateProductRequest struct {
	UserID           int64    `json:"user_id" binding:"required"`
	CategoryID       int64    `json:"category_id" binding:"required"`
	SubwayID         int64    `json:"subway_id" binding:"required"`
	PhotosURLs       []string `json:"photos_urls" binding:"required"`
	Name             string   `json:"name" binding:"required"`
	Description      string   `json:"description" binding:"required"`
	ProductCondition string   `json:"product_condition" binding:"required"`
	Price            int64    `json:"price" binding:"required"`
}
