package repository

import (
	"context"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RoleRepository manages CRUD for roles and versions.
type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

// ensureRoleDataColumn keeps the schema aligned when migrations were skipped.
func (r *RoleRepository) ensureRoleDataColumn(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `ALTER TABLE roles ADD COLUMN IF NOT EXISTS data JSONB`)
	return err
}

// ensureFavoritesTable is a defensive guard so new instances work even if migrations were not applied.
func (r *RoleRepository) ensureFavoritesTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
        CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
        CREATE TABLE IF NOT EXISTS role_favorites (
            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
            created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
            UNIQUE(user_id, role_id)
        );
        CREATE INDEX IF NOT EXISTS idx_role_favorites_user ON role_favorites(user_id);
        CREATE INDEX IF NOT EXISTS idx_role_favorites_role ON role_favorites(role_id);
    `)
	return err
}

func (r *RoleRepository) List(ctx context.Context, status string, limit int) ([]model.Role, error) {
	if limit <= 0 {
		limit = 20
	}
	if err := r.ensureRoleDataColumn(ctx); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
        SELECT id, creator_id, name, description, avatar_url, tags, abilities, allow_clone, status, coalesce(role_version,''), data, created_at, updated_at
        FROM roles
        WHERE ($1 = '' OR status = $1)
        ORDER BY updated_at DESC
        LIMIT $2
    `, status, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var tags []string
		var abilities []string
		if err := rows.Scan(&role.ID, &role.CreatorID, &role.Name, &role.Description, &role.AvatarURL, &tags, &abilities, &role.AllowClone, &role.Status, &role.Version, &role.Data, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		role.Tags = tags
		role.Abilities = abilities
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	if err := r.ensureRoleDataColumn(ctx); err != nil {
		return nil, err
	}
	row := r.pool.QueryRow(ctx, `
        SELECT r.id, r.creator_id, r.name, r.description, r.avatar_url, r.tags, r.abilities, r.allow_clone, r.status, coalesce(r.role_version,''), r.data, r.created_at, r.updated_at,
               COALESCE(f.cnt, 0) AS favorite_count
        FROM roles r
        LEFT JOIN (
            SELECT role_id, COUNT(*) AS cnt FROM role_favorites WHERE role_id = $1 GROUP BY role_id
        ) f ON r.id = f.role_id
        WHERE r.id = $1
    `, id)
	var role model.Role
	var tags []string
	var abilities []string
	if err := row.Scan(&role.ID, &role.CreatorID, &role.Name, &role.Description, &role.AvatarURL, &tags, &abilities, &role.AllowClone, &role.Status, &role.Version, &role.Data, &role.CreatedAt, &role.UpdatedAt, &role.FavoriteCnt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	role.Tags = tags
	role.Abilities = abilities
	return &role, nil
}

func (r *RoleRepository) Save(ctx context.Context, role *model.Role) error {
	if role.ID == "" {
		role.ID = uuid.NewString()
	}
	if err := r.ensureRoleDataColumn(ctx); err != nil {
		return err
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO roles (id, creator_id, name, description, avatar_url, tags, abilities, allow_clone, status, role_version, data)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (id) DO UPDATE SET
            name = EXCLUDED.name,
            description = EXCLUDED.description,
            avatar_url = EXCLUDED.avatar_url,
            tags = EXCLUDED.tags,
            abilities = EXCLUDED.abilities,
            allow_clone = EXCLUDED.allow_clone,
            status = EXCLUDED.status,
            role_version = EXCLUDED.role_version,
            data = EXCLUDED.data,
            updated_at = now()
        RETURNING created_at, updated_at
    `, role.ID, role.CreatorID, role.Name, role.Description, role.AvatarURL, role.Tags, role.Abilities, role.AllowClone, role.Status, role.Version, role.Data)
	return row.Scan(&role.CreatedAt, &role.UpdatedAt)
}

func (r *RoleRepository) UpdateStatus(ctx context.Context, roleID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE roles SET status = $2, updated_at = now() WHERE id = $1`, roleID, status)
	return err
}

func (r *RoleRepository) Favorite(ctx context.Context, userID, roleID string) error {
	if err := r.ensureFavoritesTable(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `
        INSERT INTO role_favorites(user_id, role_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, role_id) DO NOTHING
    `, userID, roleID)
	return err
}

func (r *RoleRepository) Unfavorite(ctx context.Context, userID, roleID string) error {
	if err := r.ensureFavoritesTable(ctx); err != nil {
		return err
	}
	_, err := r.pool.Exec(ctx, `DELETE FROM role_favorites WHERE user_id = $1 AND role_id = $2`, userID, roleID)
	return err
}

func (r *RoleRepository) IsFavorited(ctx context.Context, userID, roleID string) (bool, error) {
	if err := r.ensureFavoritesTable(ctx); err != nil {
		return false, err
	}
	row := r.pool.QueryRow(ctx, `SELECT 1 FROM role_favorites WHERE user_id = $1 AND role_id = $2 LIMIT 1`, userID, roleID)
	var dummy int
	if err := row.Scan(&dummy); err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RoleRepository) CountFavorites(ctx context.Context, roleID string) (int, error) {
	if err := r.ensureFavoritesTable(ctx); err != nil {
		return 0, err
	}
	row := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM role_favorites WHERE role_id = $1`, roleID)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *RoleRepository) ListFavorites(ctx context.Context, userID string, limit int) ([]model.Role, error) {
	if err := r.ensureFavoritesTable(ctx); err != nil {
		return nil, err
	}
	if err := r.ensureRoleDataColumn(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 30
	}
	rows, err := r.pool.Query(ctx, `
        SELECT r.id, r.creator_id, r.name, r.description, r.avatar_url, r.tags, r.abilities, r.allow_clone, r.status, coalesce(r.role_version,''), r.data, r.created_at, r.updated_at,
               COUNT(f_all.user_id) AS favorite_count,
               MAX(f.created_at) AS favorited_at
        FROM role_favorites f
        JOIN roles r ON r.id = f.role_id
        LEFT JOIN role_favorites f_all ON f_all.role_id = r.id
        WHERE f.user_id = $1
        GROUP BY r.id, r.creator_id, r.name, r.description, r.avatar_url, r.tags, r.abilities, r.allow_clone, r.status, r.role_version, r.data, r.created_at, r.updated_at
        ORDER BY favorited_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var tags []string
		var abilities []string
		var favoritedAt time.Time
		if err := rows.Scan(&role.ID, &role.CreatorID, &role.Name, &role.Description, &role.AvatarURL, &tags, &abilities, &role.AllowClone, &role.Status, &role.Version, &role.Data, &role.CreatedAt, &role.UpdatedAt, &role.FavoriteCnt, &favoritedAt); err != nil {
			return nil, err
		}
		role.Tags = tags
		role.Abilities = abilities
		role.IsFavorited = true
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) ListByCreator(ctx context.Context, creatorID string) ([]model.Role, error) {
	if err := r.ensureRoleDataColumn(ctx); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
        SELECT id, creator_id, name, description, avatar_url, tags, abilities, allow_clone, status, coalesce(role_version,''), data, created_at, updated_at
        FROM roles WHERE creator_id = $1 ORDER BY updated_at DESC
    `, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []model.Role
	for rows.Next() {
		var role model.Role
		var tags []string
		var abilities []string
		if err := rows.Scan(&role.ID, &role.CreatorID, &role.Name, &role.Description, &role.AvatarURL, &tags, &abilities, &role.AllowClone, &role.Status, &role.Version, &role.Data, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		role.Tags = tags
		role.Abilities = abilities
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) CreateVersion(ctx context.Context, roleID, prompt string) error {
	_, err := r.pool.Exec(ctx, `
        INSERT INTO role_versions(id, role_id, prompt)
        VALUES($1,$2,$3)
    `, uuid.NewString(), roleID, prompt)
	return err
}

func (r *RoleRepository) CountByCreator(ctx context.Context, creatorID string) (total int, published int, err error) {
	row := r.pool.QueryRow(ctx, `
        SELECT COUNT(*) AS total,
               COUNT(*) FILTER (WHERE status = 'published') AS published
        FROM roles
        WHERE creator_id = $1
    `, creatorID)
	if err = row.Scan(&total, &published); err != nil {
		return 0, 0, err
	}
	return total, published, nil
}
