package main

import (
	"SuperQueue/logger"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/etcd-io/etcd/clientv3"
	"google.golang.org/grpc"
)

type SuperQueue struct {
	DelayMapMap     *MapMap
	InFlightItems   *map[string]*QueueItem
	DelayConsumer   *MapMapConsumer
	Outbox          *Outbox
	InFlightMapLock sync.RWMutex
	Namespace       string
	Partition       string
	CloseChan       chan struct{}
}

func NewSuperQueue(namespace, partition string, bucketMS, queueLen int64) *SuperQueue {
	dmm := NewMapMap(bucketMS)
	q := &SuperQueue{
		DelayMapMap:     dmm, // 5ms default
		InFlightItems:   &map[string]*QueueItem{},
		InFlightMapLock: sync.RWMutex{},
		Partition:       partition,
		CloseChan:       make(chan struct{}),
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

	// If ETCD_HOSTS exists, start reporting for service discovery
	if ETCD_HOSTS != "" {
		logger.Debug("Starting etcd based service discovery")
		hosts := strings.Split(ETCD_HOSTS, ",")
		logger.Debug("Using hosts: ", hosts)
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   hosts,
			DialTimeout: 2 * time.Second,
			DialOptions: []grpc.DialOption{grpc.WithBlock()}, // Need this to actually fail on connect
		})
		if err != nil {
			logger.Error("Failed to connect to etcd!")
			logger.Error(err)
			panic(err)
		} else {
			logger.Debug("Connected to etcd")
		}
		go func() {
			<-q.CloseChan
			logger.Info("Closing service discovery ticker")
		}()
		defer cli.Close()
	}

	return q
}

func (sq *SuperQueue) Close() {
	logger.Info("Closing superqueue")
	sq.DelayConsumer.Stop()
	close(sq.CloseChan)
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
	// Update metrics
	atomic.AddInt64(&InFlightMessages, 1)
	atomic.AddInt64(&TotalInFlightMessages, 1)
	atomic.AddInt64(&QueuedMessages, -1)
	// Return
	return item, nil
}
