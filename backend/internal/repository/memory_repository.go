package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MemoryRepository stores persona specific memories.
type MemoryRepository struct {
	pool *pgxpool.Pool
}

func NewMemoryRepository(pool *pgxpool.Pool) *MemoryRepository {
	return &MemoryRepository{pool: pool}
}

func (r *MemoryRepository) List(ctx context.Context, userID, roleID string) ([]model.MemoryCapsule, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, user_id, role_id, content, created_at
        FROM memory_capsules
        WHERE user_id = $1 AND role_id = $2
        ORDER BY created_at DESC
    `, userID, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.MemoryCapsule
	for rows.Next() {
		var capsule model.MemoryCapsule
		if err := rows.Scan(&capsule.ID, &capsule.UserID, &capsule.RoleID, &capsule.Content, &capsule.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, capsule)
	}
	return items, nil
}

func (r *MemoryRepository) Create(ctx context.Context, capsule *model.MemoryCapsule) error {
	if capsule.ID == "" {
		capsule.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO memory_capsules(id, user_id, role_id, content)
        VALUES($1,$2,$3,$4)
        RETURNING created_at
    `, capsule.ID, capsule.UserID, capsule.RoleID, capsule.Content)
	return row.Scan(&capsule.CreatedAt)
}

func (r *MemoryRepository) Delete(ctx context.Context, id, userID string) error {
	res, err := r.pool.Exec(ctx, `DELETE FROM memory_capsules WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
