package main

import (
	"SuperQueue/logger"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	SQ *SuperQueue
)

func main() {
	logger.Logger.SetLevel(logrus.DebugLevel)
	if os.Getenv("TEST_MODE") == "true" {
		logger.Warn("TEST_MODE true, enabling cpu profiling")
		f, perr := os.Create("cpu.pprof")
		if perr != nil {
			panic(perr)
		}
		perr = pprof.StartCPUProfile(f)
		if perr != nil {
			panic(perr)
		}
		defer f.Close()
		defer pprof.StopCPUProfile()
	}
	CheckFlags()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	logger.Info("Starting SuperQueue")

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
