package main

type ItemID string

type QueueItem struct {
	ID   ItemID
	Body []byte
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
