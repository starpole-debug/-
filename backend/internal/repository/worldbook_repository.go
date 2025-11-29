package repository

import (
	"context"
	"encoding/json"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorldbookRepository manages world/setting snippets tied to roles.
type WorldbookRepository struct {
	pool *pgxpool.Pool
}

func NewWorldbookRepository(pool *pgxpool.Pool) *WorldbookRepository {
	return &WorldbookRepository{pool: pool}
}

func (r *WorldbookRepository) FindByRole(ctx context.Context, roleID string) (*model.Worldbook, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, role_id, data_json, created_at, updated_at
        FROM worldbooks WHERE role_id = $1
    `, roleID)
	var (
		record  model.Worldbook
		dataRaw []byte
	)
	if err := row.Scan(&record.ID, &record.RoleID, &dataRaw, &record.CreatedAt, &record.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(dataRaw) > 0 {
		if err := json.Unmarshal(dataRaw, &record.Data); err != nil {
			record.Data = map[string]any{}
		}
	}
	return &record, nil
}
