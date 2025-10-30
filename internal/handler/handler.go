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

// Register configures HTTP routes for managing wallet resources.
func (h *handler) Register() {
	v1 := h.router.Group("/api/v1")
	v1.POST("/wallet", h.WalletTransaction)
	v1.POST("/wallets", h.WalletCreate)
	v1.GET("/wallets/:uuid", h.WalletBalance)
}

// WalletTransaction processes incoming requests to perform financial transactions on wallets.
// It first binds and validates the request payload, ensuring proper input structure.
// Then, it delegates the actual transaction processing to the wallet service layer.
// Upon completion, it either returns the result or an appropriate error code if something goes wrong.
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

// WalletCreate handles HTTP POST requests to create a new wallet.
// It calls the wallet service to generate a unique identifier for the newly created wallet
// and responds with this ID along with a success message. If an error occurs during the process,
// it sends back an appropriate error response.
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

// WalletBalance fetches the current balance of a wallet identified by its UUID.
// It extracts the UUID from the request parameters, validates it, then queries the wallet service
// to retrieve the corresponding balance. On successful execution, it returns the balance details.
// In case of errors such as invalid UUID format or missing wallet entry, appropriate error responses are sent.
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

// sendMsg sends a JSON response containing a success indicator and additional message or data.
func (h *handler) sendMsg(c *gin.Context, success bool, status int, message any) {
	c.JSON(status, gin.H{
		"success": success,
		"message": message,
	})
}
