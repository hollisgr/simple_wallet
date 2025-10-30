package handler

import (
	"cmd/app/main.go/internal/dto"
	"cmd/app/main.go/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler interface {
	Register()
}

type handler struct {
	router        *gin.Engine
	walletService service.Wallet
	validator     *validator.Validate
}

func New(r *gin.Engine, ws service.Wallet) Handler {
	return &handler{
		router:        r,
		walletService: ws,
		validator:     validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (h *handler) Register() {
	v1 := h.router.Group("/api/v1")
	v1.POST("/wallet", h.WalletTransaction)
	v1.POST("/wallets", h.WalletCreate)
	v1.GET("/wallets/:uuid", h.WalletBalance)
}

func (h *handler) WalletTransaction(c *gin.Context) {
	req := dto.WalletTransactionRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.sendMsg(c, false, http.StatusBadRequest, "Invalid request JSON")
		return
	}

	err = h.validator.Struct(req)
	if err != nil {
		h.sendMsg(c, false, http.StatusBadRequest, fmt.Sprint("validation err: ", err))
		return
	}

	res, err := h.walletService.Transaction(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.sendMsg(c, false, http.StatusNotFound, "wallet not found")
			return
		}
		h.sendMsg(c, false, http.StatusInternalServerError, "wallet service err")
		return
	}
	h.sendMsg(c, true, http.StatusOK, res)
}

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
		}
		h.sendMsg(c, false, http.StatusInternalServerError, "wallet service err")
		return
	}
	h.sendMsg(c, true, http.StatusOK, res)
}

func (h *handler) sendMsg(c *gin.Context, success bool, status int, message any) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
	})
}
