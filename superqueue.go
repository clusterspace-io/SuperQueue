package main

import (
	"SuperQueue/logger"
	"time"
)

type SuperQueue struct {
	DelayMapMap   *MapMap
	InFlightItems *map[string]QueueItem
	DelayConsumer *MapMapConsumer
	Outbox        *Outbox
}

func NewSuperQueue(bucketMS, memoryQueueLen int64) *SuperQueue {
	dmm := NewMapMap(bucketMS)
	q := &SuperQueue{
		DelayMapMap:   dmm, // 5ms default
		InFlightItems: &map[string]QueueItem{},
	}

	q.Outbox = NewOutbox(q, memoryQueueLen)

	// Self reference trick
	dc := &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(dmm.Interval) * time.Millisecond),
		endChan:     make(chan struct{}),
		lastConsume: time.Now().UnixMilli(),
		MapMap:      dmm,
		ConsumerFunc: func(bucket int64, m map[string]*QueueItem) {
			logger.Debug("Consuming bucket ", bucket)
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
	logger.Debug("Enqueueing item ", item.ID)
	if delayMS > 0 {
		// If delayed, put in mapmap
		// TODO: Add to DB
		err := item.addItemToItemsTable()
		if err != nil {
			panic(err)
		}
		sq.DelayMapMap.AddItem(item, time.Now().UnixMilli()+delayMS)
	} else {
		// Otherwise put it right in outbox
		// TODO: Add to DB
	}
}
