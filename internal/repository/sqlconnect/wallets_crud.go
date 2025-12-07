package sqlconnect

import (
	"database/sql"
	"errors"
	"restapi/internal/models"
)

type SQLWalletRepo struct{}

func (SQLWalletRepo) GetWalletByID(id int) (models.Wallet, error) {
	return GetWalletByID(id)
}

func GetWalletByID(id int) (models.Wallet, error) {
	db, err := ConnectDb()
	if err != nil {
		return models.Wallet{}, err
	}
	defer db.Close()
	var wallet models.Wallet
	err = db.QueryRow("SELECT valletid, amount FROM wallets WHERE valletid = $1", id).Scan(&wallet.ValletID, &wallet.Amount)
	if err == sql.ErrNoRows {
		return models.Wallet{}, errors.New("wallet not found")
	} else if err != nil {
		return models.Wallet{}, err
	}
	return wallet, nil
}

// TODO: CRUD for wallets POST
