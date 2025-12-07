package handlers

import (
	"encoding/json"
	"net/http"
	"restapi/internal/repository/sqlconnect"
	"strconv"
)

type WalletHandler struct {
	Repo sqlconnect.WalletRepository
}

func (h WalletHandler) GetOneWalletHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("WALLET_UUID")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	wallet, err := h.Repo.GetWalletByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(wallet)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// TODO: Method POST
