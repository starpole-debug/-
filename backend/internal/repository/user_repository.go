package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository wraps CRUD operations for the users table.
type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO users (id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned)
        VALUES ($1, $2, $3, $4, $5, $6, $7, COALESCE($8, FALSE))
        RETURNING created_at, updated_at
    `, user.ID, user.Username, user.Email, user.PasswordHash, user.Nickname, user.AvatarURL, user.IsAdmin, user.IsBanned)
	return row.Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned, created_at, updated_at, deleted_at
        FROM users WHERE username = $1 AND deleted_at IS NULL
    `, username)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL, &user.IsAdmin, &user.IsBanned, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned, created_at, updated_at, deleted_at
        FROM users WHERE id = $1 AND deleted_at IS NULL
    `, id)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL, &user.IsAdmin, &user.IsBanned, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, query string, limit, offset int) ([]model.User, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	base := `
        SELECT id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned, created_at, updated_at, deleted_at
        FROM users
        WHERE deleted_at IS NULL
    `
	args := []interface{}{}
	idx := 1
	builder := strings.Builder{}
	builder.WriteString(base)
	if strings.TrimSpace(query) != "" {
		query = strings.TrimSpace(query)
		builder.WriteString(fmt.Sprintf(" AND (username ILIKE $%d OR email ILIKE $%d OR nickname ILIKE $%d OR id::TEXT ILIKE $%d)", idx, idx+1, idx+2, idx+3))
		pattern := "%" + query + "%"
		args = append(args, pattern, pattern, pattern, pattern)
		idx += 4
	}
	builder.WriteString(fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1))
	args = append(args, limit, offset)
	rows, err := r.pool.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL, &user.IsAdmin, &user.IsBanned, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) SetBanStatus(ctx context.Context, id string, banned bool) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE users SET is_banned = $2, updated_at = now()
        WHERE id = $1 AND deleted_at IS NULL
    `, id, banned)
	return err
}

func (r *UserRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE users SET deleted_at = now(), updated_at = now()
        WHERE id = $1 AND deleted_at IS NULL
    `, id)
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned, created_at, updated_at, deleted_at
        FROM users WHERE email = $1 AND deleted_at IS NULL
    `, email)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL, &user.IsAdmin, &user.IsBanned, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) updateAdmin(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE users SET
            username = $2,
            email = $3,
            password_hash = $4,
            nickname = $5,
            avatar_url = $6,
            is_admin = TRUE,
            is_banned = FALSE,
            deleted_at = NULL,
            updated_at = now()
        WHERE id = $1
    `, user.ID, user.Username, user.Email, user.PasswordHash, user.Nickname, user.AvatarURL)
	return err
}

// UpsertAdmin ensures an admin user exists with the provided credentials.
func (r *UserRepository) UpsertAdmin(ctx context.Context, user *model.User) error {
	if user.ID == "" {
		user.ID = uuid.NewString()
	}
	if user.Nickname == "" {
		user.Nickname = user.Username
	}
	user.IsAdmin = true
	existing, err := r.FindByUsername(ctx, user.Username)
	if err != nil {
		return err
	}
	if existing == nil {
		existing, err = r.FindByEmail(ctx, user.Email)
		if err != nil {
			return err
		}
	}
	if existing != nil {
		user.ID = existing.ID
		if user.Nickname == "" {
			user.Nickname = existing.Nickname
		}
		if err := r.updateAdmin(ctx, user); err != nil {
			return err
		}
		user.CreatedAt = existing.CreatedAt
		user.UpdatedAt = time.Now()
		return nil
	}
	return r.Create(ctx, user)
}

// UpdateProfile updates nickname/avatar for a given user.
func (r *UserRepository) UpdateProfile(ctx context.Context, userID, nickname, avatarURL string) (*model.User, error) {
	row := r.pool.QueryRow(ctx, `
        UPDATE users
        SET nickname = CASE WHEN $2 <> '' THEN $2 ELSE nickname END,
            avatar_url = CASE WHEN $3 <> '' THEN $3 ELSE avatar_url END,
            updated_at = now()
        WHERE id = $1
        RETURNING id, username, email, password_hash, nickname, avatar_url, is_admin, is_banned, created_at, updated_at, deleted_at
    `, userID, nickname, avatarURL)
	var user model.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Nickname, &user.AvatarURL, &user.IsAdmin, &user.IsBanned, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword updates the password hash for a user.
func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE users
        SET password_hash = $2, updated_at = now()
        WHERE id = $1 AND deleted_at IS NULL
    `, userID, hash)
	return err
}
