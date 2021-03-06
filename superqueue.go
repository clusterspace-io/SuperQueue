package main

import (
	"SuperQueue/logger"
	"sync"
	"sync/atomic"
	"time"
)

type SuperQueue struct {
	DelayMapMap     *MapMap
	InFlightItems   *map[string]*QueueItem
	DelayConsumer   *MapMapConsumer
	Outbox          *Outbox
	InFlightMapLock sync.RWMutex
	Name            string
	Partition       string
	CloseChan       chan struct{}
}

func NewSuperQueue(queueName, partition string, bucketMS, queueLen int64) *SuperQueue {
	dmm := NewMapMap(bucketMS)
	q := &SuperQueue{
		DelayMapMap:     dmm, // 5ms default
		InFlightItems:   &map[string]*QueueItem{},
		InFlightMapLock: sync.RWMutex{},
		Partition:       partition,
		CloseChan:       make(chan struct{}),
		Name:            queueName,
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
				// TODO: Decrement delayed message if needed
				i.ReEnqueueItem(SQ, true, nil) // Not sure I like calling this `timedout`, but it kind of is (the delay timed out) and it really only gets considered if it was inflight anyway
			}
			// logger.Debug("Deleting bucket ", bucket)
		},
	}
	q.DelayConsumer = dc

	return q
}

func (sq *SuperQueue) Close() {
	logger.Info("Closing superqueue")
	sq.DelayConsumer.Stop()
	if EtcdClient != nil {
		logger.Info("Closing etcd client")
		EtcdClient.Close()
	}
	close(sq.CloseChan)
}

func (sq *SuperQueue) Enqueue(item *QueueItem, delayTime *time.Time) error {
	logger.Debug("Enqueueing item ", item.ID)
	err := item.addItemToItemsTable(sq.Name)
	if err != nil {
		logger.Error("Error inserting item into table on Enqueue:")
		logger.Error(err)
		return err
	}

	if delayTime != nil {
		// If delayed, put in mapmap
		err = item.addItemState(sq.Name, "delayed", item.CreatedAt, delayTime, nil, nil)
		if err != nil {
			logger.Error("Error inserting delayed item state into table on Enqueue:")
			logger.Error(err)
			return err
		}
		item.DelayEnqueueItem(sq, *delayTime)
	} else {
		// Otherwise put it right in outbox
		err = item.addItemState(sq.Name, "queued", item.CreatedAt, nil, nil, nil)
		if err != nil {
			logger.Error("Error inserting item state into table on Enqueue:")
			logger.Error(err)
			return err
		}
		return item.EnqueueItem(sq)
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
	item.addItemState(sq.Name, "in-flight", time.Now(), nil, nil, nil)
	item.InFlight = true
	// Put in in-flight map with in-flight timeout
	sq.InFlightMapLock.Lock()
	(*sq.InFlightItems)[item.ID] = item
	sq.InFlightMapLock.Unlock()
	// Add to delay map
	sq.DelayMapMap.AddItem(item, time.Now().Add(time.Duration(item.InFlightTimeoutSeconds)*time.Second).UnixMilli())
	// Update metrics
	atomic.AddInt64(&InFlightMessages, 1)
	InFlightMessagesMetric.Inc()
	// atomic.AddInt64(&TotalInFlightMessages, 1)
	atomic.AddInt64(&QueuedMessages, -1)
	QueuedMessagesMetric.Dec()
	// Return
	return item, nil
}
