package main

import (
	"SuperQueue/logger"
	"time"
)

var (
	DelayMapMap   *MapMap
	InFlightItems *map[ItemID]QueueItem
	Outbox        chan *QueueItem
	DelayTicker   time.Ticker
)

func main() {
	logger.Info("Starting SuperQueue")
}
