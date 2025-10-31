package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mocks "cmd/app/main.go/internal/service/mock"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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
