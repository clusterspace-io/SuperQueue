package main

import (
	"fmt"
	"time"
)

type SuperQueue struct {
	DelayMapMap   *MapMap
	InFlightItems *map[ItemID]QueueItem
	Outbox        chan *QueueItem
	DelayConsumer *MapMapConsumer
}

func NewSuperQueue() *SuperQueue {
	dmm := NewMapMap(5)
	q := &SuperQueue{
		DelayMapMap:   dmm, // 5ms default
		InFlightItems: &map[ItemID]QueueItem{},
		Outbox:        make(chan *QueueItem),
	}

	// Self reference trick
	dc := &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(dmm.Interval) * time.Millisecond),
		endChan:     make(chan struct{}),
		lastConsume: time.Now().UnixMilli(),
		MapMap:      dmm,
		ConsumerFunc: func(bucket int64, m map[ItemID]*QueueItem) {
			fmt.Println("consuming bucket", bucket, "in range", dmm.CalculateBucket(q.DelayConsumer.lastConsume), "through", dmm.CalculateBucket(time.Now().UnixMilli()))
		},
	}
	q.DelayConsumer = dc

	return q
}

func (sq *SuperQueue) Enqueue(item *QueueItem, delayMS int64) {
	sq.DelayMapMap.AddItem(item, time.Now().UnixMilli()+delayMS)
}
