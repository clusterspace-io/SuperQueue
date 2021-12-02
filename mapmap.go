package main

import (
	"math"
	"sync"
	"time"
)

type MapMap struct {
	Map map[int64]map[ItemID]*QueueItem

	m sync.Mutex

	// The bucketing interval
	Interval int64
}

func NewMapMap(interval int64) *MapMap {
	return &MapMap{
		Interval: interval,
		Map:      map[int64]map[ItemID]*QueueItem{},
	}
}

// For every bucket between the lower and upper (inclusive) bounds, run a function on the resulting map
func (m *MapMap) ConsumeRange(lowerbound, upperbound time.Time, consumerFunc func(map[ItemID]*QueueItem)) {
	for i := lowerbound.UnixMilli(); i < upperbound.UnixMilli(); i += 5 {
		// Calculate the bucket to get
		bkey := m.calculateBucket(i)
		consumerFunc(m.Map[bkey])
	}
}

// Adds a new item, creating the bucket if needed (thread safe). `bucketer` should reflect the bucket in which you want this item consumed
func (m *MapMap) AddItem(item *QueueItem, future time.Time) {
	bucket := m.calculateBucket(future.UnixMilli())
	if _, e := m.Map[bucket]; !e {
		m.m.Lock()
		defer m.m.Unlock()
		m.Map[bucket] = map[ItemID]*QueueItem{}
	}
	m.Map[bucket][item.ID] = item
}

// Calculates what bucket a bucketer should be in
func (m *MapMap) calculateBucket(bucketer int64) int64 {
	return int64(math.Round(float64(bucketer)/float64(m.Interval))) * m.Interval
}
