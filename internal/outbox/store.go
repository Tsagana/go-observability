package outbox

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID          int64
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
	PublishedAt *time.Time
}

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

// InsertTx writes an outbox event inside an existing transaction.
// The caller owns the transaction and is responsible for commit/rollback.
func (s *Store) InsertTx(ctx context.Context, tx pgx.Tx, eventType string, payload []byte) error {
	query := `
        INSERT INTO outbox (event_type, payload)
        VALUES ($1, $2)`

	_, err := tx.Exec(ctx, query, eventType, payload)

	if err != nil {
		return err
	}

	return nil
}

// BeginTx starts a new transaction. The caller owns commit/rollback.
func (s *Store) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}

// FetchUnpublishedTx returns up to limit unpublished events within the caller's transaction.
// Uses FOR UPDATE SKIP LOCKED so multiple publisher instances are safe.
// The caller must hold the transaction open until after MarkPublishedTx and Commit.
func (s *Store) FetchUnpublishedTx(ctx context.Context, tx pgx.Tx, limit int) ([]Event, error) {
	query := `
		SELECT id, event_type, payload, created_at, published_at
		FROM outbox WHERE published_at IS NULL ORDER BY id LIMIT $1 FOR UPDATE SKIP LOCKED
	`

	rows, err := tx.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload, &e.CreatedAt, &e.PublishedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// MarkPublishedTx sets published_at = now() for the given event IDs within the caller's transaction.
func (s *Store) MarkPublishedTx(ctx context.Context, tx pgx.Tx, ids []int64) error {
	query := `
		UPDATE outbox SET published_at = NOW()
		WHERE id = ANY($1)
	`
	_, err := tx.Exec(ctx, query, ids)
	return err
}
