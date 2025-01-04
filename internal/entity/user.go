package entity

import (
	"time"
)

type User struct {
	TgID       int64     `gorm:"primaryKey;column:tg_id"`
	Name       string    `gorm:"column:name"`
	Rating     float64   `gorm:"column:rating"`
	TgUsername string    `gorm:"column:tg_username"`
	AvatarURI  string    `gorm:"column:avatar_uri"`
	Banned     bool      `gorm:"column:banned"`
	RegDate    time.Time `gorm:"column:reg_date"`
}
