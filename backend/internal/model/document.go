package model

import "time"

// Document stores unstructured knowledge base content for retrieval.
type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
}

// DocumentChunk stores embeddings for similarity search.
type DocumentChunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Embedding  []float32 `json:"embedding"`
	CreatedAt  time.Time `json:"created_at"`
}

// MemoryCapsule stores persona-specific preferences.
type MemoryCapsule struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	RoleID    string    `json:"role_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
