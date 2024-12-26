package models

import "time"

type AuthRecord struct {
	ID        uint      `gorm:"primaryKey"`
	TgID      int64     `gorm:"column:tg_id"`
	Token     string    `gorm:"unique;not null"`
	JWT       string    `gorm:"column:jwt"`
	Status    string    `gorm:"column:status"`
	CreatedAt time.Time `gorm:"column:created_at"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}
