package events

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

// MarkProcessed inserts event_id into processed_events.
// Returns (true, nil) on first insert, (false, nil) if already exists (duplicate).
func (s *Store) MarkProcessed(ctx context.Context, eventID int64) (bool, error) {
	query := `INSERT INTO processed_events (event_id) VALUES ($1) ON CONFLICT DO NOTHING`

	tag, err := s.db.Exec(ctx, query, eventID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() == 1, nil
}
