package service

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/dto"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func PostProduct(c *gin.Context) error {
	var request *dto.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error(service_errs.ErrBindingJSON.Error())
		return service_errs.ErrBindingJSON
	}
	product, err := ReqToProduct(request)
	if err != nil {
		logrus.WithError(err).Error(service_errs.ErrInvalidProduct)
		return service_errs.ErrInvalidProduct
	}

	if pgdb.AddProduct(product) != nil {
		logrus.WithError(err).Error(service_errs.ErrAddingProduct)
		return service_errs.ErrAddingProduct
	}
	return nil
}

func ReqToProduct(request *dto.CreateProductRequest) (*entity.Product, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}
	return &entity.Product{
		SellerID:         request.UserID,
		CategoryId:       request.CategoryID,
		SubwayID:         request.SubwayID,
		PhotosURLs:       request.PhotosURLs,
		Name:             request.Name,
		Description:      request.Description,
		ProductCondition: request.ProductCondition,
		Price:            request.Price,
	}, nil
}
