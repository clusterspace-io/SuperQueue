package main

import (
	"SuperQueue/logger"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	SQ *SuperQueue
)

func main() {
	logger.Logger.Logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting SuperQueue")
	partition := os.Getenv("PARTITION")
	if partition == "" {
		logger.Error("Failed to provide a partition using the PARTITION env var, exiting")
		os.Exit(1)
	}
	var err error
	QueueMaxLen, err = strconv.ParseInt(GetEnvOrFail("QUEUE_LEN"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to calculate int64 for QUEUE_LEN of %s", GetEnvOrFail("QUEUE_LEN")))
	}

	SQ = NewSuperQueue("test-ns", partition, 5, QueueMaxLen)
	go func() {
		StartHTTPServer()
	}()
	SQ.DelayConsumer.Start()

	logger.Info("Setting up DB")
	// err := ConnectToDB(os.Getenv("CONN_STRING"))
	DBConnectWithoutKeyspace()
	DBKeyspaceSetup()
	DBConnect()
	DBTableSetup()
	// if err != nil {
	// 	panic(err)
	// }
	// err = CreateTables()
	// if err != nil {
	// 	panic(err)
	// }
	logger.Info("Done setting up db")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info("Closing server")
	Server.Echo.Close()
}
