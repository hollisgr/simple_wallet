package model

import "github.com/google/uuid"

type Wallet struct {
	UUID    uuid.UUID
	Balance float64
}
