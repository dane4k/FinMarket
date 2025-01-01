package model

type InvalidJWT struct {
	ID       uint   `gorm:"primaryKey"`
	JWTToken string `gorm:"column:jwt_token"`
}
