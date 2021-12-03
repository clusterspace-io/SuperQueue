package main

type ItemID string

type QueueItem struct {
	ID   ItemID
	Body []byte
}

// Creates a struct and validates the input
func NewQueueItem(itemID string, body []byte, delayMS *int64) error {

}

// Adds a new queue item to the DB and immediately queues it
func (i *QueueItem) AddQueueItem() error {
	// Create a new
}

// Adds a new queue item to the DB and delays it
func (i *QueueItem) AddDelayedQueueItem() error {

}

func (i *QueueItem) DequeueItem() error {
	// Write inflight to DB
	// Move to inflight table
	// Move to mapmap for timeout capture
}

func (i *QueueItem) AckItem() error {
	// Write ack to DB
	// Remove from inflight table
}

func (i *QueueItem) NackItem() error {
	// Write nack to DB
	// Remove from inflight table
	// Remove from old spot in delayed mapmap
	// Discard if max attempts exceeded
	// Add to new spot in delayed mapmap
}

func (i *QueueItem) TimeoutItem() error {
	// Write timeout to DB
	// Remove from inflight table
	// Discard if max attempts exceeded
	// Move to delayed mapmap
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

}

func (i *QueueItem) discardItem() error {

}
