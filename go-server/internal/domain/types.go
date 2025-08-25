package domain

import (
	"time"
)

// Wallet represents a spending wallet with threshold approvals
type Wallet struct {
	ID string
	Name string
	Owners []string
	DailyLimitCents int64
	RequiredApprovals int
	CreatedAt time.Time
}

// Transaction represents a spending request that may need approvals
type Transaction struct {
	ID string
	WalletID string
	AmountCents int64
	Memo string
	CreatedBy string
	Approvers []string
	Status string // pending, approved, rejected, executed
	CreatedAt time.Time
}
