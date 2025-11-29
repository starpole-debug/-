package repository

import (
	"context"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// VerificationRepository manages email verification / reset codes.
type VerificationRepository struct {
	pool *pgxpool.Pool
}

func NewVerificationRepository(pool *pgxpool.Pool) *VerificationRepository {
	return &VerificationRepository{pool: pool}
}

// Upsert stores a code for a given email/purpose, replacing any existing pending code.
func (r *VerificationRepository) Upsert(ctx context.Context, email, purpose, code string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO email_verification_codes (email, purpose, code, expires_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (email, purpose)
        DO UPDATE SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at, consumed_at = NULL, created_at = now()
    `, email, purpose, code, expiresAt)
	return err
}

// Find retrieves a pending code for validation.
func (r *VerificationRepository) Find(ctx context.Context, email, purpose string) (*model.VerificationCode, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, email, code, purpose, expires_at, consumed_at, created_at
        FROM email_verification_codes
        WHERE email = $1 AND purpose = $2
    `, email, purpose)
	var v model.VerificationCode
	if err := row.Scan(&v.ID, &v.Email, &v.Code, &v.Purpose, &v.ExpiresAt, &v.ConsumedAt, &v.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

// Consume marks a code as used.
func (r *VerificationRepository) Consume(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE email_verification_codes
        SET consumed_at = now()
        WHERE id = $1 AND consumed_at IS NULL
    `, id)
	return err
}
