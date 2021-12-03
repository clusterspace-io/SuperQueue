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
		SQ.Enqueue(&QueueItem{
			ID: "test",
		}, int64(i*1000))
		fmt.Println("Currnet time", time.Now().UnixMilli())
	}
	time.Sleep(time.Second * 100)
}
