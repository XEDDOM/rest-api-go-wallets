package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"restapi/internal/models"
	"testing"
)

type MockRepo struct {
	GetWalletByIDFunc func(int) (models.Wallet, error)
	UpdateWalletFunc  func(map[string]any) error
}

func (m MockRepo) GetWalletByID(id int) (models.Wallet, error) {
	return m.GetWalletByIDFunc(id)
}

func (m MockRepo) UpdateWallet(data map[string]any) error {
	return m.UpdateWalletFunc(data)
}

func TestGetOneWalletHandler_OK(t *testing.T) {
	mock := MockRepo{
		GetWalletByIDFunc: func(id int) (models.Wallet, error) {
			return models.Wallet{
				ValletID: 1,
				Amount:   500,
			}, nil
		},
	}
	h := WalletHandler{Repo: mock}
	req := httptest.NewRequest("GET", "/api/v1/wallets/1", nil)
	req.SetPathValue("WALLET_UUID", "1")
	rr := httptest.NewRecorder()
	h.GetOneWalletHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var w models.Wallet
	if err := json.Unmarshal(rr.Body.Bytes(), &w); err != nil {
		t.Fatal("cannot decode response JSON")
	}
	if w.Amount != 500 {
		t.Fatalf("expected amount 500, got %d", w.Amount)
	}
}

func TestGetOneWalletHandler_InvalidID(t *testing.T) {
	h := WalletHandler{Repo: MockRepo{}}
	req := httptest.NewRequest("GET", "/api/v1/wallets/abc", nil)
	req.SetPathValue("WALLET_UUID", "abc")
	rr := httptest.NewRecorder()
	h.GetOneWalletHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetOneWalletHandler_DBError(t *testing.T) {
	mock := MockRepo{
		GetWalletByIDFunc: func(id int) (models.Wallet, error) {
			return models.Wallet{}, errors.New("db error")
		},
	}
	h := WalletHandler{Repo: mock}
	req := httptest.NewRequest("GET", "/api/v1/wallets/1", nil)
	req.SetPathValue("WALLET_UUID", "1")
	rr := httptest.NewRecorder()
	h.GetOneWalletHandler(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestUpdateWalletHandler_OK(t *testing.T) {
	mock := MockRepo{
		UpdateWalletFunc: func(data map[string]any) error {
			return nil
		},
	}
	h := WalletHandler{Repo: mock}
	body := []byte(`{"valletId":1,"operationType":"DEPOSIT","amount":100}`)
	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	h.UpdateWalletHandler(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
}

func TestUpdateWalletHandler_InvalidJSON(t *testing.T) {
	h := WalletHandler{Repo: MockRepo{}}
	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer([]byte("INVALID")))
	rr := httptest.NewRecorder()
	h.UpdateWalletHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateWalletHandler_DBError(t *testing.T) {
	mock := MockRepo{
		UpdateWalletFunc: func(data map[string]any) error {
			return errors.New("update error")
		},
	}
	h := WalletHandler{Repo: mock}
	body := []byte(`{"valletId":1,"operationType":"DEPOSIT","amount":100}`)
	req := httptest.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	h.UpdateWalletHandler(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
