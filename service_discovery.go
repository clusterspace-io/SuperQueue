package main

import (
	"SuperQueue/logger"
	"strings"
	"time"

	"github.com/etcd-io/etcd/clientv3"
	"google.golang.org/grpc"
)

var (
	EtcdClient *clientv3.Client
	SDTicker   *time.Ticker
)

// Tries to start etcd based service discovery. Returns whether SD was configured and setup.
func TryEtcdSD(sq *SuperQueue) bool {
	// If ETCD_HOSTS exists, start reporting for service discovery
	if ETCD_HOSTS != "" {
		logger.Debug("Starting etcd based service discovery")
		hosts := strings.Split(ETCD_HOSTS, ",")
		logger.Debug("Using hosts: ", hosts)
		var err error
		EtcdClient, err = clientv3.New(clientv3.Config{
			Endpoints:   hosts,
			DialTimeout: 2 * time.Second,
			DialOptions: []grpc.DialOption{grpc.WithBlock()}, // Need this to actually block and fail on connect error
		})
		if err != nil {
			logger.Error("Failed to connect to etcd!")
			logger.Error(err)
			panic(err)
		} else {
			logger.Debug("Connected to etcd")
		}
		UpdateSD()
		SDTicker = time.NewTicker(10 * time.Second)
		go func() {
			for {
				select {
				case <-SDTicker.C:
					logger.Debug("SD Tick")
					UpdateSD()

				case <-sq.CloseChan:
					logger.Info("Closing service discovery ticker")
					SDTicker.Stop()
					return
				}
			}
		}()
		return true
	}
	return false
}

// Updates the service discovery entry for this partition
func UpdateSD() error {
	return nil
}
