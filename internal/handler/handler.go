package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register()
}

type handler struct {
	router *gin.Engine
}

func New(r *gin.Engine) Handler {
	return &handler{
		router: r,
	}
}

func (h *handler) Register() {
	h.router.GET("/hello", h.Hello)
}

func (h *handler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello, world",
		"success": true,
	})
}
