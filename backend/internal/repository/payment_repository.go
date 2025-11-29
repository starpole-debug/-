package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PaymentRepository persists recharge/payment orders.
type PaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) *PaymentRepository {
	return &PaymentRepository{pool: pool}
}

func (r *PaymentRepository) Create(ctx context.Context, order *model.PaymentOrder) error {
	row := r.pool.QueryRow(ctx, `
        INSERT INTO payment_orders(user_id, out_trade_no, provider_trade_no, pay_type, status, money_cents, coins, notify_payload)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8)
        RETURNING id, created_at, updated_at
    `, order.UserID, order.OutTradeNo, order.ProviderTradeNo, order.PayType, order.Status, order.MoneyCents, order.Coins, order.NotifyPayload)
	return row.Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *PaymentRepository) FindByOutTradeNo(ctx context.Context, outTradeNo string) (*model.PaymentOrder, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, user_id, out_trade_no, provider_trade_no, pay_type, status, money_cents, coins, notify_payload, created_at, updated_at
        FROM payment_orders WHERE out_trade_no = $1
    `, outTradeNo)
	var order model.PaymentOrder
	if err := row.Scan(&order.ID, &order.UserID, &order.OutTradeNo, &order.ProviderTradeNo, &order.PayType, &order.Status, &order.MoneyCents, &order.Coins, &order.NotifyPayload, &order.CreatedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) MarkPaid(ctx context.Context, outTradeNo, providerTradeNo string, payload map[string]string) (*model.PaymentOrder, error) {
	data, _ := json.Marshal(payload)
	row := r.pool.QueryRow(ctx, `
        UPDATE payment_orders
        SET status = 'paid',
            provider_trade_no = COALESCE(NULLIF($2,''), provider_trade_no),
            notify_payload = $3,
            updated_at = now()
        WHERE out_trade_no = $1 AND status <> 'paid'
        RETURNING id, user_id, out_trade_no, provider_trade_no, pay_type, status, money_cents, coins, notify_payload, created_at, updated_at
    `, outTradeNo, providerTradeNo, data)
	var order model.PaymentOrder
	if err := row.Scan(&order.ID, &order.UserID, &order.OutTradeNo, &order.ProviderTradeNo, &order.PayType, &order.Status, &order.MoneyCents, &order.Coins, &order.NotifyPayload, &order.CreatedAt, &order.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, outTradeNo, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE payment_orders SET status = $2, updated_at = now() WHERE out_trade_no = $1`, outTradeNo, status)
	return err
}

// ListByUser returns recent payment orders for a user, optionally filtered by status ("paid" | "pending" | "failed" | "all/empty").
func (r *PaymentRepository) ListByUser(ctx context.Context, userID string, limit, offset int, status string) ([]model.PaymentOrder, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	query := `
        SELECT id, user_id, out_trade_no, provider_trade_no, pay_type, status, money_cents, coins, notify_payload, created_at, updated_at
        FROM payment_orders
        WHERE user_id = $1
    `
	args := []interface{}{userID}
	if status != "" && status != "all" {
		query += " AND status = $2"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT $3 OFFSET $4"
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []model.PaymentOrder
	for rows.Next() {
		var o model.PaymentOrder
		if err := rows.Scan(&o.ID, &o.UserID, &o.OutTradeNo, &o.ProviderTradeNo, &o.PayType, &o.Status, &o.MoneyCents, &o.Coins, &o.NotifyPayload, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

// ListAll returns recent payment orders for admin audit.
func (r *PaymentRepository) ListAll(ctx context.Context, limit, offset int) ([]model.PaymentOrder, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	rows, err := r.pool.Query(ctx, `
        SELECT id, user_id, out_trade_no, provider_trade_no, pay_type, status, money_cents, coins, notify_payload, created_at, updated_at
        FROM payment_orders
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []model.PaymentOrder
	for rows.Next() {
		var o model.PaymentOrder
		if err := rows.Scan(&o.ID, &o.UserID, &o.OutTradeNo, &o.ProviderTradeNo, &o.PayType, &o.Status, &o.MoneyCents, &o.Coins, &o.NotifyPayload, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
