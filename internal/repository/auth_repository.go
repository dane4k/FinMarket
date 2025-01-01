package repository

import (
	"database/sql"
	"errors"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/default_error"
	"github.com/dane4k/FinMarket/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func GetAuthRecord(authToken string) (*model.AuthRecord, error) {
	var record *model.AuthRecord

	if err := db.DB.Where("token = ?", authToken).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("%s: %s for token %s", default_error.GetTokenContext, default_error.ErrTokenNotFound, authToken)
			return nil, err
		}
		logrus.WithError(err).Errorf("%s: %s", default_error.GetTokenContext, default_error.ErrDatabaseError)
		return nil, err
	}
	return record, nil
}

func saveAuthRecord(record *model.AuthRecord) error {
	if err := db.DB.Save(&record).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", default_error.SaveTokenContext, default_error.ErrDatabaseError)
		return err
	}
	return nil
}

func ConfirmToken(token string, tgId int) error {
	record, err := GetAuthRecord(token)
	if err != nil {
		return err
	}
	record.Status = "confirmed"
	record.TgID = int64(tgId)
	return saveAuthRecord(record)
}

func InvalidateAllTokens(userId int64) error {
	var tokens []sql.NullString
	if err := db.DB.Model(&model.AuthRecord{}).Where("tg_id = ?", userId).Pluck("jwt", &tokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("%ss: %s for user id %v", default_error.GetTokenContext, default_error.ErrTokenNotFound, userId)
			return err
		}
		logrus.WithError(err).Errorf("%ss: %s", default_error.GetTokenContext, default_error.ErrDatabaseError)
	}

	var invalidatedTokens []model.InvalidJWT
	for _, token := range tokens {
		if token.Valid && token.String != "" {
			invalidatedTokens = append(invalidatedTokens, model.InvalidJWT{
				JWTToken: token.String,
			})
		}
	}

	if len(invalidatedTokens) > 0 {
		if err := db.DB.Create(&invalidatedTokens).Error; err != nil {
			logrus.WithError(err).Errorf("%ss: %s", default_error.SaveTokenContext, default_error.ErrDatabaseError)
			return err
		}
		return nil
	}
	if err := db.DB.Model(&model.AuthRecord{}).Where("tg_id = ?", userId).Update("jwt", nil).Error; err != nil {
		logrus.WithError(err).Errorf("%ss: %s", default_error.InvalidateTokensContext, default_error.ErrDatabaseError)
		return err
	}
	return nil
}

func PinJTI(authRecordID uint, jti string) error {
	if err := db.DB.Model(&model.AuthRecord{}).Where("id = ?", authRecordID).Update("jwt", jti).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Errorf("%ss: %s", default_error.UpdateJTIContext, default_error.ErrTokenNotFound)
			return err
		}
		logrus.WithError(err).Errorf("%s: %v", default_error.ErrAddJTI, authRecordID)
		return err
	}
	return nil
}

func IsJTIValid(tokenJTI string) (bool, error) {
	var invalidJWT *model.InvalidJWT

	err := db.DB.Where("jwt_token = ?", tokenJTI).First(&invalidJWT).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		logrus.WithError(err).Errorf("%ss: %s", default_error.GetJTIContext, default_error.ErrDatabaseError)
		return false, err
	}
	return false, nil
}

func CreateAuthRecord(token string) error {
	authRecord := &model.AuthRecord{
		Token:     token,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return saveAuthRecord(authRecord)
}

func InvalidateAuthRecord(jti string) error {
	invalidJWT := model.InvalidJWT{
		JWTToken: jti,
	}
	if err := db.DB.Create(&invalidJWT).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", default_error.InvalidateSessionContext, default_error.ErrDatabaseError)
		return err
	}
	return nil
}
