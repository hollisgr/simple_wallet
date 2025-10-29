package handler

import (
	"cmd/app/main.go/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return
	}
	data := map[string]any{
		"walletId": uuid,
	}
	h.sendMsg(c, true, http.StatusOK, data)

}

func (h *handler) WalletBalance(c *gin.Context) {
	uuidStr := c.Params.ByName("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		h.sendMsg(c, false, http.StatusBadRequest, "incorrect wallet uuid")
		return
	}
	res, err := h.walletService.Balance(c.Request.Context(), uuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.sendMsg(c, false, http.StatusNotFound, "wallet not found")
			return
		} else {
			h.sendMsg(c, false, http.StatusInternalServerError, "wallet service err")
			return
		}
	}
	data := map[string]any{
		"balance":  res,
		"walletId": uuid,
	}
	h.sendMsg(c, true, http.StatusOK, data)
}

func (h *handler) sendMsg(c *gin.Context, success bool, status int, message any) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
	})
}
