package main

import (
	"SuperQueue/logger"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	SQ *SuperQueue
)

func main() {
	logger.Logger.Logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting SuperQueue")
	SQ = NewSuperQueue(5, 2<<20)
	go func() {
		StartHTTPServer()
	}()
	SQ.DelayConsumer.Start()

	logger.Info("Setting up DB")
	err := ConnectToDB("postgresql://dan:thisisabadpassword@free-tier.gcp-us-central1.cockroachlabs.cloud:26257/defaultdb?sslmode=require&options=--cluster%3Dportly-impala-2852")
	if err != nil {
		panic(err)
	}
	err = CreateTables()
	if err != nil {
		panic(err)
	}
	logger.Info("Done setting up db")

	for i := 0; i < 10; i++ {
		itemID := fmt.Sprintf("test-%d", time.Now().UnixNano())
		p := []byte("hey")
		SQ.Enqueue(&QueueItem{
			ID:                     itemID,
			Payload:                p,
			CreatedAt:              time.Now(),
			Bucket:                 "fake-bucket",
			ExpireAt:               time.Now().Add(4 * time.Hour),
			InFlightTimeoutSeconds: 30,
			BackoffMinMS:           300,
			BackoffMultiplier:      2,
			Version:                0,
		}, int64((i+1)*1000))
	}
	time.Sleep(time.Second * 100)
	Server.Echo.Close()
}
