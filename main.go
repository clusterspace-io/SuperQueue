package main

import (
	"SuperQueue/logger"
	"time"
)

var (
	SQ *SuperQueue
)

func main() {
	logger.Info("Starting SuperQueue")
	SQ = NewSuperQueue(5, 2<<20)
	SQ.DelayConsumer.Start()
	for i := 0; i < 10; i++ {
		SQ.Enqueue(&QueueItem{
			ID: "test",
		}, int64(i*1000))
	}
	time.Sleep(time.Second * 100)
}
