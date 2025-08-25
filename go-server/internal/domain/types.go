package domain

import (
	"time"
)

// Wallet represents a spending wallet with threshold approvals
type Wallet struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Owners             []string  `json:"owners"`
	DailyLimitCents    int64     `json:"dailyLimitCents"`
	RequiredApprovals  int       `json:"requiredApprovals"`
	CreatedAt          time.Time `json:"createdAt"`
}

// Transaction represents a spending request that may need approvals
type Transaction struct {
	ID          string    `json:"id"`
	WalletID    string    `json:"walletId"`
	AmountCents int64     `json:"amountCents"`
	Memo        string    `json:"memo"`
	CreatedBy   string    `json:"createdBy"`
	Approvers   []string  `json:"approvers"`
	Status      string    `json:"status"` // pending, approved, rejected, executed
	CreatedAt   time.Time `json:"createdAt"`
}
