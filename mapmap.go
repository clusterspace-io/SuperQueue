package main

import (
	"math"
	"sync"
)

type MapMap struct {
	Map map[int64]map[string]*QueueItem

	m sync.RWMutex

	// The bucketing interval
	Interval int64
}

func NewMapMap(intervalMS int64) *MapMap {
	return &MapMap{
		Interval: intervalMS,
		Map:      map[int64]map[string]*QueueItem{},
		m:        sync.RWMutex{},
	}
}

// For every bucket between the lower and upper (inclusive) bounds, run a function on the resulting map. Returns the last item checked, -1 if none were found
func (m *MapMap) ConsumeRange(lowerbound, upperbound int64, consumerFunc func(int64, map[string]*QueueItem)) int64 {
	var lastItem int64 = -1
	upperBucket := m.CalculateBucket(upperbound)
	// lowerBucket := m.CalculateBucket(lowerbound)
	var i int64
	for i = 0; i <= upperBucket; i += m.Interval {
		// Calculate the bucket to get
		bkey := m.CalculateBucket(i)
		m.m.Lock()
		if _, exists := m.Map[bkey]; exists {
			consumerFunc(bkey, m.Map[bkey])
		}
		delete(m.Map, bkey)
		m.m.Unlock()
		lastItem = bkey
	}
	return lastItem
}

// Adds a new item, creating the bucket if needed (thread safe). `bucketer` should reflect the bucket in which you want this item consumed
func (m *MapMap) AddItem(item *QueueItem, executeTimeMS int64) {
	item.TimeBucket = m.CalculateBucket(executeTimeMS)
	m.m.Lock()
	defer m.m.Unlock()
	if _, e := m.Map[item.TimeBucket]; !e {
		m.Map[item.TimeBucket] = map[string]*QueueItem{}
	}
	m.Map[item.TimeBucket][item.ID] = item
}

// Calculates what bucket a bucketer should be in
func (m *MapMap) CalculateBucket(bucketer int64) int64 {
	return int64(math.Floor(float64(bucketer)/float64(m.Interval))) * m.Interval
}

func (m *MapMap) DeleteItem(item *QueueItem) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.Map[item.TimeBucket], item.ID)
}
