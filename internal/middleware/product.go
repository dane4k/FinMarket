package middleware

import (
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/dane4k/FinMarket/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProductMiddleware struct {
	productService *service.ProductService
}

func NewProductMiddleware(productService *service.ProductService) *ProductMiddleware {
	return &ProductMiddleware{productService}
}

func (pmw *ProductMiddleware) ProductBelongsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuthed, _ := c.Get("isAuthed")
		if !isAuthed.(bool) {
			c.Set("belongs", false)
			c.Next()
			return
		}
		user, _ := c.Get("user")
		usr, ok := user.(entity.User)
		if !ok {
			c.Set("belongs", false)
			c.Next()
			return
		}
		productID, err := strconv.ParseInt(c.Param("productID"), 10, 64)
		if err != nil {
			c.Set("belongs", false)
			c.Next()
			return
		}
		userID, err := pmw.productService.GetSeller(productID)
		if err != nil {
			c.Set("belongs", false)
			c.Next()
			return
		}
		if usr.TgID == userID {
			c.Set("belongs", true)
			c.Next()
			return
		}
		c.Set("belongs", false)
		c.Next()
	}
}
