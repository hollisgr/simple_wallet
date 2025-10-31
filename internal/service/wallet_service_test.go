package service

import (
	mocks "cmd/app/main.go/internal/db/mock"
	"cmd/app/main.go/internal/dto"
	"cmd/app/main.go/internal/model"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestWalletServiceCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeDB := mocks.NewMockStorage(ctrl)
	ws := New(fakeDB)

	t.Run("TestWalletServiceCreate_Success", func(t *testing.T) {
		fakeDB.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
		_, err := ws.Create(t.Context())
		if err != nil {
			t.Error("create err")
		}
	})

	t.Run("TestWalletServiceCreate_Fail", func(t *testing.T) {
		fakeDB.EXPECT().Create(gomock.Any(), gomock.Any()).Return(fmt.Errorf("db random err"))
		_, err := ws.Create(t.Context())
		expErr := fmt.Errorf("service create wallet error")
		if err.Error() != expErr.Error() {
			t.Errorf("create err. Expected: %v, recieved: %v", expErr, err)
		}
	})
}

func TestWalletServiceBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeDB := mocks.NewMockStorage(ctrl)
	ws := New(fakeDB)

	t.Run("TestWalletServiceBalance_Success", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeDB.EXPECT().Balance(gomock.Any(), fakeUUID).Return(fakeWallet, nil)
		wallet, err := ws.Balance(t.Context(), fakeUUID)
		if err != nil {
			t.Error("create err")
		}
		if wallet != fakeWallet {
			t.Errorf("Expected: %v, recieved: %v", fakeWallet, wallet)
		}
	})

	t.Run("TestWalletServiceBalance_Fail", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeErr := fmt.Errorf("random db err")
		fakeDB.EXPECT().Balance(gomock.Any(), fakeUUID).Return(fakeWallet, fakeErr)
		_, err := ws.Balance(t.Context(), fakeUUID)
		if err.Error() != fakeErr.Error() {
			t.Errorf("create err. Expected: %v, recieved: %v", fakeErr, err)
		}
	})
}

func TestWalletServiceTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fakeDB := mocks.NewMockStorage(ctrl)
	ws := New(fakeDB)

	t.Run("TestWalletServiceTransactionDeposit_Success", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeReq := dto.WalletTransactionRequest{
			UUID:   fakeUUID,
			Type:   "DEPOSIT",
			Amount: 100,
		}
		fakeDB.EXPECT().Deposit(gomock.Any(), fakeUUID, fakeReq.Amount).Return(fakeWallet, nil)
		wallet, err := ws.Transaction(t.Context(), fakeReq)
		if err != nil {
			t.Error("create err")
		}
		if wallet != fakeWallet {
			t.Errorf("Expected: %v, recieved: %v", fakeWallet, wallet)
		}
	})

	t.Run("TestWalletServiceTransactionDeposit_Fail", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeReq := dto.WalletTransactionRequest{
			UUID:   fakeUUID,
			Type:   "DEPOSIT",
			Amount: 100,
		}
		fakeErr := fmt.Errorf("random db err")
		fakeDB.EXPECT().Deposit(gomock.Any(), fakeUUID, fakeReq.Amount).Return(fakeWallet, fakeErr)
		_, err := ws.Transaction(t.Context(), fakeReq)
		if err.Error() != fakeErr.Error() {
			t.Errorf("Expected: %v, recieved: %v", fakeErr, err)
		}
	})

	t.Run("TestWalletServiceTransactionWithdraw_Success", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeReq := dto.WalletTransactionRequest{
			UUID:   fakeUUID,
			Type:   "WITHDRAW",
			Amount: 100,
		}
		fakeDB.EXPECT().Withdraw(gomock.Any(), fakeUUID, fakeReq.Amount).Return(fakeWallet, nil)
		wallet, err := ws.Transaction(t.Context(), fakeReq)
		if err != nil {
			t.Error("create err")
		}
		if wallet != fakeWallet {
			t.Errorf("Expected: %v, recieved: %v", fakeWallet, wallet)
		}
	})

	t.Run("TestWalletServiceTransactionWithdraw_Fail", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 100,
		}
		fakeReq := dto.WalletTransactionRequest{
			UUID:   fakeUUID,
			Type:   "WITHDRAW",
			Amount: 100,
		}
		fakeErr := fmt.Errorf("random db err")
		fakeDB.EXPECT().Withdraw(gomock.Any(), fakeUUID, fakeReq.Amount).Return(fakeWallet, fakeErr)
		_, err := ws.Transaction(t.Context(), fakeReq)
		if err.Error() != fakeErr.Error() {
			t.Errorf("Expected: %v, recieved: %v", fakeErr, err)
		}
	})
}
