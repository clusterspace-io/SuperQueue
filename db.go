package main

import (
	"SuperQueue/logger"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	PGPool *pgxpool.Pool
)

func ConnectToDB(connString string) error {
	logger.Debug("Connecting to db...")
	var err error
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return err
	}
	PGPool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return err
	}
	logger.Debug("Connected to db")
	return nil
}

func CreateTables() error {
	logger.Debug("Creating tables...")
	statement1 := `
	CREATE TABLE IF NOT EXISTS items (
		id TEXT NOT NULL,
		payload BLOB NOT NULL,
		bucket TEXT NOT NULL, -- the bucket to archive to, maybe remove if we keep this config somewhere else
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		expire_at TIMESTAMPTZ NOT NULL, -- when we say this message will be deleted
		in_flight_timeout INT NOT NULL,
		backoff_min INT NOT NULL, -- the first backoff value
		backoff_multiplier REAL NOT NULL, -- the exponential multiplier, should be >= 1.0 (enforced with code)
		PRIMARY KEY(id)
	);
	`
	_, err := PGPool.Query(context.Background(), statement1)
	if err != nil {
		logger.Error("Error making the items table")
		return err
	}
	logger.Debug("Created items table")

	statement2 := `
	CREATE TYPE item_state AS ENUM ('awaiting-queueing', 'queued', 'in-flight', 'delivered', 'failed', 'delayed', 'archiving', 'archived', 'expired');
	`
	_, err = PGPool.Query(context.Background(), statement2)
	if err != nil {
		logger.Error("Error making the item_state ENUM")
		return err
	}
	logger.Debug("Created item_state enum")

	statement3 := `
	CREATE TYPE delivery_error AS ENUM ('max retries exceeded', 'unknown', 'expired', 'nack');
	`
	_, err = PGPool.Query(context.Background(), statement3)
	if err != nil {
		logger.Error("Error making the delivery_error ENUM")
		return err
	}
	logger.Debug("Created delivery_error enum")

	statement4 := `
	CREATE TABLE IF NOT EXISTS item_states (
		id TEXT NOT NULL,
		generation TEXT NOT NULL,
		state item_state NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		attempts INT NOT NULL,
		delay_to TIMESTAMPTZ, -- for either delayed messages, nacked, or timed out retries
		error delivery_error,
		error_message TEXT,
		PRIMARY KEY(id, generation)
	);
	`
	_, err = PGPool.Query(context.Background(), statement4)
	if err != nil {
		logger.Error("Error making the item_states table")
		return err
	}
	logger.Debug("Created item_states enum")
	return nil
}
