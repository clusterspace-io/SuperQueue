package main

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

var (
	MAPMAP_SCALE    = 10000000
	MAPMAP_RANGE    = 1000000
	MAPMAP_INTERVAL = 5
	CONSUME_SEC     = 10
)

func TestMapMapScale(t *testing.T) {
	fmt.Println("\n## Testing MapMap Scale with", MAPMAP_SCALE, "items and an interval of", MAPMAP_INTERVAL)
	m := NewMapMap(int64(MAPMAP_INTERVAL))
	fmt.Println("Adding", MAPMAP_SCALE, "items...")
	s := time.Now()
	for i := 0; i < MAPMAP_SCALE; i++ {
		m.AddItem(&QueueItem{
			ID: ItemID(fmt.Sprintf("%d", i)),
		}, time.Now().Add(time.Duration(i)*time.Millisecond).UnixMilli())
	}
	fmt.Println("MapMap filled", MAPMAP_SCALE, "items in", time.Since(s))

	fmt.Println("\nTesting range consumption of", MAPMAP_RANGE, "items")
	s = time.Now()
	m.ConsumeRange(time.Now().UnixMilli(), time.Now().Add(time.Duration(MAPMAP_RANGE)*time.Millisecond).UnixMilli(), func(bkey int64, mi map[ItemID]*QueueItem) {
		for _, _ = range mi {
			// spin
		}
	})

	fmt.Println("consumed in", time.Since(s))

	fmt.Println("\nCreating consumer and running for", CONSUME_SEC, "seconds")
	endChan := make(chan struct{})
	var Dc *MapMapConsumer
	var agg int64 = 0
	Dc = &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(m.Interval) * time.Millisecond),
		endChan:     endChan,
		lastConsume: time.Now().UnixMilli(),
		MapMap:      m,
		ConsumerFunc: func(bucket int64, mi map[ItemID]*QueueItem) {
			atomic.AddInt64(&agg, 1)
		},
	}
	Dc.Start()

	time.Sleep(time.Duration(CONSUME_SEC) * time.Second)
	close(endChan)
	fmt.Println("Done consuming", agg, "buckets")
	if agg != 2000 {
		HandleTestError(t, fmt.Errorf("Did not consume the expected 2000 buckets!"))
	}
}
