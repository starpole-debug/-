package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ConfigRepository handles dictionaries + model configs.
type ConfigRepository struct {
	pool *pgxpool.Pool
}

func NewConfigRepository(pool *pgxpool.Pool) *ConfigRepository {
	return &ConfigRepository{pool: pool}
}

func (r *ConfigRepository) ListModels(ctx context.Context, includeDisabled bool) ([]model.ModelConfig, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id,
               name,
               description,
               provider,
               base_url,
               model_name,
               is_default,
               is_enabled,
               status,
               max_context_tokens,
               price_coins,
               share_role_pct,
               share_preset_pct,
               coalesce(price_hint,''),
               temperature,
               max_tokens,
               api_key <> '' AS has_api_key,
               created_at,
               updated_at
        FROM models
        WHERE ($1 = true OR status = 'active')
        ORDER BY created_at DESC
    `, includeDisabled)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var models []model.ModelConfig
	for rows.Next() {
		var m model.ModelConfig
		if err := rows.Scan(&m.ID, &m.Name, &m.Description, &m.Provider, &m.BaseURL, &m.ModelName, &m.IsDefault, &m.IsEnabled, &m.Status, &m.MaxContextTokens, &m.PriceCoins, &m.ShareRolePct, &m.SharePresetPct, &m.PriceHint, &m.Temperature, &m.MaxTokens, &m.HasAPIKey, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		models = append(models, m)
	}
	return models, nil
}

func (r *ConfigRepository) FindModel(ctx context.Context, id string) (*model.ModelConfig, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id,
               name,
               description,
               provider,
               base_url,
               model_name,
               api_key,
               temperature,
               max_tokens,
               status,
               is_default,
               is_enabled,
               max_context_tokens,
               price_coins,
               share_role_pct,
               share_preset_pct,
               coalesce(price_hint,''),
               api_key <> '' AS has_api_key,
               created_at,
               updated_at
        FROM models WHERE id = $1
    `, id)
	var m model.ModelConfig
	if err := row.Scan(&m.ID, &m.Name, &m.Description, &m.Provider, &m.BaseURL, &m.ModelName, &m.APIKey, &m.Temperature, &m.MaxTokens, &m.Status, &m.IsDefault, &m.IsEnabled, &m.MaxContextTokens, &m.PriceCoins, &m.ShareRolePct, &m.SharePresetPct, &m.PriceHint, &m.HasAPIKey, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *ConfigRepository) DefaultModel(ctx context.Context) (*model.ModelConfig, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id,
               name,
               description,
               provider,
               base_url,
               model_name,
               api_key,
               temperature,
               max_tokens,
               status,
               is_default,
               is_enabled,
               max_context_tokens,
               price_coins,
               share_role_pct,
               share_preset_pct,
               coalesce(price_hint,''),
               api_key <> '' AS has_api_key,
               created_at,
               updated_at
        FROM models WHERE is_default = true LIMIT 1
    `)
	var m model.ModelConfig
	if err := row.Scan(&m.ID, &m.Name, &m.Description, &m.Provider, &m.BaseURL, &m.ModelName, &m.APIKey, &m.Temperature, &m.MaxTokens, &m.Status, &m.IsDefault, &m.IsEnabled, &m.MaxContextTokens, &m.PriceCoins, &m.ShareRolePct, &m.SharePresetPct, &m.PriceHint, &m.HasAPIKey, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			models, listErr := r.ListModels(ctx, true)
			if listErr != nil {
				return nil, listErr
			}
			if len(models) == 0 {
				return nil, err
			}
			return r.FindModel(ctx, models[0].ID)
		}
		return nil, err
	}
	return &m, nil
}

func (r *ConfigRepository) SaveModel(ctx context.Context, m *model.ModelConfig) error {
	if m.ID == "" {
		m.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO models(
            id,
            name,
            description,
            provider,
            base_url,
            model_name,
            api_key,
            temperature,
            max_tokens,
            status,
            is_default,
            is_enabled,
            max_context_tokens,
            price_coins,
            share_role_pct,
            share_preset_pct,
            price_hint
        )
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            description = EXCLUDED.description,
            provider = EXCLUDED.provider,
            base_url = EXCLUDED.base_url,
            model_name = EXCLUDED.model_name,
            api_key = CASE WHEN EXCLUDED.api_key = '' THEN models.api_key ELSE EXCLUDED.api_key END,
            temperature = EXCLUDED.temperature,
            max_tokens = EXCLUDED.max_tokens,
            status = EXCLUDED.status,
            is_default = EXCLUDED.is_default,
            is_enabled = EXCLUDED.is_enabled,
            max_context_tokens = EXCLUDED.max_context_tokens,
            share_role_pct = EXCLUDED.share_role_pct,
            share_preset_pct = EXCLUDED.share_preset_pct,
            price_coins = EXCLUDED.price_coins,
            price_hint = EXCLUDED.price_hint,
            updated_at = now()
        RETURNING created_at, updated_at, api_key <> '' AS has_api_key
    `, m.ID, m.Name, m.Description, m.Provider, m.BaseURL, m.ModelName, m.APIKey, m.Temperature, m.MaxTokens, m.Status, m.IsDefault, m.IsEnabled, m.MaxContextTokens, m.PriceCoins, m.ShareRolePct, m.SharePresetPct, m.PriceHint)
	if err := row.Scan(&m.CreatedAt, &m.UpdatedAt, &m.HasAPIKey); err != nil {
		return err
	}
	m.APIKey = ""
	return nil
}

func (r *ConfigRepository) DeleteModel(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM models WHERE id = $1`, id)
	return err
}

func (r *ConfigRepository) ListDictionary(ctx context.Context, group string) ([]model.DictionaryItem, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, group_key, key, label, description, ord, enabled, created_at, updated_at
        FROM config_dictionary
        WHERE ($1 = '' OR group_key = $1)
        ORDER BY ord ASC
    `, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []model.DictionaryItem
	for rows.Next() {
		var item model.DictionaryItem
		if err := rows.Scan(&item.ID, &item.Group, &item.Key, &item.Label, &item.Description, &item.Order, &item.Enabled, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ConfigRepository) SaveDictionaryItem(ctx context.Context, item *model.DictionaryItem) error {
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO config_dictionary(id, group_key, key, label, description, ord, enabled)
        VALUES($1,$2,$3,$4,$5,$6,$7)
        ON CONFLICT (id) DO UPDATE SET
            group_key = EXCLUDED.group_key,
            key = EXCLUDED.key,
            label = EXCLUDED.label,
            description = EXCLUDED.description,
            ord = EXCLUDED.ord,
            enabled = EXCLUDED.enabled,
            updated_at = now()
        RETURNING created_at, updated_at
    `, item.ID, item.Group, item.Key, item.Label, item.Description, item.Order, item.Enabled)
	return row.Scan(&item.CreatedAt, &item.UpdatedAt)
}

func (r *ConfigRepository) DeleteDictionaryItem(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM config_dictionary WHERE id = $1`, id)
	return err
}
