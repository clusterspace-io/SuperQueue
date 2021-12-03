package main

import "time"

type ItemID string

type QueueItem struct {
	ID   ItemID
	Body []byte
}

type QueueItemDB struct {
	ID                string
	Payload           []byte
	Bucket            string
	CreatedAt         time.Time
	ExpireAt          time.Time
	InFlightTimeout   int
	BackoffMin        int
	BackoffMultiplier float64
}

type QueueItemStateDB struct {
	ID string
	// SERIAL monotomically incrementing integer
	Version int
	// Item state, ENUM ('queued', 'in-flight', 'delivered', 'discarded', 'delayed', 'timedout', 'nacked', 'discarded', 'expired')
	State     string
	CreatedAt time.Time
	DelayTo   time.Time
	Attempts  int
	// The error type, ENUM ('max retries exceeded', 'unknown', 'expired', 'nack')
	Error        string
	ErrorMessage string
}

// Adds a new queue item to the DB and immediately queues it
func (i *QueueItem) EnqueueItem() error {
	// Add item to DB
	// Add a newly created queue item to the queue
	return nil
}

// Adds a new queue item to the DB and delays it
func (i *QueueItem) DelayEnqueueItem(delayMS int64) error {
	// Add item to DB as delayed
	// Add item to delay mapmap
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

// Moving from `delayed`, `nacked`, or `timedout` to `queued` (anything in the delayed mapmap)
func (i *QueueItem) RequeueItem() {
	// Write queued to DB
	// Remove item from delayed mapmap
	// Add item to queue
}

// -----------------------------------------------------------------------------
// Internal functions
// -----------------------------------------------------------------------------

func (i *QueueItem) addItemToItemsTable() error {

	return nil
}

func (i *QueueItem) discardItem() error {

	return nil
}

func debugReadItem(itemID ItemID) *QueueItemDB {

}

func debugReadItemState(itemID ItemID) *QueueItemStateDB {

}
