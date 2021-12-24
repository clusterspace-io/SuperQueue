package main

import (
	"SuperQueue/logger"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	PARTITION     = os.Getenv("PARTITION")
	ETCD_HOSTS    = os.Getenv("ETCD_HOSTS")
	ADVERTISE_URL = os.Getenv("ADVERTISE_URL")
	SCYLLA_HOSTS  = os.Getenv("SCYLLA_HOSTS")
	QUEUE_LEN     = os.Getenv("QUEUE_LEN")
	HTTP_PORT     = os.Getenv("HTTP_PORT")
)

func GetEnvOrDefault(env, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	} else {
		return e
	}
}

func GetEnvOrFail(env string) string {
	e := os.Getenv(env)
	if e == "" {
		panic(fmt.Sprintf("Failed to find env var '%s'", env))
	} else {
		return e
	}
}

func CheckFlags() {
	// These vars can use themselves as defaults because now they have
	flag.StringVar(&PARTITION, "partition", PARTITION, "Specify the partition")

	flag.StringVar(&HTTP_PORT, "port", HTTP_PORT, "Specify the http port to listen on")

	flag.StringVar(&ADVERTISE_URL, "advertise-url", ADVERTISE_URL, "Specifies the url to advertise to service discovery, should include http(s)://")

	flag.StringVar(&QUEUE_LEN, "queue-len", QUEUE_LEN, "Specifies the max queue length")

	flag.StringVar(&SCYLLA_HOSTS, "scylla-hosts", SCYLLA_HOSTS, "Specifies the scylla hosts")

	flag.StringVar(&ETCD_HOSTS, "etcd-hosts", ETCD_HOSTS, "Specifies the etcd hosts")

	flag.Parse()

	if PARTITION == "" {
		logger.Error("Failed to provide a partition using the PARTITION env var or -partition cli flag, exiting")
		os.Exit(1)
	}

	if HTTP_PORT == "" {
		logger.Error("Failed to provide a port using the HTTP_PORT env var or -port cli flag, exiting")
		os.Exit(1)
	}

	if QUEUE_LEN == "" {
		logger.Error("Failed to provide a queue length using the QUEUE_LEN env var or -queue-len cli flag, exiting")
		os.Exit(1)
	}
	var err error
	QueueMaxLen, err = strconv.ParseInt(QUEUE_LEN, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Failed to calculate int64 for QUEUE_LEN of %s", QUEUE_LEN))
	}
}
