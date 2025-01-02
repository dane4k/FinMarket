package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadHomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "home page",
	})
}
