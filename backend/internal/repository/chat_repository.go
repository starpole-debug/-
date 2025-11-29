package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChatRepository persists sessions and messages.
type ChatRepository struct {
	pool *pgxpool.Pool
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{pool: pool}
}

func (r *ChatRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	if session.ID == "" {
		session.ID = uuid.NewString()
	}
	if session.Mode == "" {
		session.Mode = "sfw"
	}
	if session.Status == "" {
		session.Status = "active"
	}
	settingsJSON, err := json.Marshal(session.Settings)
	if err != nil {
		return err
	}
	row := r.pool.QueryRow(ctx, `
        INSERT INTO chat_sessions(
			id, user_id, role_id, model_key, title, mode, status, settings, settings_json
		)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$8)
        RETURNING created_at, updated_at
    `, session.ID, session.UserID, session.RoleID, session.ModelKey, session.Title, session.Mode, session.Status, settingsJSON)
	return row.Scan(&session.CreatedAt, &session.UpdatedAt)
}

func (r *ChatRepository) FindSession(ctx context.Context, id string) (*model.ChatSession, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT
			id,
			user_id,
			role_id,
			model_key,
			title,
			mode,
			status,
			COALESCE(settings, settings_json, '{}'::jsonb) AS settings,
			created_at,
			updated_at
        FROM chat_sessions
		WHERE id = $1
    `, id)
	var session model.ChatSession
	if err := scanSession(row, &session); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *ChatRepository) ListSessionsByUser(ctx context.Context, userID string, limit int) ([]model.ChatSession, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.pool.Query(ctx, `
        SELECT
            cs.id,
            cs.user_id,
            cs.role_id,
            cs.model_key,
            cs.title,
            cs.mode,
            cs.status,
            COALESCE(cs.settings, cs.settings_json, '{}'::jsonb) AS settings,
            cs.created_at,
            cs.updated_at,
            lm.content AS last_message
        FROM chat_sessions cs
        LEFT JOIN LATERAL (
            SELECT content
            FROM chat_messages cm
            WHERE cm.session_id = cs.id
            ORDER BY cm.created_at DESC
            LIMIT 1
        ) lm ON TRUE
        WHERE cs.user_id = $1
        ORDER BY cs.updated_at DESC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sessions []model.ChatSession
	for rows.Next() {
		var session model.ChatSession
		var settingsRaw []byte
		var last sql.NullString
		if err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RoleID,
			&session.ModelKey,
			&session.Title,
			&session.Mode,
			&session.Status,
			&settingsRaw,
			&session.CreatedAt,
			&session.UpdatedAt,
			&last,
		); err != nil {
			return nil, err
		}
		if len(settingsRaw) > 0 {
			_ = json.Unmarshal(settingsRaw, &session.Settings)
		}
		session.LastMsg = last.String
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (r *ChatRepository) AddMessage(ctx context.Context, msg *model.ChatMessage) error {
	if msg.ID == "" {
		msg.ID = uuid.NewString()
	}
	metaJSON, _ := json.Marshal(msg.Metadata)
	row := r.pool.QueryRow(ctx, `
        INSERT INTO chat_messages(id, session_id, role, content, is_important, metadata)
        VALUES($1,$2,$3,$4,$5,$6)
        RETURNING created_at
    `, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.IsImportant, metaJSON)
	if err := row.Scan(&msg.CreatedAt); err != nil {
		return err
	}
	_, _ = r.pool.Exec(ctx, `UPDATE chat_sessions SET updated_at = now() WHERE id = $1`, msg.SessionID)
	return nil
}

func (r *ChatRepository) ListMessages(ctx context.Context, sessionID string, limit int) ([]model.ChatMessage, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.pool.Query(ctx, `
        SELECT id, session_id, role, content, is_important, metadata, created_at
        FROM chat_messages WHERE session_id = $1 ORDER BY created_at ASC LIMIT $2
    `, sessionID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []model.ChatMessage
	for rows.Next() {
		var msg model.ChatMessage
		var metaRaw []byte
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.IsImportant, &metaRaw, &msg.CreatedAt); err != nil {
			return nil, err
		}
		if len(metaRaw) > 0 {
			_ = json.Unmarshal(metaRaw, &msg.Metadata)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *ChatRepository) UpdateMessageContent(ctx context.Context, id, sessionID, content string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE chat_messages SET content = $3, metadata = metadata, created_at = created_at
        WHERE id = $1 AND session_id = $2
    `, id, sessionID, content)
	return err
}

func (r *ChatRepository) DeleteMessage(ctx context.Context, id, sessionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM chat_messages WHERE id = $1 AND session_id = $2`, id, sessionID)
	return err
}

func (r *ChatRepository) DeleteMessagesBySession(ctx context.Context, sessionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM chat_messages WHERE session_id = $1`, sessionID)
	return err
}

func (r *ChatRepository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM chat_sessions WHERE id = $1`, sessionID)
	return err
}

func (r *ChatRepository) FindSessionByMessage(ctx context.Context, messageID string) (*model.ChatSession, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT
			cs.id,
			cs.user_id,
			cs.role_id,
			cs.model_key,
			cs.title,
			cs.mode,
			cs.status,
			COALESCE(cs.settings, cs.settings_json, '{}'::jsonb) AS settings,
			cs.created_at,
			cs.updated_at
        FROM chat_messages cm
        JOIN chat_sessions cs ON cs.id = cm.session_id
        WHERE cm.id = $1
    `, messageID)
	var session model.ChatSession
	if err := scanSession(row, &session); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *ChatRepository) CountSessions(ctx context.Context, userID string) (int, error) {
	row := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM chat_sessions WHERE user_id = $1`, userID)
	var total int
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *ChatRepository) UpdateSettings(ctx context.Context, sessionID, mode, modelKey string, settings model.ChatSessionSettings) (*model.ChatSession, error) {
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}
	query := `
		UPDATE chat_sessions
		SET mode = $2, model_key = $3, settings = $4, settings_json = $4, updated_at = now()
		WHERE id = $1
		RETURNING id, user_id, role_id, model_key, title, summary, mode, status, settings, created_at, updated_at
	`
	var s model.ChatSession
	var settingsBytes []byte
	var summary string // Added for scanning the summary field
	err = r.pool.QueryRow(ctx, query, sessionID, mode, modelKey, settingsJSON).Scan(
		&s.ID, &s.UserID, &s.RoleID, &s.ModelKey, &s.Title, &summary, &s.Mode, &s.Status, &settingsBytes, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	s.Summary = summary // Assign scanned summary
	if err := json.Unmarshal(settingsBytes, &s.Settings); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *ChatRepository) UpdateSummary(ctx context.Context, sessionID, summary string) error {
	query := `
		UPDATE chat_sessions
		SET summary = $2, updated_at = now()
		WHERE id = $1
	`
	_, err := r.pool.Exec(ctx, query, sessionID, summary)
	return err
}

func scanSession(row pgx.Row, session *model.ChatSession) error {
	var settingsRaw []byte
	if err := row.Scan(&session.ID, &session.UserID, &session.RoleID, &session.ModelKey, &session.Title, &session.Mode, &session.Status, &settingsRaw, &session.CreatedAt, &session.UpdatedAt); err != nil {
		return err
	}
	if len(settingsRaw) == 0 {
		session.Settings = model.DefaultChatSessionSettings()
		return nil
	}
	if err := json.Unmarshal(settingsRaw, &session.Settings); err != nil {
		session.Settings = model.DefaultChatSessionSettings()
	}
	return nil
}
