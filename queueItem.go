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
	StorageBucket          string
	CreatedAt              time.Time
	ExpireAt               time.Time
	InFlightTimeoutSeconds int
	// Used to determine what to do in delaymapmap parsing
	InFlight          bool
	BackoffMinMS      int
	BackoffMultiplier float64
	// The time that is used to bucket for the delaymapmap
	TimeBucket int64

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

// Moving from `delayed`
func (i *QueueItem) ReEnqueueItem(sq *SuperQueue) error {
	// If in-flight then mark timedout
	if i.InFlight {
		timeoutMSG := "timedout"
		err := i.addItemState("timedout", time.Now(), nil, &timeoutMSG, &timeoutMSG)
		if err != nil {
			logger.Error("Error adding item state during timeout:")
			logger.Error(err)
			return err
		}
		i.InFlight = false
		delete(*sq.InFlightItems, i.ID)
	}
	// Write queued state to DB
	err := i.addItemState("queued", time.Now(), nil, nil, nil)
	if err != nil {
		logger.Error("Error adding item state during requeue:")
		logger.Error(err)
		return err
	}
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

func (i *QueueItem) AckItem(sq *SuperQueue) error {
	// Write ack to DB
	i.addItemState("acked", time.Now(), nil, nil, nil)
	// Remove from inflight table
	delete(*sq.InFlightItems, i.ID)
	// Remove from delay mapmap
	sq.DelayMapMap.DeleteItem(i)
	return nil
}

func (i *QueueItem) NackItem(sq *SuperQueue) error {
	// Write nack to DB
	nackMSG := "nacked"
	err := i.addItemState("nacked", time.Now(), nil, &nackMSG, &nackMSG)
	if err != nil {
		return err
	}
	// Remove from inflight table
	delete(*sq.InFlightItems, i.ID)
	// Remove from old spot in delayed mapmap
	sq.DelayMapMap.DeleteItem(i)
	// TODO: Discard if max attempts exceeded
	// Add to new spot in delayed mapmap
	i.InFlight = false
	i.ReEnqueueItem(sq)
	return nil
}

// -----------------------------------------------------------------------------
// Internal functions
// -----------------------------------------------------------------------------

func (i *QueueItem) addItemToItemsTable() error {
	_, err := PGPool.Exec(context.Background(), `
		INSERT INTO items (id, payload, bucket, created_at, expire_at, in_flight_timeout, backoff_min, backoff_multiplier)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, i.ID, i.Payload, i.StorageBucket, i.CreatedAt, i.ExpireAt, i.InFlightTimeoutSeconds, i.BackoffMinMS, i.BackoffMultiplier)
	return err
}

func (i *QueueItem) addItemState(state string, createdAt time.Time, delayTo *time.Time, itemError, itemErrorMessage *string) error {
	i.Version++
	_, err := PGPool.Exec(context.Background(), `
		INSERT INTO item_states (id, version, state, created_at, attempts, delay_to, error, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, i.ID, i.Version, state, createdAt, i.Attempts, delayTo, itemError, itemErrorMessage)
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
