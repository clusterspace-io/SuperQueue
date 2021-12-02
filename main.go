package main

import (
	"SuperQueue/logger"
	"fmt"
	"time"
)

var (
	DelayMapMap   *MapMap
	InFlightItems *map[ItemID]QueueItem
	Outbox        chan *QueueItem
	DelayTicker   time.Ticker
	DelayConsumer *MapMapConsumer
)

func main() {
	logger.Info("Starting SuperQueue")
	DelayMapMap = NewMapMap(5)
	endChan := make(chan struct{})
	DelayConsumer = &MapMapConsumer{
		ticker:      *time.NewTicker(time.Duration(DelayMapMap.Interval) * time.Millisecond),
		endChan:     endChan,
		lastConsume: time.Now().UnixMilli(),
		MapMap:      DelayMapMap,
		ConsumerFunc: func(bucket int64, m map[ItemID]*QueueItem) {
			fmt.Println("consuming bucket", bucket, "in range", DelayMapMap.CalculateBucket(DelayConsumer.lastConsume), "through", DelayMapMap.CalculateBucket(time.Now().UnixMilli()))
		},
	}
	DelayConsumer.Start()
	for i := 0; i < 10; i++ {
		DelayMapMap.AddItem(&QueueItem{
			ID: "test",
		}, time.Now().Add(time.Duration(i+2)*time.Second).UnixMilli())
	}
	time.Sleep(time.Second * 100)
}
