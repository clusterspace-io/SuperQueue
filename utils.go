package main

import (
	"fmt"
	"os"
)

var (
	PARTITION         = os.Getenv("PARTITION")
	ETCD_HOSTS        = os.Getenv("ETCD_HOSTS")
	ADVERTISE_ADDRESS = os.Getenv("ADVERTISE_ADDRESS")
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
