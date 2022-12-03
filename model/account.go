package model

import "time"

type Account struct {
	ID                string             `firestore:"id"`
	Owner             string             `firestore:"owner"`
	Balances          map[string]float64 `firestore:"balances"`
	Reputation        int                `firestore:"reputation"`
	CreationTimestamp time.Time          `firestore:"creationTimestamp"`
}
