package main

import (
	"SuperQueue/logger"
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type QueueItem struct {
	ID                     string
	Payload                []byte
	Bucket                 string
	CreatedAt              time.Time
	ExpireAt               time.Time
	InFlightTimeoutSeconds int
	BackoffMinMS           int
	BackoffMultiplier      float64

	// Not stored in the DB
	Attempts int
	Version  int
}

type QueueItemStateDB struct {
	ID string
	// SERIAL monotomically incrementing integer
	Version int
	// Item state, ENUM ('queued', 'in-flight', 'delivered', 'discarded', 'delayed', 'timedout', 'nacked', 'discarded', 'expired')
	State     string
	CreatedAt time.Time
	DelayTo   *time.Time
	Attempts  int
	// The error type, ENUM ('max retries exceeded', 'unknown', 'expired', 'nack')
	Error        *string
	ErrorMessage *string
}

// Adds a new queue item to the DB and immediately queues it
func (i *QueueItem) EnqueueItem(sq *SuperQueue) error {
	// Add item to the queue
	sq.Outbox.Add(i)
	return nil
}

// Moving from `delayed`, `nacked`, or `timedout` to `queued` (anything in the delayed mapmap)
func (i *QueueItem) RequeueItem(sq *SuperQueue) error {
	// Write queued state to DB
	err := i.addItemState("queued", i.CreatedAt, i.Attempts, nil, nil, nil)
	if err != nil {
		logger.Error("Error adding item state during requeue:")
		logger.Error(err)
		return err
	}
	// Remove item from delayed mapmap
	return i.EnqueueItem(sq)
}

// Adds a new queue item to the DB and delays it
func (i *QueueItem) DelayEnqueueItem(sq *SuperQueue, delayTime time.Time) error {
	// Add item to DB as delayed
	// Add item to delay mapmap
	sq.DelayMapMap.AddItem(i, delayTime.UnixMilli())
	return nil
}

func (i *QueueItem) DequeueItem() error {
	// Write inflight to DB
	// Move to inflight table
	// Move to mapmap for timeout capture
	return nil
}

func (i *QueueItem) AckItem() error {
	// Write ack to DB
	// Remove from inflight table
	return nil
}

func (i *QueueItem) NackItem() error {
	// Write nack to DB
	// Remove from inflight table
	// Remove from old spot in delayed mapmap
	// Discard if max attempts exceeded
	// Add to new spot in delayed mapmap
	return nil
}

func (i *QueueItem) TimeoutItem() error {
	// Write timeout to DB
	// Remove from inflight table
	// Discard if max attempts exceeded
	// Move to delayed mapmap
	return nil
}

// -----------------------------------------------------------------------------
// Internal functions
// -----------------------------------------------------------------------------

func (i *QueueItem) addItemToItemsTable() error {
	_, err := PGPool.Exec(context.Background(), `
		INSERT INTO items (id, payload, bucket, created_at, expire_at, in_flight_timeout, backoff_min, backoff_multiplier)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, i.ID, i.Payload, i.Bucket, i.CreatedAt, i.ExpireAt, i.InFlightTimeoutSeconds, i.BackoffMinMS, i.BackoffMultiplier)
	return err
}

func (i *QueueItem) addItemState(state string, createdAt time.Time, attempts int, delayTo *time.Time, itemError, itemErrorMessage *string) error {
	i.Version++
	_, err := PGPool.Exec(context.Background(), `
		INSERT INTO item_states (id, version, state, created_at, attempts, delay_to, error, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, i.ID, i.Version, state, createdAt, attempts, delayTo, itemError, itemErrorMessage)
	return err
}

func (i *QueueItem) discardItem() error {

	return nil
}

func debugReadItem(itemID string) (*QueueItem, error) {
	var item QueueItem
	rows, err := PGPool.Query(context.Background(), `
		SELECT *
		FROM items
		WHERE id = $1
	`, itemID)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanOne(&item, rows)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func debugReadItemState(itemID string) (*QueueItemStateDB, error) {
	var item QueueItemStateDB
	rows, err := PGPool.Query(context.Background(), `
		SELECT *
		FROM item_states
		WHERE id = $1
	`, itemID)
	if err != nil {
		return nil, err
	}
	err = pgxscan.ScanOne(&item, rows)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
