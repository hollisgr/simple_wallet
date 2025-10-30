package model

import "github.com/google/uuid"

type Wallet struct {
	UUID    uuid.UUID `json:"walletId"`
	Balance float64   `json:"balance"`
}
