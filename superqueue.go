package main

import (
	"time"
)

type SuperQueue struct {
	DelayMapMap   *MapMap
	InFlightItems *map[ItemID]QueueItem
	DelayConsumer *MapMapConsumer
	Outbox        *Outbox
}

func NewSuperQueue(bucketMS, memoryQueueLen int64) *SuperQueue {
	dmm := NewMapMap(bucketMS)
	q := &SuperQueue{
		DelayMapMap:   dmm, // 5ms default
		InFlightItems: &map[ItemID]QueueItem{},
	}

	q.Outbox = NewOutbox(q, memoryQueueLen)

	// Self reference trick
	dc := &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(dmm.Interval) * time.Millisecond),
		endChan:     make(chan struct{}),
		lastConsume: time.Now().UnixMilli(),
		MapMap:      dmm,
		ConsumerFunc: func(bucket int64, m map[ItemID]*QueueItem) {
			for _, i := range m {
				// Move on disk

				// Put in outbox
				q.Outbox.Add(i)
			}
		},
	}
	q.DelayConsumer = dc

	return q
}

func (sq *SuperQueue) Enqueue(item *QueueItem, delayMS int64) {
	if delayMS > 0 {
		// If delayed, put in mapmap
		// TODO: Add to disk
		sq.DelayMapMap.AddItem(item, time.Now().UnixMilli()+delayMS)
	} else {
		// Otherwise put it right in outbox
		// TODO: Add to disk
	}
}

// Creates a new item in DB
func (sq *SuperQueue) addItemDB(item *QueueItem) error {

}

// Updates an item in DB
func (sq *SuperQueue) updateItemDB(item *QueueItem) error {

}

// Deletes an item in DB
func (sq *SuperQueue) deleteItemDB(item *QueueItem) error {

}
