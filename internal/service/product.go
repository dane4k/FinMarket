package service

import (
	"errors"
	"fmt"
	"github.com/dane4k/FinMarket/internal/config"
	product2 "github.com/dane4k/FinMarket/internal/dto/product"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/imgur"
	"github.com/dane4k/FinMarket/internal/repo/pgdb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
)

type ProductService struct {
	cfg  *config.Config
	repo pgdb.ProductRepository
}

func NewProductService(repo pgdb.ProductRepository, cfg *config.Config) *ProductService {
	return &ProductService{repo: repo, cfg: cfg}
}

func (ps *ProductService) PostProduct(c *gin.Context) (*entity.Product, error) {
	var request *product2.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error(ErrBindingJSON.Error())
		return nil, ErrBindingJSON
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return nil, ErrValidatingProduct
	}

	product, err := ps.CreateReqToEntity(request)
	if err != nil {
		logrus.WithError(err).Error(ErrInvalidProduct)
		return nil, ErrInvalidProduct
	}

	if ps.repo.AddProduct(product) != nil {
		logrus.WithError(err).Error(ErrAddingProduct)
		return nil, ErrAddingProduct
	}
	return product, nil
}

func (ps *ProductService) UpdateProduct(c *gin.Context) error {
	var request *product2.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logrus.WithError(err).Error(ErrBindingJSON.Error())
		return ErrBindingJSON
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return ErrValidatingProduct
	}

	product, err := ps.repo.GetProduct(request.ProductID)
	if err != nil {
		return ErrInvalidProduct
	}

	product.CategoryId = request.CategoryID
	product.SubwayID = request.SubwayID
	product.PhotosURLs = request.PhotosURLs
	product.Name = request.Name
	product.Description = request.Description
	product.ProductCondition = request.ProductCondition
	product.Price = request.Price

	if err = ps.repo.UpdateProduct(product); err != nil {
		logrus.WithError(err).Error("error updating product")
		return ErrUpdatingProduct
	}
	return nil
}

func (ps *ProductService) GetProductData(c *gin.Context) (*product2.GetProductResponse, error) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		logrus.WithError(err).Error(ErrInvalidID)
		return nil, ErrInvalidID
	}
	product, err := ps.repo.GetProduct(productID)
	if err != nil {
		return nil, ErrProductNotFound
	}

	return &product2.GetProductResponse{
		ID:               productID,
		SellerID:         product.SellerID,
		CategoryID:       product.CategoryId,
		SubwayID:         product.SubwayID,
		PhotosURLs:       product.PhotosURLs,
		Name:             product.Name,
		Description:      product.Description,
		NViews:           product.NViews,
		ProductCondition: product.ProductCondition,
		DatePosted:       product.DatePosted,
		Price:            product.Price,
		Active:           product.Active,
	}, nil
}

func (ps *ProductService) CreateReqToEntity(request *product2.CreateProductRequest) (*entity.Product, error) {
	if request == nil {
		return nil, errors.New("request is nil")
	}

	description := request.Description
	if description == "" {
		description = "Нет описания"
	}

	photosURLs := make([]string, len(request.PhotosBytes))

	for i, photo := range request.PhotosBytes {
		url, err := imgur.UploadImageToImgur(photo, ps.cfg.Imgur.AccessToken)
		if err != nil {
			return nil, err
		}
		photosURLs[i] = url
	}

	return &entity.Product{
		SellerID:         request.UserID,
		CategoryId:       request.CategoryID,
		SubwayID:         request.SubwayID,
		PhotosURLs:       photosURLs,
		Name:             request.Name,
		Description:      description,
		ProductCondition: request.ProductCondition,
		Price:            request.Price,
		DatePosted:       time.Now(),
		Active:           true,
	}, nil
}

func (ps *ProductService) GetProducts(pageNumber string, pageSize string, categoryID string) (*product2.GetProductsResponse, error) {
	pageNum, err := strconv.Atoi(pageNumber)
	if err != nil {
		logrus.WithError(err).Error("invalid query parameter")
		return nil, fmt.Errorf("invalid query parameter: %s", pageNumber)
	}
	pageSz, err := strconv.Atoi(pageSize)
	if err != nil {
		logrus.WithError(err).Error("invalid query parameter")
		return nil, fmt.Errorf("invalid query parameter: %s", pageNumber)
	}

	categoryIDInt, err := strconv.Atoi(categoryID)
	if err != nil {
		logrus.WithError(err).Error("invalid query parameter")
		return nil, fmt.Errorf("invalid query parameter: %s", pageNumber)
	}

	products, err := ps.repo.GetProducts(pageSz, pageNum, categoryIDInt)
	if err != nil {
		return nil, err
	}
	productsArr := make([]product2.ProdPreviewResponse, len(products))

	for i, product := range products {
		productsArr[i] = product2.ProdPreviewResponse{
			ID:       product.ID,
			PhotoURI: product.PhotosURLs[0],
			Name:     product.Name,
			Price:    product.Price,
		}
	}
	log.Println(len(products))
	return &product2.GetProductsResponse{
		Page:     pageNum,
		PageSize: pageSz,
		Products: productsArr,
	}, nil
}

func (ps *ProductService) GetSeller(productID int64) (int64, error) {
	return ps.repo.GetProductSeller(productID)
}
