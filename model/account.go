package model

import "time"

type Account struct {
	ID                string
	User              User
	Funds             map[Currency]float64
	TransactionVolume int
	CreationDate      time.Time
}
