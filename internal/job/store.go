package job

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

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
