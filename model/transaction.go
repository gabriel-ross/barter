package model

import "time"

type Transaction struct {
	ID         string
	Timestamp  time.Time
	Quantities map[Currency]float64
	Sender     Account
	Recipient  Account
}
