package entity

import "time"

type Review struct {
	ID       int64     `gorm:"primaryKey;column:id"`
	UserID   int64     `gorm:"column:user_id"`
	SellerID int64     `gorm:"column:seller_id"`
	Rating   int64     `gorm:"column:rating"`
	Comment  string    `gorm:"column:comment"`
	Date     time.Time `gorm:"column:date"`
}
