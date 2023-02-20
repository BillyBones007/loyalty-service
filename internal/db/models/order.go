package models

import (
	"time"
)

const (
	New        string = "NEW"
	Registered string = "REGISTERED"
	Processing string = "PROCESSING"
	Invalid    string = "INVALID"
	Processed  string = "PROCESSED"
)

// Model for list orders information
type ListOrdersInfo struct {
	Number     string    `json:"numder"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	IntAccrual int       `json:"-"`
	TimeOrder  time.Time `json:"-"`
	Uploaded   string    `json:"uploaded_at"`
}

// Model current balance
type CurrentBalance struct {
	Current      float64 `json:"current"`
	Withdrawn    float64 `json:"withdrawn"`
	IntCurrent   int     `json:"-"`
	IntWithdrawn int     `json:"-"`
}

// Model withdrawal points
type BalanceWithdraw struct {
	Order  string  `json:"order"`
	Sum    float64 `json:"sum"`
	IntSum int     `json:"-"`
}

// Model withdrawals
type Withdrawals struct {
	Order     string    `json:"order"`
	Sum       float64   `json:"sum"`
	IntSum    int       `json:"-"`
	TimeOrder time.Time `json:"-"`
	Processed string    `json:"processed_at"`
}

// Model order information
type OrderInfo struct {
	Order      string  `json:"order"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	IntAccrual int     `json:"-"`
	UUID       string  `json:"-"`
}

// Model unprocessed order
type UnprocessedOrder struct {
	Order string
	UUID  string
}
