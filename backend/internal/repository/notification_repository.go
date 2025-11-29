package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepository provides CRUD for notifications.
type NotificationRepository struct {
	pool *pgxpool.Pool
}

func NewNotificationRepository(pool *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{pool: pool}
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID string, limit int) ([]model.Notification, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, user_id, type, title, content, is_read, created_at
        FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Content, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *NotificationRepository) Create(ctx context.Context, n *model.Notification) error {
	if n.ID == "" {
		n.ID = uuid.NewString()
	}
	_, err := r.pool.Exec(ctx, `
        INSERT INTO notifications(id, user_id, type, title, content, is_read)
        VALUES($1,$2,$3,$4,$5,$6)
    `, n.ID, n.UserID, n.Type, n.Title, n.Content, n.IsRead)
	return err
}

func (r *NotificationRepository) MarkRead(ctx context.Context, userID, notificationID string) error {
	res, err := r.pool.Exec(ctx, `
        UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2
    `, notificationID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *NotificationRepository) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET is_read = true WHERE user_id = $1`, userID)
	return err
}
