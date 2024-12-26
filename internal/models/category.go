package models

type Category struct {
	ID   int64  `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name"`
}
