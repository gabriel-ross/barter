package model

import "time"

type Account struct {
	ID                string             `firestore:"id"`
	UserID            string             `firestore:"user"`
	Funds             map[string]float64 `firestore:"funds"`
	Reputation        int                `firestore:"reputation"`
	CreationTimestamp time.Time          `firestore:"creationTimestamp"`
}

func NewAccount() Account {
	return Account{
		ID:                "",
		UserID:            "",
		Funds:             map[string]float64{},
		Reputation:        100,
		CreationTimestamp: time.Now(),
	}
}
