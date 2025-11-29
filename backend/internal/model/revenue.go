package model

import "time"

// CreatorWallet keeps tracked balances for payouts.
type CreatorWallet struct {
	CreatorID        string    `json:"creator_id"`
	AvailableBalance int64     `json:"available_balance"`
	FrozenBalance    int64     `json:"frozen_balance"`
	TotalEarned      int64     `json:"total_earned"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// RevenueEvent records monetisation triggers (tip/purchase/etc.).
type RevenueEvent struct {
	ID        string    `json:"id"`
	CreatorID string    `json:"creator_id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	EventType string    `json:"event_type"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PayoutRecord captures manual withdrawal requests.
type PayoutRecord struct {
	ID        string    `json:"id"`
	CreatorID string    `json:"creator_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	Channel   string    `json:"channel"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RevenueRule contains configurable monetisation coefficients.
type RevenueRule struct {
	ID        string  `json:"id"`
	EventType string  `json:"event_type"`
	Rate      float64 `json:"rate"`
	Amount    int64   `json:"amount"`
	Enabled   bool    `json:"enabled"`
}
