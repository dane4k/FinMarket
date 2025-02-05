package handler

import (
	"encoding/json"
	"github.com/dane4k/FinMarket/internal/dto/product"
	"github.com/dane4k/FinMarket/internal/entity"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var Categories []entity.Category

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (hh *HomeHandler) LoadHomeHandler(c *gin.Context) {
	response, err := http.Get("http://localhost:8080/api/products?page=1&page_size=10&category=-1")
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err},
		)
	}
	defer response.Body.Close()

	var resp product.GetProductsResponse
	if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err},
		)
	}

	products := resp.Products

	isAuthed, _ := c.Get("isAuthed")
	IsAuthenticated := isAuthed.(bool)

	log.Println(products)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"IsAuthenticated": IsAuthenticated,
		"Products":        products,
		"Categories":      Categories,
	})
}
