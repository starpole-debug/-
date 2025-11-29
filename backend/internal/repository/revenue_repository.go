package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RevenueRepository persists monetisation ledgers.
type RevenueRepository struct {
	pool *pgxpool.Pool
}

func NewRevenueRepository(pool *pgxpool.Pool) *RevenueRepository {
	return &RevenueRepository{pool: pool}
}

func (r *RevenueRepository) GetWallet(ctx context.Context, creatorID string) (*model.CreatorWallet, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT creator_id, available_balance, frozen_balance, total_earned, updated_at
        FROM creator_wallets WHERE creator_id = $1
    `, creatorID)
	var wallet model.CreatorWallet
	if err := row.Scan(&wallet.CreatorID, &wallet.AvailableBalance, &wallet.FrozenBalance, &wallet.TotalEarned, &wallet.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			wallet.CreatorID = creatorID
			return &wallet, nil
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *RevenueRepository) UpsertWallet(ctx context.Context, wallet *model.CreatorWallet) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO creator_wallets(creator_id, available_balance, frozen_balance, total_earned)
        VALUES($1,$2,$3,$4)
        ON CONFLICT (creator_id) DO UPDATE SET
            available_balance = EXCLUDED.available_balance,
            frozen_balance = EXCLUDED.frozen_balance,
            total_earned = EXCLUDED.total_earned,
            updated_at = now()
    `, wallet.CreatorID, wallet.AvailableBalance, wallet.FrozenBalance, wallet.TotalEarned)
	return err
}

func (r *RevenueRepository) CreateEvent(ctx context.Context, event *model.RevenueEvent) error {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO revenue_events(id, creator_id, user_id, role_id, event_type, amount, status)
        VALUES($1,$2,$3,$4,$5,$6,$7)
        RETURNING created_at
    `, event.ID, event.CreatorID, event.UserID, event.RoleID, event.EventType, event.Amount, event.Status)
	return row.Scan(&event.CreatedAt)
}

func (r *RevenueRepository) ListEvents(ctx context.Context, creatorID string, limit int) ([]model.RevenueEvent, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, creator_id, user_id, role_id, event_type, amount, status, created_at
        FROM revenue_events WHERE creator_id = $1 ORDER BY created_at DESC LIMIT $2
    `, creatorID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var events []model.RevenueEvent
	for rows.Next() {
		var event model.RevenueEvent
		if err := rows.Scan(&event.ID, &event.CreatorID, &event.UserID, &event.RoleID, &event.EventType, &event.Amount, &event.Status, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *RevenueRepository) CreatePayout(ctx context.Context, payout *model.PayoutRecord) error {
	if payout.ID == "" {
		payout.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO payout_records(id, creator_id, amount, status, channel)
        VALUES($1,$2,$3,$4,$5)
        RETURNING created_at, updated_at
    `, payout.ID, payout.CreatorID, payout.Amount, payout.Status, payout.Channel)
	return row.Scan(&payout.CreatedAt, &payout.UpdatedAt)
}

func (r *RevenueRepository) ListPayouts(ctx context.Context, creatorID string) ([]model.PayoutRecord, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, creator_id, amount, status, channel, created_at, updated_at
        FROM payout_records WHERE creator_id = $1 ORDER BY created_at DESC
    `, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payouts []model.PayoutRecord
	for rows.Next() {
		var payout model.PayoutRecord
		if err := rows.Scan(&payout.ID, &payout.CreatorID, &payout.Amount, &payout.Status, &payout.Channel, &payout.CreatedAt, &payout.UpdatedAt); err != nil {
			return nil, err
		}
		payouts = append(payouts, payout)
	}
	return payouts, nil
}

func (r *RevenueRepository) FindPayout(ctx context.Context, id string) (*model.PayoutRecord, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, creator_id, amount, status, channel, created_at, updated_at
        FROM payout_records WHERE id = $1
    `, id)
	var payout model.PayoutRecord
	if err := row.Scan(&payout.ID, &payout.CreatorID, &payout.Amount, &payout.Status, &payout.Channel, &payout.CreatedAt, &payout.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &payout, nil
}

// ListAllPayouts provides paginated payout requests for admin review.
func (r *RevenueRepository) ListAllPayouts(ctx context.Context, status string, limit, offset int) ([]model.PayoutRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	query := `
		SELECT id, creator_id, amount, status, channel, created_at, updated_at
		FROM payout_records
		WHERE ($1 = '' OR status = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pool.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payouts []model.PayoutRecord
	for rows.Next() {
		var payout model.PayoutRecord
		if err := rows.Scan(&payout.ID, &payout.CreatorID, &payout.Amount, &payout.Status, &payout.Channel, &payout.CreatedAt, &payout.UpdatedAt); err != nil {
			return nil, err
		}
		payouts = append(payouts, payout)
	}
	return payouts, nil
}

// UpdatePayoutStatus updates status and returns affected payout.
func (r *RevenueRepository) UpdatePayoutStatus(ctx context.Context, id, status string) (*model.PayoutRecord, error) {
	row := r.pool.QueryRow(ctx, `
		UPDATE payout_records
		SET status = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, creator_id, amount, status, channel, created_at, updated_at
	`, id, status)
	var payout model.PayoutRecord
	if err := row.Scan(&payout.ID, &payout.CreatorID, &payout.Amount, &payout.Status, &payout.Channel, &payout.CreatedAt, &payout.UpdatedAt); err != nil {
		return nil, err
	}
	return &payout, nil
}

func (r *RevenueRepository) ListRules(ctx context.Context) ([]model.RevenueRule, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, event_type, rate, amount, enabled FROM revenue_rules ORDER BY event_type
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rules []model.RevenueRule
	for rows.Next() {
		var rule model.RevenueRule
		if err := rows.Scan(&rule.ID, &rule.EventType, &rule.Rate, &rule.Amount, &rule.Enabled); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *RevenueRepository) SaveRule(ctx context.Context, rule *model.RevenueRule) error {
	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	_, err := r.pool.Exec(ctx, `
        INSERT INTO revenue_rules(id, event_type, rate, amount, enabled)
        VALUES($1,$2,$3,$4,$5)
        ON CONFLICT (id) DO UPDATE SET
            event_type = EXCLUDED.event_type,
            rate = EXCLUDED.rate,
            amount = EXCLUDED.amount,
            enabled = EXCLUDED.enabled
    `, rule.ID, rule.EventType, rule.Rate, rule.Amount, rule.Enabled)
	return err
}
