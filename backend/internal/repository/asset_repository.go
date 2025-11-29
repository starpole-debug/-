package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserAssetRepository provides access to the user_assets table.
type UserAssetRepository struct {
	pool *pgxpool.Pool
}

func NewUserAssetRepository(pool *pgxpool.Pool) *UserAssetRepository {
	return &UserAssetRepository{pool: pool}
}

// GetByUser returns the user's asset record or a zero struct if missing.
func (r *UserAssetRepository) GetByUser(ctx context.Context, userID string) (*model.UserAsset, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT user_id, balance, monthly_tickets, created_at, updated_at
        FROM user_assets WHERE user_id = $1
    `, userID)
	var asset model.UserAsset
	if err := row.Scan(&asset.UserID, &asset.Balance, &asset.MonthlyTickets, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return &model.UserAsset{UserID: userID}, nil
		}
		return nil, err
	}
	return &asset, nil
}

// Upsert updates or inserts the asset row.
func (r *UserAssetRepository) Upsert(ctx context.Context, asset *model.UserAsset) error {
	row := r.pool.QueryRow(ctx, `
        INSERT INTO user_assets(user_id, balance, monthly_tickets)
        VALUES($1,$2,$3)
        ON CONFLICT (user_id) DO UPDATE
        SET balance = EXCLUDED.balance,
            monthly_tickets = EXCLUDED.monthly_tickets,
            updated_at = now()
        RETURNING created_at, updated_at
    `, asset.UserID, asset.Balance, asset.MonthlyTickets)
	return row.Scan(&asset.CreatedAt, &asset.UpdatedAt)
}
