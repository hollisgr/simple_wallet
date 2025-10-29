package dto

import "github.com/google/uuid"

type WalletFromDB struct {
	UUID    uuid.UUID `db:"uuid"`
	Balance float64   `db:"balance"`
}

type WalletToWeb struct {
	UUID    uuid.UUID `json:"valletId"`
	Balance float64   `json:"balance"`
}

type WalletTransactionRequest struct {
	UUID   uuid.UUID `json:"valletId" validate:"required,uuid"`
	Type   string    `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount float64   `json:"amount" validate:"required,gte=0.01"`
}
