package models

import "time"

type Favorite struct {
	UserID    int64     `gorm:"primaryKey;column:user_id"`
	ProductID int64     `gorm:"primaryKey;column:product_id"`
	AddedAt   time.Time `gorm:"column:added_at"`
}
