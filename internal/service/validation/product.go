package validation

import (
	"github.com/go-playground/validator/v10"
)

var Conditions = map[string]bool{
	"Новый":      true,
	"Как новый":  true,
	"Хороший":    true,
	"Нормальный": true,
	"Плохой":     true,
}

func ConditionValidator(fl validator.FieldLevel) bool {
	condition := fl.Field().String()
	return Conditions[condition]
}

func IsPriceValid(fl validator.FieldLevel) bool {
	price := fl.Field().Int()

	if price < 0 || price > 100000 {
		return false
	}
	return true
}
