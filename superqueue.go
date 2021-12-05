package main

import (
	"SuperQueue/logger"
	"sync"
	"time"
)

type SuperQueue struct {
	DelayMapMap     *MapMap
	InFlightItems   *map[string]*QueueItem
	DelayConsumer   *MapMapConsumer
	Outbox          *Outbox
	InFlightMapLock sync.RWMutex
	// The combination of the queue name and partition
	Namespace string
}

func NewSuperQueue(namespace string, bucketMS, queueLen int64) *SuperQueue {
	dmm := NewMapMap(bucketMS)
	q := &SuperQueue{
		DelayMapMap:     dmm, // 5ms default
		InFlightItems:   &map[string]*QueueItem{},
		InFlightMapLock: sync.RWMutex{},
	}

	q.Outbox = NewOutbox(q, queueLen)

	// Self reference trick
	dc := &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(dmm.Interval) * time.Millisecond),
		endChan:     make(chan struct{}),
		lastConsume: time.Now().UnixMilli(),
		MapMap:      dmm,
		ConsumerFunc: func(bucket int64, m map[string]*QueueItem) {
			// logger.Debug("Consuming bucket ", bucket)
			for _, i := range m {
				// logger.Debug("Found item: ", i)
				// In testing this as a goroutine had no difference in processing speed
				i.ReEnqueueItem(SQ)
			}
			// logger.Debug("Deleting bucket ", bucket)
			dmm.m.Lock()
			delete(dmm.Map, bucket)
			dmm.m.Unlock()
		},
	}
	q.DelayConsumer = dc

	return q
}

func (sq *SuperQueue) Enqueue(item *QueueItem, delayTime *time.Time) error {
	logger.Debug("Enqueueing item ", item.ID)
	err := item.addItemToItemsTable(sq.Namespace)
	if err != nil {
		logger.Error("Error inserting item into table on Enqueue:")
		logger.Error(err)
		return err
	}

	if delayTime != nil {
		// If delayed, put in mapmap
		err = item.addItemState(sq.Namespace, "delayed", item.CreatedAt, delayTime, nil, nil)
		if err != nil {
			logger.Error("Error inserting delayed item state into table on Enqueue:")
			logger.Error(err)
			return err
		}
		item.DelayEnqueueItem(sq, *delayTime)
	} else {
		// Otherwise put it right in outbox
		err = item.addItemState(sq.Namespace, "queued", item.CreatedAt, nil, nil, nil)
		if err != nil {
			logger.Error("Error inserting item state into table on Enqueue:")
			logger.Error(err)
			return err
		}
		item.EnqueueItem(sq)
	}

	return nil
}

func (sq *SuperQueue) Dequeue() (*QueueItem, error) {
	// Get from the outbox
	item := sq.Outbox.Pop()
	// Empty
	if item == nil {
		return nil, nil
	}
	// Increment delivery attempts
	item.Attempts++
	// Write inflight state to db
	item.addItemState(sq.Namespace, "in-flight", time.Now(), nil, nil, nil)
	item.InFlight = true
	// Put in in-flight map with in-flight timeout
	sq.InFlightMapLock.Lock()
	(*sq.InFlightItems)[item.ID] = item
	sq.InFlightMapLock.Unlock()
	// Add to delay map
	sq.DelayMapMap.AddItem(item, time.Now().Add(time.Duration(item.InFlightTimeoutSeconds)*time.Second).UnixMilli())
	// Return
	return item, nil
}
