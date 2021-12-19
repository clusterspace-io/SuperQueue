package main

import (
	"SuperQueue/logger"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/etcd-io/etcd/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	EtcdClient *clientv3.Client
	SDTicker   *time.Ticker
)

type PartitionSDRecord struct {
	QueueName  string
	Partition  string
	UpdatedAt  time.Time
	Address    string
	IsDraining bool
}

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
		UpdateSD(context.Background())
		SDTicker = time.NewTicker(10 * time.Second)
		go func() {
			for {
				select {
				case <-SDTicker.C:
					logger.Debug("SD Tick")
					UpdateSD(context.Background())

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
func UpdateSD(c context.Context) error {
	logger.Debug("Updating service discovery...")
	t := time.Now()
	ctx, cancelFunc := context.WithTimeout(c, time.Second*1) // 1 second timeout
	defer cancelFunc()
	r := &PartitionSDRecord{
		QueueName:  SQ.Name,
		Partition:  SQ.Partition,
		UpdatedAt:  time.Now(),
		IsDraining: false,
		Address:    ADVERTISE_ADDRESS,
	}
	b, err := json.Marshal(r)
	if err != nil {
		logger.Error("Error marshalling service discovery record!")
		logger.Error(err)
		return err
	}
	EtcdClient.KV.Put(ctx, fmt.Sprintf("q_%s_%s", SQ.Name, SQ.Partition), string(b))
	logger.Debug("Updated service discovery in ", time.Since(t))
	return nil
}
