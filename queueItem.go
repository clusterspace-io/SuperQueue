package main

import (
	"SuperQueue/logger"
	"fmt"
	"sync/atomic"
	"time"
)

type QueueItem struct {
	ID                     string
	Payload                []byte
	StorageBucket          string
	CreatedAt              time.Time
	ExpireAt               time.Time
	InFlightTimeoutSeconds int
	// Used to determine what to do in delaymapmap parsing
	InFlight bool
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
	full := sq.Outbox.Add(i)
	if full {
		atomic.AddInt64(&FullQueueResponses, 1)
		FullQueueResponsesCounter.Inc()
		return fmt.Errorf("outbox full")
	} else {
		QueuedMessagesMetric.Inc()
	}
	// atomic.AddInt64(&TotalQueuedMessages, 1)
	return nil
}

func (i *QueueItem) ReEnqueueItem(sq *SuperQueue, timedout bool, delayMS *int64) error {
	if i.InFlight {
		atomic.AddInt64(&InFlightMessages, -1)
		InFlightMessagesMetric.Dec()
		if timedout {
			timeoutMSG := "timedout"
			err := i.addItemState(sq.Name, "timedout", time.Now(), nil, &timeoutMSG, &timeoutMSG)
			if err != nil {
				logger.Error("Error adding item state during timeout:")
				logger.Error(err)
				return err
			}
			atomic.AddInt64(&TimedoutMessages, 1)
			TimedoutMessagesMetric.Inc()
		}
		i.InFlight = false
		sq.InFlightMapLock.Lock()
		delete(*sq.InFlightItems, i.ID)
		sq.InFlightMapLock.Unlock()
	} else {
		// Otherwise we are delayed
		atomic.AddInt64(&DelayedMessages, -1)
		DelayedMessagesMetric.Dec()
	}
	if delayMS != nil {
		// If delayed, put in mapmap
		dt := time.Now().Add(time.Millisecond * time.Duration(*delayMS))
		delayTime := &dt
		err := i.addItemState(sq.Name, "delayed", time.Now(), delayTime, nil, nil)
		if err != nil {
			logger.Error("Error inserting delayed item state into table on Enqueue:")
			logger.Error(err)
			return err
		}
		return i.DelayEnqueueItem(sq, *delayTime)
	} else {
		// Write queued state to DB
		err := i.addItemState(sq.Name, "queued", time.Now(), nil, nil, nil)
		if err != nil {
			logger.Error("Error adding item state during requeue:")
			logger.Error(err)
			return err
		}
		return i.EnqueueItem(sq)
	}
}

// Adds a new queue item to the DB and delays it
func (i *QueueItem) DelayEnqueueItem(sq *SuperQueue, delayTime time.Time) error {
	// Add item to DB as delayed
	// Add item to delay mapmap
	sq.DelayMapMap.AddItem(i, delayTime.UnixMilli())
	atomic.AddInt64(&DelayedMessages, 1)
	DelayedMessagesMetric.Inc()
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
	i.addItemState(sq.Name, "acked", time.Now(), nil, nil, nil)
	// Remove from inflight table
	sq.InFlightMapLock.Lock()
	delete(*sq.InFlightItems, i.ID)
	sq.InFlightMapLock.Unlock()
	// Remove from delay mapmap
	sq.DelayMapMap.DeleteItem(i)
	atomic.AddInt64(&AckedMessages, 1)
	AckedMessagesCounter.Inc()
	atomic.AddInt64(&InFlightMessages, -1)
	InFlightMessagesMetric.Dec()
	return nil
}

func (i *QueueItem) NackItem(sq *SuperQueue, delayMS *int64) error {
	// Write nack to DB
	nackMSG := "nacked"
	err := i.addItemState(sq.Name, "nacked", time.Now(), nil, &nackMSG, &nackMSG)
	if err != nil {
		return err
	}
	// Remove from inflight table
	sq.InFlightMapLock.Lock()
	delete(*sq.InFlightItems, i.ID)
	sq.InFlightMapLock.Unlock()
	// Remove from old spot in delayed mapmap
	sq.DelayMapMap.DeleteItem(i)
	// Add to new spot in delayed mapmap
	// Check whether we are delayed
	i.ReEnqueueItem(sq, false, delayMS)
	atomic.AddInt64(&NackedMessages, 1)
	NackedMessagesCounter.Inc()
	return nil
}

// -----------------------------------------------------------------------------
// Internal functions
// -----------------------------------------------------------------------------

func (i *QueueItem) addItemToItemsTable(namespace string) error {
	// _, err := PGPool.Exec(context.Background(), `
	// 	INSERT INTO items (id, payload, bucket, created_at, expire_at, in_flight_timeout)
	// 	VALUES ($1, $2, $3, $4, $5, $6)
	// `, i.ID, i.Payload, i.StorageBucket, i.CreatedAt, i.ExpireAt, i.InFlightTimeoutSeconds)
	// time.Sleep(time.Millisecond * 3)
	// return nil
	q := DBSession.Query(ItemsTable.Insert()).BindMap(map[string]interface{}{
		"namespace":         namespace,
		"id":                i.ID,
		"payload":           i.Payload,
		"bucket":            i.StorageBucket,
		"created_at":        i.CreatedAt,
		"expire_at":         i.ExpireAt,
		"in_flight_timeout": i.InFlightTimeoutSeconds,
	})
	err := q.ExecRelease()
	return err
}

func (i *QueueItem) addItemState(namespace, state string, createdAt time.Time, delayTo *time.Time, itemError, itemErrorMessage *string) error {
	i.Version++
	// _, err := PGPool.Exec(context.Background(), `
	// 	INSERT INTO item_states (id, version, state, created_at, attempts, delay_to, error, error_message)
	// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	// `, i.ID, i.Version, state, createdAt, i.Attempts, delayTo, itemError, itemErrorMessage)
	// time.Sleep(time.Millisecond * 3)
	// return nil
	q := DBSession.Query(ItemStatesTable.Insert()).BindMap(map[string]interface{}{
		"namespace":     namespace,
		"id":            i.ID,
		"version":       i.Version,
		"state":         state,
		"created_at":    createdAt,
		"attempts":      i.Attempts,
		"delay_to":      delayTo,
		"error":         itemError,
		"error_message": itemErrorMessage,
	})
	err := q.ExecRelease()
	return err
}

func (i *QueueItem) discardItem() error {

	return nil
}
