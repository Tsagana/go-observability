package job

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storer interface {
	Claim(ctx context.Context) (*Job, error)
	Complete(ctx context.Context, id string, result []byte) error
	Fail(ctx context.Context, job *Job) (*Job, error)
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, payload []byte) (*Job, error) {
	query := `
        INSERT INTO jobs (payload)
        VALUES ($1)
        RETURNING id, status, payload, result, error, retry_count, retry_after, created_at, updated_at
    `
	var j Job

	err := s.db.QueryRow(ctx, query, payload).Scan(
		&j.ID, &j.Status, &j.Payload, &j.Result, &j.Error, &j.RetryCount, &j.RetryAfter, &j.CreatedAt, &j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (s *Store) Get(ctx context.Context, id string) (*Job, error) {
	query := `
		SELECT id, status, payload, result, error, retry_count, retry_after, created_at, updated_at
		FROM jobs WHERE id = $1
    `
	var j Job

	err := s.db.QueryRow(ctx, query, id).Scan(
		&j.ID, &j.Status, &j.Payload, &j.Result, &j.Error, &j.RetryCount, &j.RetryAfter, &j.CreatedAt, &j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (s *Store) Claim(ctx context.Context) (*Job, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	var j Job
	query := `
		SELECT id, status, payload, result, error, retry_count, retry_after, created_at, updated_at
		FROM jobs WHERE status = 'pending' 
		AND (retry_after IS NULL OR retry_after < NOW())
		ORDER BY created_at
		LIMIT 1
		FOR UPDATE SKIP LOCKED;
    `

	err = tx.QueryRow(ctx, query).Scan(
		&j.ID, &j.Status, &j.Payload, &j.Result, &j.Error, &j.RetryCount, &j.RetryAfter, &j.CreatedAt, &j.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	updateQuery := `
		UPDATE jobs SET status = 'processing', updated_at = NOW()
		WHERE id=$1
    `
	_, err = tx.Exec(ctx, updateQuery, j.ID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &j, nil
}

func (s *Store) Complete(ctx context.Context, id string, result []byte) error {
	updateQuery := `
		UPDATE jobs SET status = 'completed', result = $1, updated_at = NOW()
		WHERE id=$2
    `
	_, updateErr := s.db.Exec(ctx, updateQuery, result, id)
	if updateErr != nil {
		return updateErr
	}
	return nil
}

func (s *Store) Fail(ctx context.Context, job *Job) (*Job, error) {
	nextRetryCnt := job.RetryCount + 1
	if nextRetryCnt < MaxRetry {
		backoff := calculateBackoff(nextRetryCnt)
		updateQuery := `
		UPDATE jobs SET status = 'pending', retry_count = $1, retry_after = NOW() + ($2 * interval '1 second'), updated_at = NOW()
		WHERE id=$3
    `
		_, updateErr := s.db.Exec(ctx, updateQuery, nextRetryCnt, backoff, job.ID)
		if updateErr != nil {
			return nil, updateErr
		}
		return job, nil
	}

	updateQuery := `
		UPDATE jobs SET status = 'failed', retry_count = $1, updated_at = NOW()
		WHERE id=$2
    `
	_, updateErr := s.db.Exec(ctx, updateQuery, nextRetryCnt, job.ID)
	if updateErr != nil {
		return nil, updateErr
	}
	log.Println("Job failed ", job.ID)
	return job, nil

}
