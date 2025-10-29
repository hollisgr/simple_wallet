package service

import (
	"context"
	"fmt"
	"log"

	"cmd/app/main.go/internal/db"

	"github.com/google/uuid"
)

type Wallet interface {
	Create(ctx context.Context) (uuid.UUID, error)
	Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error)
	Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error)
	Balance(ctx context.Context, uuid uuid.UUID) (float64, error)
}

type wallet struct {
	storage db.Storage
}

func New(s db.Storage) Wallet {
	return &wallet{
		storage: s,
	}
}

func (ws *wallet) Create(ctx context.Context) (uuid.UUID, error) {
	uuid := uuid.New()
	err := ws.storage.Create(ctx, uuid)
	if err != nil {
		log.Println(err)
		return uuid, fmt.Errorf("service create wallet error")
	}
	return uuid, nil
}

func (ws *wallet) Deposit(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
func (ws *wallet) Withdraw(ctx context.Context, uuid uuid.UUID, amount float64) (float64, error) {
	var res float64
	return res, nil
}
func (ws *wallet) Balance(ctx context.Context, uuid uuid.UUID) (float64, error) {
	res, err := ws.storage.Balance(ctx, uuid)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}
