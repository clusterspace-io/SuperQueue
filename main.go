package main

import (
	"SuperQueue/logger"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	SQ *SuperQueue
)

func main() {
	if os.Getenv("TEST_MODE") == "true" {
		logger.Warn("TEST_MODE true, enabling cpu profiling")
		f, perr := os.Create("cpu.pprof")
		if perr != nil {
			panic(perr)
		}
		runtime.SetCPUProfileRate(100)
		perr = pprof.StartCPUProfile(f)
		if perr != nil {
			panic(perr)
		}
		defer f.Close()
		defer pprof.StopCPUProfile()
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	logger.Logger.Logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting SuperQueue")

	if PARTITION == "" {
		logger.Error("Failed to provide a partition using the PARTITION env var, exiting")
		os.Exit(1)
	}
	if ADVERTISE_ADDRESS == "" {
		logger.Error("Failed to provide a advertise address using the ADVERTISE_ADDRESS env var, exiting")
		os.Exit(1)
	}
	var err error
	QueueMaxLen, err = strconv.ParseInt(GetEnvOrFail("QUEUE_LEN"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to calculate int64 for QUEUE_LEN of %s", GetEnvOrFail("QUEUE_LEN")))
	}

	SQ = NewSuperQueue("test-ns", PARTITION, 5, QueueMaxLen)
	// Try to setup service discovery
	TryEtcdSD(SQ)
	go func() {
		StartHTTPServer()
	}()
	SQ.DelayConsumer.Start()

	logger.Debug("Setting up DB")
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
	logger.Debug("Done setting up db")

	<-c
	SQ.Close()
	logger.Info("Closing server")
	Server.Echo.Close()
}
