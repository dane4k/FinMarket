package handler

import (
	"encoding/json"
	"fmt"
	"github.com/dane4k/FinMarket/internal/dto/product"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (ch *CategoryHandler) GetProductsByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")
	query := fmt.Sprintf("http://localhost:8080/api/products?page=1&page_size=10&category=%s", categoryID)
	response, err := http.Get(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err},
		)
	}
	defer response.Body.Close()

	var resp product.GetProductsResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err},
		)
	}

	products := resp.Products

	isAuthed, _ := c.Get("isAuthed")
	IsAuthenticated := isAuthed.(bool)

	c.HTML(http.StatusOK, "products.html", gin.H{
		"IsAuthenticated": IsAuthenticated,
		"Products":        products,
		"Categories":      Categories,
	})
}
