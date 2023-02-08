package models

// Model for list orders information
type ListOrdersInfo struct {
	Number   string `json:"numder"`
	Status   string `json:"status"`
	Accrual  uint32 `json:"accrual,omitempty"`
	Uploaded string `json:"uploaded_at"`
}

// Model order information
type OrderInfo struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual uint32 `json:"accrual,omitempty"`
}

// Model current balance
type CurrentBalance struct {
	Current   uint32 `json:"current"`
	Withdrawn uint32 `json:"withdrawn"`
}

// Model balance withdraw
type BalanceWithdraw struct {
	Order string `json:"order"`
	Sum   uint32 `json:"sum"`
}

// Model withdrawals
type Withdrawals struct {
	Order     string `json:"order"`
	Sum       uint32 `json:"sum"`
	Processed string `json:"processed_at"`
}
