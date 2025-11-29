package repository

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PresetRepository struct {
	pool *pgxpool.Pool
	once sync.Once
	err  error
}

func NewPresetRepository(pool *pgxpool.Pool) *PresetRepository {
	return &PresetRepository{pool: pool}
}

// ensureTable defensively creates the presets table if migrations were not applied.
func (r *PresetRepository) ensureTable(ctx context.Context) error {
	r.once.Do(func() {
		_, r.err = r.pool.Exec(ctx, `
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
			CREATE TABLE IF NOT EXISTS presets (
				id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
				creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				name TEXT NOT NULL,
				description TEXT,
				model_key TEXT,
				blocks JSONB NOT NULL DEFAULT '[]',
				gen_params JSONB,
				is_public BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
			);
			CREATE INDEX IF NOT EXISTS idx_presets_creator ON presets(creator_id);
			CREATE INDEX IF NOT EXISTS idx_presets_public ON presets(is_public);
		`)
	})
	return r.err
}

func (r *PresetRepository) Create(ctx context.Context, preset *model.Preset) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	if preset.ID == "" {
		preset.ID = uuid.NewString()
	}
	blocksJSON, err := json.Marshal(preset.Blocks)
	if err != nil {
		return err
	}
	genParamsJSON, err := json.Marshal(preset.GenParams)
	if err != nil {
		return err
	}
	query := `
		INSERT INTO presets (id, creator_id, name, description, model_key, blocks, gen_params, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`
	return r.pool.QueryRow(ctx, query,
		preset.ID, preset.CreatorID, preset.Name, preset.Description, preset.ModelKey, blocksJSON, genParamsJSON, preset.IsPublic, time.Now(), time.Now(),
	).Scan(&preset.CreatedAt, &preset.UpdatedAt)
}

func (r *PresetRepository) Update(ctx context.Context, preset *model.Preset) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	blocksJSON, err := json.Marshal(preset.Blocks)
	if err != nil {
		return err
	}
	genParamsJSON, err := json.Marshal(preset.GenParams)
	if err != nil {
		return err
	}
	query := `
		UPDATE presets
		SET name = $2, description = $3, model_key = $4, blocks = $5, gen_params = $6, is_public = $7, updated_at = now()
		WHERE id = $1
		RETURNING updated_at
	`
	return r.pool.QueryRow(ctx, query,
		preset.ID, preset.Name, preset.Description, preset.ModelKey, blocksJSON, genParamsJSON, preset.IsPublic,
	).Scan(&preset.UpdatedAt)
}

func (r *PresetRepository) Delete(ctx context.Context, id string) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, "DELETE FROM presets WHERE id = $1", id)
	return err
}

func (r *PresetRepository) UpdatePublic(ctx context.Context, id string, isPublic bool) error {
	if err := r.ensureTable(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `
		UPDATE presets SET is_public = $2, updated_at = now() WHERE id = $1
	`, id, isPublic)
	return err
}

func (r *PresetRepository) FindByID(ctx context.Context, id string) (*model.Preset, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	query := `
		SELECT id, creator_id, name, description, model_key, blocks, gen_params, is_public, created_at, updated_at
		FROM presets WHERE id = $1
	`
	var p model.Preset
	var blocksJSON []byte
	var genParamsJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CreatorID, &p.Name, &p.Description, &p.ModelKey, &blocksJSON, &genParamsJSON, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if len(blocksJSON) > 0 {
		_ = json.Unmarshal(blocksJSON, &p.Blocks)
	}
	if len(genParamsJSON) > 0 {
		_ = json.Unmarshal(genParamsJSON, &p.GenParams)
	}
	return &p, nil
}

func (r *PresetRepository) ListByCreator(ctx context.Context, creatorID string) ([]model.Preset, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	query := `
		SELECT id, creator_id, name, description, model_key, blocks, gen_params, is_public, created_at, updated_at
		FROM presets WHERE creator_id = $1 ORDER BY updated_at DESC
	`
	rows, err := r.pool.Query(ctx, query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var presets []model.Preset
	for rows.Next() {
		var p model.Preset
		var blocksJSON []byte
		var genParamsJSON []byte
		if err := rows.Scan(&p.ID, &p.CreatorID, &p.Name, &p.Description, &p.ModelKey, &blocksJSON, &genParamsJSON, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if len(blocksJSON) > 0 {
			_ = json.Unmarshal(blocksJSON, &p.Blocks)
		}
		if len(genParamsJSON) > 0 {
			_ = json.Unmarshal(genParamsJSON, &p.GenParams)
		}
		presets = append(presets, p)
	}
	return presets, nil
}

func (r *PresetRepository) ListPublic(ctx context.Context, limit int) ([]model.Preset, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id, creator_id, name, description, model_key, blocks, gen_params, is_public, created_at, updated_at
		FROM presets WHERE is_public = TRUE
		ORDER BY updated_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var presets []model.Preset
	for rows.Next() {
		var p model.Preset
		var blocksJSON []byte
		var genParamsJSON []byte
		if err := rows.Scan(&p.ID, &p.CreatorID, &p.Name, &p.Description, &p.ModelKey, &blocksJSON, &genParamsJSON, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if len(blocksJSON) > 0 {
			_ = json.Unmarshal(blocksJSON, &p.Blocks)
		}
		if len(genParamsJSON) > 0 {
			_ = json.Unmarshal(genParamsJSON, &p.GenParams)
		}
		presets = append(presets, p)
	}
	return presets, nil
}
