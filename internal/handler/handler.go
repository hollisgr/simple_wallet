package handler

import (
	"cmd/app/main.go/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register()
}

type handler struct {
	router        *gin.Engine
	walletService service.Wallet
}

func New(r *gin.Engine, ws service.Wallet) Handler {
	return &handler{
		router:        r,
		walletService: ws,
	}
}

func (h *handler) Register() {
	main := h.router.Group("/api/v1")
	main.POST("/wallet", h.WalletTransaction)
	main.GET("/wallets/:uuid", h.WalletBalance)
}

func (h *handler) WalletTransaction(c *gin.Context) {}

func (h *handler) WalletBalance(c *gin.Context) {}
