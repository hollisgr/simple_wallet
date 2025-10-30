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

func (ws *wallet) Create(ctx context.Context) (uuid.UUID, error) {
	uuid := uuid.New()
	err := ws.storage.Create(ctx, uuid)
	if err != nil {
		log.Println(err)
		return uuid, fmt.Errorf("service create wallet error")
	}
	return uuid, nil
}

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

func (ws *wallet) Balance(ctx context.Context, uuid uuid.UUID) (model.Wallet, error) {
	res, err := ws.storage.Balance(ctx, uuid)
	if err != nil {
		log.Println(err)
		return res, err
	}
	return res, nil
}
