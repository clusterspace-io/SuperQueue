package main

import (
	"SuperQueue/logger"
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

	logger.Warn("Sleeping for 1000 seconds before shutdown")
	time.Sleep(time.Second * 1000)
	logger.Warn("Shutting down")
	Server.Echo.Close()
}
