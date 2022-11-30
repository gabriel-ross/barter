package model

import "time"

type Account struct {
	ID                string             `firestore:"id"`
	UserID            string             `firestore:"user"`
	Funds             map[string]float64 `firestore:"funds"`
	TransactionVolume int                `firestore:"transactionVolume"`
	CreationTimestamp time.Time          `firestore:"creationTimestamp"`
}

func NewAccount() Account {
	return Account{
		ID:                "",
		UserID:            "",
		Funds:             map[string]float64{},
		TransactionVolume: 0,
		CreationTimestamp: time.Now(),
	}
}
