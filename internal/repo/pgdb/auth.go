package pgdb

import (
	"database/sql"
	"errors"
	"github.com/dane4k/FinMarket/db"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/repoerrs"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func GetAuthRecord(authToken string) (*entity.AuthRecord, error) {
	var record *entity.AuthRecord

	if err := db.DB.Where("token = ?", authToken).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("%s: %s for token %s", repoerrs.GetTokenContext, repoerrs.TokenNotFoundContext, authToken)
			return nil, err
		}
		logrus.WithError(err).Errorf("%s: %s", repoerrs.GetTokenContext, repoerrs.ErrDatabaseError)
		return nil, err
	}
	return record, nil
}

func saveAuthRecord(record *entity.AuthRecord) error {
	if err := db.DB.Save(&record).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", repoerrs.SaveTokenContext, repoerrs.ErrDatabaseError)
		return err
	}
	return nil
}

func ConfirmToken(token string, tgId int64) error {
	record, err := GetAuthRecord(token)
	if err != nil {
		return err
	}
	record.Status = "confirmed"
	record.TgID = tgId
	return saveAuthRecord(record)
}

func InvalidateAllTokens(userId int64) error {
	var tokens []sql.NullString
	if err := db.DB.Model(&entity.AuthRecord{}).Where("tg_id = ?", userId).Pluck("jwt", &tokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Errorf("%ss: %s for user id %v", repoerrs.GetTokenContext, repoerrs.TokenNotFoundContext, userId)
			return err
		}
		logrus.WithError(err).Errorf("%ss: %s", repoerrs.GetTokenContext, repoerrs.ErrDatabaseError)
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
		if err := db.DB.Create(&invalidatedTokens).Error; err != nil {
			logrus.WithError(err).Errorf("%ss: %s", repoerrs.SaveTokenContext, repoerrs.ErrDatabaseError)
			return err
		}
		return nil
	}
	if err := db.DB.Model(&entity.AuthRecord{}).Where("tg_id = ?", userId).Update("jwt", nil).Error; err != nil {
		logrus.WithError(err).Errorf("%ss: %s", repoerrs.InvalidateTokensContext, repoerrs.ErrDatabaseError)
		return err
	}
	return nil
}

func PinJTI(authRecordID uint, jti string) error {
	if err := db.DB.Model(&entity.AuthRecord{}).Where("id = ?", authRecordID).Update("jwt", jti).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.WithError(err).Errorf("%ss: %s", repoerrs.UpdateJTIContext, repoerrs.TokenNotFoundContext)
			return err
		}
		logrus.WithError(err).Errorf("%s: %v", repoerrs.ErrAddJTI, authRecordID)
		return err
	}
	return nil
}

func IsJTIValid(tokenJTI string) (bool, error) {
	var invalidJWT *entity.InvalidJWT

	err := db.DB.Where("jwt_token = ?", tokenJTI).First(&invalidJWT).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		logrus.WithError(err).Errorf("%ss: %s", repoerrs.GetJTIContext, repoerrs.ErrDatabaseError)
		return false, err
	}
	return false, nil
}

func CreateAuthRecord(token string) error {
	authRecord := &entity.AuthRecord{
		Token:     token,
		Status:    "pending",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return saveAuthRecord(authRecord)
}

func InvalidateAuthRecord(jti string) error {
	invalidJWT := entity.InvalidJWT{
		JWTToken: jti,
	}
	if err := db.DB.Create(&invalidJWT).Error; err != nil {
		logrus.WithError(err).Errorf("%s: %s", repoerrs.InvalidateSessionContext, repoerrs.ErrDatabaseError)
		return err
	}
	return nil
}
