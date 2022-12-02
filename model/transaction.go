package model

import "time"

type Transaction struct {
	ID                 string             `firestore:"id"`
	Quantities         map[string]float64 `firestore:"quantities"` // The keys are currency IDs
	SenderAccountID    string             `firestore:"sender"`     // Sender account ID
	RecipientAccountID string             `firestore:"recipient"`  // Recipient Account ID
	Timestamp          time.Time          `firestore:"timestamp"`
}
