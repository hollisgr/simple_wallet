package service

import (
	"context"
	"fmt"
	"log"

	"cmd/app/main.go/internal/db"
	"cmd/app/main.go/internal/dto"
	"cmd/app/main.go/internal/model"

	"github.com/google/uuid"
)

type Wallet interface {
	Create(ctx context.Context) (uuid.UUID, error)
	Transaction(ctx context.Context, req dto.WalletTransactionRequest) (model.Wallet, error)
	Balance(ctx context.Context, uuid uuid.UUID) (model.Wallet, error)
}

type wallet struct {
	storage db.Storage
}

func New(s db.Storage) Wallet {
	return &wallet{
		storage: s,
	}
}

// Create generates a new wallet with a unique identifier and saves it to persistent storage.
// It generates a UUID, attempts to store the wallet, and returns the UUID upon success or an error otherwise.
func (ws *wallet) Create(ctx context.Context) (uuid.UUID, error) {
	uuid := uuid.New()
	err := ws.storage.Create(ctx, uuid)
	if err != nil {
		log.Println(err)
		return uuid, fmt.Errorf("service create wallet error")
	}
	return uuid, nil
}

// Transaction performs deposit or withdrawal operations on a wallet based on the request type.
// It determines the action (deposit or withdraw) and delegates the task to the storage layer accordingly.
// Any errors encountered during the process are logged and returned.
func (ws *wallet) Transaction(ctx context.Context, req dto.WalletTransactionRequest) (model.Wallet, error) {
	var res model.Wallet
	var err error

	switch req.Type {
	case "DEPOSIT":
		res, err = ws.storage.Deposit(ctx, req.UUID, req.Amount)

	case "WITHDRAW":
		res, err = ws.storage.Withdraw(ctx, req.UUID, req.Amount)
	}
	if err != nil {
		log.Println("wallet service transaction err: ", err)
		return res, err
	}
	return res, nil
}

// Balance retrieves the current balance of a wallet identified by its UUID.
// It forwards the request to the storage layer and returns the retrieved balance or an error.
func (ws *wallet) Balance(ctx context.Context, uuid uuid.UUID) (model.Wallet, error) {
	res, err := ws.storage.Balance(ctx, uuid)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}
