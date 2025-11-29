package repository

import (
	"context"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ImagePresetRepository stores prompt presets.
type ImagePresetRepository struct {
	pool *pgxpool.Pool
}

func NewImagePresetRepository(pool *pgxpool.Pool) *ImagePresetRepository {
	return &ImagePresetRepository{pool: pool}
}

func (r *ImagePresetRepository) Active(ctx context.Context) (*model.ImagePreset, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, name, preset_json, prompt_model_key, status, created_at, updated_at
        FROM image_presets
        WHERE status = 'active'
        ORDER BY updated_at DESC
        LIMIT 1
    `)
	var p model.ImagePreset
	if err := row.Scan(&p.ID, &p.Name, &p.PresetJSON, &p.PromptModelKey, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *ImagePresetRepository) List(ctx context.Context) ([]model.ImagePreset, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, name, preset_json, prompt_model_key, status, created_at, updated_at
        FROM image_presets
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var presets []model.ImagePreset
	for rows.Next() {
		var p model.ImagePreset
		if err := rows.Scan(&p.ID, &p.Name, &p.PresetJSON, &p.PromptModelKey, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		presets = append(presets, p)
	}
	return presets, nil
}

func (r *ImagePresetRepository) Save(ctx context.Context, p *model.ImagePreset) (*model.ImagePreset, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.pool.Exec(ctx, `
        INSERT INTO image_presets(id, name, preset_json, prompt_model_key, status, created_at, updated_at)
        VALUES($1,$2,$3,$4,$5,$6,$7)
        ON CONFLICT (id) DO UPDATE SET
            name=EXCLUDED.name,
            preset_json=EXCLUDED.preset_json,
            prompt_model_key=EXCLUDED.prompt_model_key,
            status=EXCLUDED.status,
            updated_at=EXCLUDED.updated_at
    `, p.ID, p.Name, p.PresetJSON, p.PromptModelKey, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ImagePresetRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM image_presets WHERE id = $1`, id)
	return err
}

func (r *ImagePresetRepository) FindByID(ctx context.Context, id string) (*model.ImagePreset, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, name, preset_json, prompt_model_key, status, created_at, updated_at
        FROM image_presets WHERE id = $1
    `, id)
	var p model.ImagePreset
	if err := row.Scan(&p.ID, &p.Name, &p.PresetJSON, &p.PromptModelKey, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
