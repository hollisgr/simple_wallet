package dto

import "github.com/google/uuid"

type WalletTransactionRequest struct {
	UUID   uuid.UUID `json:"valletId" validate:"required,uuid"`
	Type   string    `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount float64   `json:"amount" validate:"required,gte=0.01"`
}
