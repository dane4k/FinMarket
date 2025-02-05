package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService}
}

func (ph *ProductHandler) CreateProductHandler(c *gin.Context) {
	product, err := ph.productService.PostProduct(c)
	if err != nil {
		if errors.Is(err, service.ErrInvalidProduct) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// v product.id budet new id, nuzhen redirect na /product/:id
	c.JSON(http.StatusOK, gin.H{
		"data":       product.ID,
		"Categories": Categories,
	})
}

func (ph *ProductHandler) EditProductHandler(c *gin.Context) {
	err := ph.productService.UpdateProduct(c)
	if err != nil {
		if errors.Is(err, service.ErrBindingJSON) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else if errors.Is(err, service.ErrValidatingProduct) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неправильная цена или состояние"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "товар не существует"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": c.Param("id")})
	return
}

func (ph *ProductHandler) GetProductHandler(c *gin.Context) {
	response, err := ph.productService.GetProductData(c)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"message": "Internal Server Error: " + err.Error(),
		})
		return
	}
	if response == nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"message": "Product not found",
		})
		return
	}
	c.HTML(http.StatusOK, "product.html", gin.H{
		"Categories": Categories,
		"Product":    response,
	})
}

func (ph *ProductHandler) InputProductHandler(c *gin.Context) {
	isAuthed, _ := c.Get("isAuthed")
	IsAuthenticated := isAuthed.(bool)
	if !IsAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	user, _ := c.Get("user")
	usr, _ := user.(*entity.User)

	c.HTML(http.StatusOK, "add.html", gin.H{
		"Categories": Categories,
		"UserID":     usr.TgID,
	})
}

func (ph *ProductHandler) GetProductsHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	category := c.DefaultQuery("category", "-1")

	productsResp, err := ph.productService.GetProducts(page, pageSize, category)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      productsResp.Page,
		"page_size": productsResp.PageSize,
		"products":  productsResp.Products,
	})
}
