package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"cmd/app/main.go/internal/dto"
	"cmd/app/main.go/internal/model"
	mocks "cmd/app/main.go/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestWalletCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := mocks.NewMockWallet(ctrl)

	router := gin.Default()
	handler := New(router, fakeService)
	handler.Register()

	t.Run("TestWalletCreate_Success", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeService.EXPECT().Create(gomock.Any()).Return(fakeUUID, nil)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/wallets", nil)
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		if recoder.Code != http.StatusCreated {
			t.Errorf("response code incorrect. Expected: %d, received: %d", http.StatusCreated, recoder.Code)
		}

		resp := make(map[string]any)
		correctResp := map[string]any{
			"message": map[string]any{
				"walletId": fakeUUID.String(),
			},
			"success": true,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		message := resp["message"].(map[string]any)
		value, ok := message["walletId"]
		if !ok || value != fakeUUID.String() {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})
	t.Run("TestWalletCreate_ServiceErr", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeErr := fmt.Errorf("db err")
		fakeService.EXPECT().Create(gomock.Any()).Return(fakeUUID, fakeErr)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/wallets", nil)

		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusInternalServerError
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)
		correctResp := map[string]any{
			"message": fmt.Sprint(fakeErr),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})
}

func TestWalletBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := mocks.NewMockWallet(ctrl)

	router := gin.Default()
	handler := New(router, fakeService)
	handler.Register()

	t.Run("TestWalletBalance_Success", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 0,
		}
		fakeService.EXPECT().Balance(gomock.Any(), fakeUUID).Return(fakeWallet, nil)

		url := fmt.Sprintf("/api/v1/wallets/%s", fakeUUID)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusOK
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		message := resp["message"].(map[string]any)
		respWalletId, ok := message["walletId"]
		if !ok || respWalletId != fakeWallet.UUID.String() {
			t.Errorf("response body incorrect. Expected walletId: %v, received: %v", fakeWallet.UUID, respWalletId)
		}

		respBalance, ok := message["balance"]
		if !ok || respBalance != fakeWallet.Balance {
			t.Errorf("response body incorrect. Expected walletId: %v, received: %v", respBalance, fakeWallet.Balance)
		}
	})
	t.Run("TestWalletBalance_BadUUID", func(t *testing.T) {
		url := fmt.Sprintf("/api/v1/wallets/%s", "123")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusBadRequest
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("incorrect wallet uuid"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})

	t.Run("TestWalletBalance_WalletNotFound", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 0,
		}
		fakeService.EXPECT().Balance(gomock.Any(), fakeUUID).Return(fakeWallet, pgx.ErrNoRows)

		url := fmt.Sprintf("/api/v1/wallets/%s", fakeUUID)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusNotFound
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("wallet not found"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})

	t.Run("TestWalletBalance_InternalErr", func(t *testing.T) {
		fakeUUID := uuid.New()
		fakeWallet := model.Wallet{
			UUID:    fakeUUID,
			Balance: 0,
		}
		fakeService.EXPECT().Balance(gomock.Any(), fakeUUID).Return(fakeWallet, fmt.Errorf("db random err"))

		url := fmt.Sprintf("/api/v1/wallets/%s", fakeUUID)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusInternalServerError
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("wallet service err"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})
}

func TestWalletTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fakeService := mocks.NewMockWallet(ctrl)

	router := gin.Default()
	handler := New(router, fakeService)
	handler.Register()

	t.Run("TestWalletTransaction_Success", func(t *testing.T) {
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

		fakeService.EXPECT().Transaction(gomock.Any(), fakeReq).Return(fakeWallet, nil)

		url := fmt.Sprintf("/api/v1/wallet")

		body, err := json.Marshal(fakeReq)

		if err != nil {
			t.Error("marshall err: ", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusOK
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		message := resp["message"].(map[string]any)
		respWalletId, ok := message["walletId"]
		if !ok || respWalletId != fakeWallet.UUID.String() {
			t.Errorf("response body incorrect. Expected walletId: %v, received: %v", fakeWallet.UUID, respWalletId)
		}

		respBalance, ok := message["balance"]
		if !ok || respBalance != fakeWallet.Balance {
			t.Errorf("response body incorrect. Expected walletId: %v, received: %v", respBalance, fakeWallet.Balance)
		}
	})

	t.Run("TestWalletTransaction_ValErr", func(t *testing.T) {

		fakeReq := dto.WalletTransactionRequest{
			Type:   "DEPOSIT",
			Amount: 100,
		}

		url := fmt.Sprintf("/api/v1/wallet")

		body, err := json.Marshal(fakeReq)

		if err != nil {
			t.Error("marshall err: ", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusBadRequest
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("validation err: Key: 'WalletTransactionRequest.UUID' Error:Field validation for 'UUID' failed on the 'required' tag"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})

	t.Run("TestWalletTransaction_WalletNotFound", func(t *testing.T) {
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

		fakeService.EXPECT().Transaction(gomock.Any(), fakeReq).Return(fakeWallet, pgx.ErrNoRows)

		url := fmt.Sprintf("/api/v1/wallet")

		body, err := json.Marshal(fakeReq)

		if err != nil {
			t.Error("marshall err: ", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusNotFound
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("wallet not found"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})

	t.Run("TestWalletTransaction_InternalErr", func(t *testing.T) {
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

		fakeService.EXPECT().Transaction(gomock.Any(), fakeReq).Return(fakeWallet, fmt.Errorf("random db err"))

		url := fmt.Sprintf("/api/v1/wallet")

		body, err := json.Marshal(fakeReq)

		if err != nil {
			t.Error("marshall err: ", err)
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
		if err != nil {
			t.Error("new request err: ", err)
		}

		recoder := httptest.NewRecorder()
		router.ServeHTTP(recoder, req)
		correctCode := http.StatusInternalServerError
		if recoder.Code != correctCode {
			t.Errorf("response code incorrect. Expected: %d, received: %d", correctCode, recoder.Code)
		}

		resp := make(map[string]any)

		correctResp := map[string]any{
			"message": fmt.Sprint("wallet service err"),
			"success": false,
		}

		err = json.Unmarshal(recoder.Body.Bytes(), &resp)
		if err != nil {
			t.Error("unmarshal body err")
		}

		ok := reflect.DeepEqual(correctResp, resp)

		if !ok {
			t.Errorf("response body incorrect. Expected: %v, received: %v", correctResp, resp)
		}
	})
}
