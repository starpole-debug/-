package repository

import (
	"context"
	"time"

	"github.com/example/ai-avatar-studio/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ImageJobRepository persists image generation jobs.
type ImageJobRepository struct {
	pool *pgxpool.Pool
}

func NewImageJobRepository(pool *pgxpool.Pool) *ImageJobRepository {
	return &ImageJobRepository{pool: pool}
}

func (r *ImageJobRepository) Create(ctx context.Context, job *model.ImageJob) (*model.ImageJob, error) {
	if job.ID == "" {
		job.ID = uuid.NewString()
	}
	now := time.Now()
	job.CreatedAt = now
	job.UpdatedAt = now
	_, err := r.pool.Exec(ctx, `
        INSERT INTO image_jobs(id, user_id, session_id, message_id, provider_id, preset_id, prompt, negative_prompt, final_prompt, status, result_url, error, created_at, updated_at)
        VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
    `, job.ID, job.UserID, job.SessionID, job.MessageID, job.ProviderID, job.PresetID, job.Prompt, job.NegativePrompt, job.FinalPrompt, job.Status, job.ResultURL, job.Error, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *ImageJobRepository) UpdateStatus(ctx context.Context, id, status, resultURL, errMsg string) error {
	_, err := r.pool.Exec(ctx, `
        UPDATE image_jobs
        SET status = $2,
            result_url = $3,
            error = $4,
            updated_at = now()
        WHERE id = $1
    `, id, status, resultURL, errMsg)
	return err
}

func (r *ImageJobRepository) Find(ctx context.Context, id string) (*model.ImageJob, error) {
	row := r.pool.QueryRow(ctx, `
        SELECT id, user_id, session_id, message_id, provider_id, preset_id, prompt, negative_prompt, final_prompt, status, result_url, error, created_at, updated_at
        FROM image_jobs WHERE id = $1
    `, id)
	var job model.ImageJob
	if err := row.Scan(&job.ID, &job.UserID, &job.SessionID, &job.MessageID, &job.ProviderID, &job.PresetID, &job.Prompt, &job.NegativePrompt, &job.FinalPrompt, &job.Status, &job.ResultURL, &job.Error, &job.CreatedAt, &job.UpdatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}
