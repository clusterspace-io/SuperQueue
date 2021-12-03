package main

import (
	"fmt"
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

func NewMapMap(intervalMS int64) *MapMap {
	return &MapMap{
		Interval: intervalMS,
		Map:      map[int64]map[ItemID]*QueueItem{},
	}
}

// For every bucket between the lower and upper (inclusive) bounds, run a function on the resulting map. Returns the last item checked, -1 if none were found
func (m *MapMap) ConsumeRange(lowerbound, upperbound int64, consumerFunc func(int64, map[ItemID]*QueueItem)) int64 {
	var lastItem int64 = -1
	upperBucket := m.CalculateBucket(upperbound)
	lowerBucket := m.CalculateBucket(lowerbound)
	for i := lowerBucket; i <= upperBucket; i += m.Interval {
		// Calculate the bucket to get
		bkey := m.CalculateBucket(i)
		if _, exists := m.Map[bkey]; exists {
			consumerFunc(bkey, m.Map[bkey])
		}
		lastItem = bkey
	}
	return lastItem
}

// Adds a new item, creating the bucket if needed (thread safe). `bucketer` should reflect the bucket in which you want this item consumed
func (m *MapMap) AddItem(item *QueueItem, executeTime int64) {
	bucket := m.CalculateBucket(executeTime)
	fmt.Println("Adding item to bucket", bucket, "currently", m.CalculateBucket(time.Now().UnixMilli()))
	if _, e := m.Map[bucket]; !e {
		m.m.Lock()
		defer m.m.Unlock()
		m.Map[bucket] = map[ItemID]*QueueItem{}
	}
	m.Map[bucket][item.ID] = item
}

// Calculates what bucket a bucketer should be in
func (m *MapMap) CalculateBucket(bucketer int64) int64 {
	return int64(math.Floor(float64(bucketer)/float64(m.Interval))) * m.Interval
}
