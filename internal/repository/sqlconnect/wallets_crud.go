package sqlconnect

import (
	"database/sql"
	"errors"
	"restapi/internal/models"
	"strconv"
)

type SQLWalletRepo struct{}

func (SQLWalletRepo) GetWalletByID(id int) (models.Wallet, error) {
	return GetWalletByID(id)
}

func (SQLWalletRepo) UpdateWallet(data map[string]any) error {
	return UpdateWallet(data)
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

func UpdateWallet(updates map[string]any) error {
	db, err := ConnectDb()
	if err != nil {
		return err
	}
	defer db.Close()
	walletID, ok := updates["valletId"]
	if !ok || walletID == nil {
		return errors.New("valletid is required")
	}
	opTypeVal, ok := updates["operationType"]
	if !ok || opTypeVal == nil {
		return errors.New("operationType is required")
	}
	operationType, ok := opTypeVal.(string)
	if !ok {
		return errors.New("operationType must be a string")
	}
	amountVal, ok := updates["amount"]
	if !ok || amountVal == nil {
		return errors.New("amount is required")
	}
	var amount int
	switch v := amountVal.(type) {
	case int:
		amount = v
	case float64:
		amount = int(v)
	case string:
		amountInt, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("amount must be a valid number")
		}
		amount = amountInt
	default:
		return errors.New("amount must be a number")
	}
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	var wallet models.Wallet
	err = db.QueryRow("SELECT valletid, amount FROM wallets WHERE valletid = $1", walletID).Scan(&wallet.ValletID, &wallet.Amount)
	if err == sql.ErrNoRows {
		return errors.New("wallet not found")
	} else if err != nil {
		return err
	}
	switch operationType {
	case "DEPOSIT":
		wallet.Amount += amount
	case "WITHDRAW":
		if wallet.Amount < amount {
			return errors.New("insufficient funds")
		}
		wallet.Amount -= amount
	default:
		return errors.New("unknown operationType. Must be DEPOSIT or WITHDRAW")
	}
	_, err = db.Exec("UPDATE wallets SET amount = $1 WHERE valletid = $2",
		wallet.Amount, wallet.ValletID)
	if err != nil {
		return err
	}
	return nil
}
