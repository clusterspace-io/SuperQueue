package main

import (
	"SuperQueue/logger"
)

var (
	DelayMapMap   *MapMap
	InFlightItems *map[ItemID]QueueItem
	Outbox        chan *QueueItem
)

func main() {
	logger.Info("Starting SuperQueue")
}
