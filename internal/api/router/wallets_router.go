package router

import (
	"net/http"
	"restapi/internal/api/handlers"
	"restapi/internal/repository/sqlconnect"
)

func WalletRouter() *http.ServeMux {
	mux := http.NewServeMux()
	h := handlers.WalletHandler{
		Repo: sqlconnect.SQLWalletRepo{},
	}
	mux.HandleFunc("GET /api/v1/wallets/{WALLET_UUID}", h.GetOneWalletHandler)
	// TODO: Method POST
	return mux
}
