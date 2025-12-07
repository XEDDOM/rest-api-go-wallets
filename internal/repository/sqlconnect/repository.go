package sqlconnect

import "restapi/internal/models"

type WalletRepository interface {
	GetWalletByID(id int) (models.Wallet, error)
	// TODO: CRUD for wallets POST
}
