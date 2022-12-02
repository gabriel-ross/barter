package model

import "time"

type Account struct {
	ID                string             `firestore:"id"`
	Owner             string             `firestore:"owner"`
	Funds             map[string]float64 `firestore:"funds"`
	Reputation        int                `firestore:"reputation"`
	CreationTimestamp time.Time          `firestore:"creationTimestamp"`
}
