package model

import "time"

// PaymentOrder records a recharge transaction.
type PaymentOrder struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	OutTradeNo      string    `json:"out_trade_no"`
	ProviderTradeNo string    `json:"provider_trade_no"`
	PayType         string    `json:"pay_type"`
	Status          string    `json:"status"` // pending | paid | failed
	MoneyCents      int64     `json:"money_cents"`
	Coins           int64     `json:"coins"`
	NotifyPayload   any       `json:"notify_payload"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

