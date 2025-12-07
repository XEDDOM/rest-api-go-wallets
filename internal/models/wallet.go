package models

type Wallet struct {
	ValletID      int    `json:"valletId,omitempty" db:"valletid,omitepmty"`
	OperationType string `json:"operationType,omitempty"`
	Amount        int    `json:"amount,omitempty" db:"amount,omitepmty"`
}
