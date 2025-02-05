package product

type CreateProductRequest struct {
	UserID           int64    `json:"user_id" binding:"required"`
	CategoryID       int64    `json:"category_id" binding:"required"`
	SubwayID         int64    `json:"subway_id" binding:"required"`
	PhotosBytes      [][]byte `json:"photos_bytes" binding:"required"`
	Name             string   `json:"name" binding:"required" validate:"min=4,max=90"`
	Description      string   `json:"description" validate:"max=1000"`
	ProductCondition string   `json:"product_condition" binding:"required,condition"`
	Price            int64    `json:"price" binding:"required,price"`
}
