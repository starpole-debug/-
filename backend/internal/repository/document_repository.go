package repository

import (
	"context"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DocumentRepository exposes read helpers for the RAG placeholder.
type DocumentRepository struct {
	pool *pgxpool.Pool
}

func NewDocumentRepository(pool *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{pool: pool}
}

// ListRecentChunks returns the most recent document snippets as a mock retrieval layer.
func (r *DocumentRepository) ListRecentChunks(ctx context.Context, limit int) ([]model.DocumentChunk, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT id, document_id, content, created_at
        FROM document_chunks ORDER BY created_at DESC LIMIT $1
    `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chunks []model.DocumentChunk
	for rows.Next() {
		var chunk model.DocumentChunk
		if err := rows.Scan(&chunk.ID, &chunk.DocumentID, &chunk.Content, &chunk.CreatedAt); err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}
