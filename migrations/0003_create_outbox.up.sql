-- outbox stores domain events atomically alongside job writes.
-- The publisher process reads from here and forwards events to Redis.
-- Named "outbox" after the transactional outbox pattern, not a mailbox.
CREATE TABLE outbox (
    id           BIGSERIAL PRIMARY KEY,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    published_at TIMESTAMPTZ
);

CREATE INDEX outbox_unpublished_idx ON outbox (id) WHERE published_at IS NULL;
