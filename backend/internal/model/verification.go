package model

import "time"

// VerificationCode represents a short-lived code for signup or password reset.
type VerificationCode struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	Purpose   string    `json:"purpose"`
	ExpiresAt time.Time `json:"expires_at"`
	ConsumedAt *time.Time `json:"consumed_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
