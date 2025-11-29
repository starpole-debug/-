package model

import "time"

// UserAsset stores platform virtual currency + tickets for a user.
type UserAsset struct {
	UserID         string    `json:"user_id"`
	Balance        int64     `json:"balance"`
	MonthlyTickets int64     `json:"monthly_tickets"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
