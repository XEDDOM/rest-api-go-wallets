package sqlconnect

import "restapi/internal/models"

type WalletRepository interface {
	GetWalletByID(id int) (models.Wallet, error)
	UpdateWallet(data map[string]any) error
}
