package repository

import (
	"context"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ImageProviderRepository manages image provider configs.
type ImageProviderRepository struct {
	pool *pgxpool.Pool
}

func NewImageProviderRepository(pool *pgxpool.Pool) *ImageProviderRepository {
	return &ImageProviderRepository{pool: pool}
}

func (r *ImageProviderRepository) ListActive(ctx context.Context) ([]model.ImageProvider, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, name, base_url, api_key, max_concurrency, weight, status, params_json, selected_model, created_at, updated_at
        FROM image_providers
        WHERE status = 'active'
        ORDER BY created_at ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var providers []model.ImageProvider
	for rows.Next() {
		var p model.ImageProvider
		if err := rows.Scan(&p.ID, &p.Name, &p.BaseURL, &p.APIKey, &p.MaxConcurrency, &p.Weight, &p.Status, &p.ParamsJSON, &p.SelectedModel, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func (r *ImageProviderRepository) List(ctx context.Context) ([]model.ImageProvider, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, name, base_url, api_key, max_concurrency, weight, status, params_json, selected_model, created_at, updated_at
        FROM image_providers
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var providers []model.ImageProvider
	for rows.Next() {
		var p model.ImageProvider
		if err := rows.Scan(&p.ID, &p.Name, &p.BaseURL, &p.APIKey, &p.MaxConcurrency, &p.Weight, &p.Status, &p.ParamsJSON, &p.SelectedModel, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}
	return providers, nil
}

func (r *ImageProviderRepository) Save(ctx context.Context, p *model.ImageProvider) (*model.ImageProvider, error) {
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	_, err := r.pool.Exec(ctx, `
        INSERT INTO image_providers(id, name, base_url, api_key, max_concurrency, weight, status, params_json, selected_model, created_at, updated_at)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (id) DO UPDATE SET
            name=EXCLUDED.name,
            base_url=EXCLUDED.base_url,
            api_key=EXCLUDED.api_key,
            max_concurrency=EXCLUDED.max_concurrency,
            weight=EXCLUDED.weight,
            status=EXCLUDED.status,
            params_json=EXCLUDED.params_json,
            selected_model=EXCLUDED.selected_model,
            updated_at=EXCLUDED.updated_at
    `, p.ID, p.Name, p.BaseURL, p.APIKey, p.MaxConcurrency, p.Weight, p.Status, p.ParamsJSON, p.SelectedModel, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ImageProviderRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM image_providers WHERE id = $1`, id)
	return err
}

func (r *ImageProviderRepository) FindByID(ctx context.Context, id string) (*model.ImageProvider, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, name, base_url, api_key, max_concurrency, weight, status, params_json, selected_model, created_at, updated_at
        FROM image_providers WHERE id = $1
    `, id)
	var p model.ImageProvider
	if err := row.Scan(&p.ID, &p.Name, &p.BaseURL, &p.APIKey, &p.MaxConcurrency, &p.Weight, &p.Status, &p.ParamsJSON, &p.SelectedModel, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}
