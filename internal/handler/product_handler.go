package handler

import (
	"errors"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/dane4k/FinMarket/internal/service/service_errs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateProductHandler(c *gin.Context) {
	err := service.PostProduct(c)
	if err != nil {
		if errors.Is(err, service_errs.ErrInvalidProduct) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "added"})
}
