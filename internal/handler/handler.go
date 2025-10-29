package handler

import (
	"cmd/app/main.go/internal/service"
	"fmt"
	"net/http"

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
	main.POST("/wallets", h.WalletCreate)
	main.GET("/wallets/:uuid", h.WalletBalance)
}

func (h *handler) WalletTransaction(c *gin.Context) {}

func (h *handler) WalletCreate(c *gin.Context) {
	uuid, err := h.walletService.Create(c.Request.Context())
	if err != nil {
		h.sendMsg(c, false, http.StatusInternalServerError, fmt.Sprint(err))
	}
	h.sendMsg(c, true, http.StatusOK, gin.H{
		"valletId": uuid,
	})

}

func (h *handler) WalletBalance(c *gin.Context) {}

func (h *handler) sendMsg(c *gin.Context, success bool, status int, message any) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
	})
}
