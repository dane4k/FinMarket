package entity

type Category struct {
	ID   int64  `gorm:"primaryKey;column:id"`
	Name string `gorm:"column:name"`
	URI  string `gorm:"column:uri"`
}
