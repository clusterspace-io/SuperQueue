package main

import (
	"SuperQueue/logger"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/etcd-io/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	SQ *SuperQueue
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	c2 := make(chan struct{})

	logger.Logger.Logger.SetLevel(logrus.DebugLevel)
	logger.Info("Starting SuperQueue")

	if PARTITION == "" {
		logger.Error("Failed to provide a partition using the PARTITION env var, exiting")
		os.Exit(1)
	}
	var err error
	QueueMaxLen, err = strconv.ParseInt(GetEnvOrFail("QUEUE_LEN"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to calculate int64 for QUEUE_LEN of %s", GetEnvOrFail("QUEUE_LEN")))
	}

	SQ = NewSuperQueue("test-ns", PARTITION, 5, QueueMaxLen)
	go func() {
		StartHTTPServer()
	}()
	SQ.DelayConsumer.Start()

	// If ETCD_HOSTS exists, start reporting for service discovery
	if ETCD_HOSTS != "" {
		logger.Debug("Starting etcd based service discovery")
		hosts := strings.Split(ETCD_HOSTS, ",")
		logger.Debug("Using hosts: ", hosts)
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   hosts,
			DialTimeout: 2 * time.Second,
			DialOptions: []grpc.DialOption{grpc.WithBlock()}, // Need this to actually fail on connect
		})
		if err != nil {
			logger.Error("Failed to connect to etcd!")
			logger.Error(err)
			panic(err)
		} else {
			logger.Debug("Connected to etcd")
		}
		go func() {
			<-c2
			logger.Info("Closing service discovery ticker")
		}()
		defer cli.Close()
	}

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
	close(c2)
	logger.Info("Closing server")
	Server.Echo.Close()
}
