CREATE TABLE IF NOT EXISTS items (
  id TEXT NOT NULL,
  payload BLOB NOT NULL,
  bucket TEXT NOT NULL, -- the bucket to archive to, maybe remove if we keep this config somewhere else
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expire_at TIMESTAMPTZ NOT NULL, -- when we say this message will be deleted
  in_flight_timeout INT NOT NULL, -- in seconds
  backoff_min INT NOT NULL, -- the first backoff value in milliseconds
  backoff_multiplier REAL NOT NULL, -- the exponential multiplier, should be >= 1.0 (enforced with code)
  PRIMARY KEY(id)
);

CREATE TYPE item_state AS ENUM ('queued', 'in-flight', 'delivered', 'discarded', 'delayed', 'timedout', 'nacked', 'expired');

CREATE TYPE delivery_error AS ENUM ('max retries exceeded', 'unknown', 'timedout', 'expired', 'nack');

CREATE TABLE IF NOT EXISTS item_states (
  id TEXT NOT NULL,
  version INT NOT NULL,
  state item_state NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  attempts INT NOT NULL,
  delay_to TIMESTAMPTZ, -- for either delayed messages, nacked, or timed out retries
  error delivery_error,
  error_message TEXT,
  PRIMARY KEY(id, version)
);
