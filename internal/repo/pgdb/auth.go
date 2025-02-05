package pgdb

import (
	"database/sql"
	"errors"
	"github.com/dane4k/FinMarket/internal/config"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type AuthRepository interface {
	GetAuthRecord(token string) (*entity.AuthRecord, error)
	SaveAuthRecord(record *entity.AuthRecord) error
	ConfirmToken(token string, tgId int64) error
	InvalidateAllTokens(userId int64) error
	PinJTI(authRecordID uint, jti string) error
	IsJTIValid(tokenJTI string) (bool, error)
	CreateAuthRecord(token string) error
	InvalidateAuthRecord(jti string) error
}

type authRepository struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewAuthRepository(db *gorm.DB, cfg *config.Config) AuthRepository {
	return &authRepository{db: db, cfg: cfg}
}

func (ar *authRepository) GetAuthRecord(token string) (*entity.AuthRecord, error) {
	var record *entity.AuthRecord

	if err := ar.db.Where("token = ?", token).First(&record).Error; err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return record, nil
}

func (ar *authRepository) SaveAuthRecord(record *entity.AuthRecord) error {
	if err := ar.db.Save(&record).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (ar *authRepository) ConfirmToken(token string, tgId int64) error {
	record, err := ar.GetAuthRecord(token)
	if err != nil {
		return err
	}
	record.Status = "confirmed"
	record.TgID = tgId
	return ar.SaveAuthRecord(record)
}

func (ar *authRepository) InvalidateAllTokens(userId int64) error {
	var tokens []sql.NullString
	if err := ar.db.Model(&entity.AuthRecord{}).Where("tg_id = ?", userId).Pluck("jwt", &tokens).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}

	var invalidatedTokens []entity.InvalidJWT
	for _, token := range tokens {
		if token.Valid && token.String != "" {
			invalidatedTokens = append(invalidatedTokens, entity.InvalidJWT{
				JWTToken: token.String,
			})
		}
	}

	if len(invalidatedTokens) > 0 {
		return ar.db.Create(&invalidatedTokens).Error
	}
	return ar.db.Model(&entity.AuthRecord{}).Where("tg_id = ?", userId).Update("jwt", nil).Error
}

func (ar *authRepository) PinJTI(authRecordID uint, JTI string) error {
	if err := ar.db.Model(&entity.AuthRecord{}).Where("id = ?", authRecordID).Update("jwt", JTI).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %v", repo.ErrAddJTI, authRecordID)
		return err
	}
	return nil
}

func (ar *authRepository) IsJTIValid(JTI string) (bool, error) {
	var invalidJWT *entity.InvalidJWT

	err := ar.db.Where("jwt_token = ?", JTI).First(&invalidJWT).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (ar *authRepository) CreateAuthRecord(token string) error {
	authRecord := &entity.AuthRecord{
		Token:     token,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Second),
	}

	return ar.SaveAuthRecord(authRecord)
}

func (ar *authRepository) InvalidateAuthRecord(JTI string) error {
	invalidJWT := entity.InvalidJWT{
		JWTToken: JTI,
	}
	return ar.db.Create(&invalidJWT).Error
}
