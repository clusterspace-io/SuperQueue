package main

import (
	"SuperQueue/logger"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	InFlightMessages       int64 = 0
	InFlightMessagesMetric       = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_messages",
		Help: "The current number of in-flight messages",
	})

	QueuedMessages       int64 = 0
	QueuedMessagesMetric       = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "queued_messages",
		Help: "The current number of queued messages",
	})

	QueueMessageSizeMetric = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "queue_message_size",
		Help:    "The size in bytes of the messages queued",
		Buckets: prometheus.ExponentialBuckets(10, 2, 8),
	})

	QueueMaxLen       int64 = 0
	QueueMaxLenMetric       = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "queue_max_len",
		Help: "The max number of queued messages allowed in the partition",
	})

	DelayedMessages       int64 = 0
	DelayedMessagesMetric       = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "delayed_messages",
		Help: "The current number of delayed messages",
	})

	TimedoutMessages       int64 = 0
	TimedoutMessagesMetric       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "timedout_messages",
		Help: "The number of timedout messages",
	})

	AckedMessages        int64 = 0
	AckedMessagesCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "acked_messages",
		Help: "The number of acked messages",
	})

	NackedMessages        int64 = 0
	NackedMessagesCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nacked_messages",
		Help: "The number of nacked messages",
	})

	PostRecordRequests       int64 = 0
	PostRecordRequestCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "post_record_reqs",
		Help: "The total number of POST /record requests",
	})

	PostRecordLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "post_record_latency",
		Help:    "The latency of POST /record requests in milliseconds",
		Buckets: prometheus.ExponentialBuckets(1, 3, 10),
	})

	GetRecordRequests       int64 = 0
	GetRecordRequestCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "get_record_reqs",
		Help: "The total number of GET /record requests",
	})

	GetRecordLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "get_record_latency",
		Help:    "The latency of GET /record requests in milliseconds",
		Buckets: prometheus.ExponentialBuckets(1, 3, 10),
	})

	AckMisses        int64 = 0
	AckMissesCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "ack_misses",
		Help: "The number of POST /ack requests that miss",
	})

	NackMisses        int64 = 0
	NackMissesCounter       = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "nack_misses",
		Help: "The number of POST /nack requests that miss",
	})

	EmptyQueueResponses int64 = 0
	FullQueueResponses  int64 = 0
	TotalRequests       int64 = 0
	HTTP500s            int64 = 0
	HTTP400s            int64 = 0
	AckLatency          int64 = 0
	NackLatency         int64 = 0

	Metrics = []string{}
)

func FormatMetric(metricName, metricType, description string, value interface{}) string {
	hn, err := os.Hostname()
	if err != nil {
		logger.Error("Failed to get hostname!")
		panic(err)
	}
	return fmt.Sprintf("#TYPE %s %s\n#HELP %s %s\n%s{host=%s, queue=%s, partition=%s} %v", metricName, metricType, metricName, description, metricName, hn, SQ.Name, SQ.Partition, value)
}

func SetupMetrics() {
	prometheus.Register(InFlightMessagesMetric)
	prometheus.Register(QueuedMessagesMetric)
	prometheus.Register(QueueMessageSizeMetric)
	prometheus.Register(QueueMaxLenMetric)
	prometheus.Register(DelayedMessagesMetric)
	prometheus.Register(TimedoutMessagesMetric)
	prometheus.Register(AckedMessagesCounter)
	prometheus.Register(NackedMessagesCounter)
	prometheus.Register(PostRecordRequestCounter)
	prometheus.Register(PostRecordLatency)
	prometheus.Register(GetRecordRequestCounter)
	prometheus.Register(GetRecordRequestCounter)
	prometheus.Register(AckMissesCounter)
	prometheus.Register(NackMissesCounter)
}

func GetMetrics() string {
	metrics := make([]string, 30)

	metrics = append(metrics, FormatMetric("empty_queue_responses", "counter", "The total number of GET /record requests that result in an empty queue response", atomic.LoadInt64(&EmptyQueueResponses)))

	metrics = append(metrics, FormatMetric("full_queue_responses", "counter", "The total number of GET /record requests that result in a full queue response", atomic.LoadInt64(&FullQueueResponses)))

	metrics = append(metrics, FormatMetric("http_total_requests", "counter", "The total number of http requests processed returning any code", atomic.LoadInt64(&TotalRequests)))

	metrics = append(metrics, FormatMetric("http_500s", "counter", "The total number of returned 500 http responses", atomic.LoadInt64(&HTTP500s)))

	metrics = append(metrics, FormatMetric("http_400s", "counter", "The total number of returned 400 http responses", atomic.LoadInt64(&HTTP400s)))

	metrics = append(metrics, FormatMetric("ack_latency", "counter", "The sum of POST /ack latency. Use both ack_misses and acked_messages to calculate latency per request", atomic.LoadInt64(&AckLatency)))

	metrics = append(metrics, FormatMetric("nack_latency", "counter", "The sum of POST /nack latency. Use both nack_misses and nacked_messages to calculate latency per request", atomic.LoadInt64(&NackLatency)))

	// Resource metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics = append(metrics, FormatMetric("mem_bytes_heap", "gauge", "The current memory usage in bytes of the process' heap", m.Alloc))
	metrics = append(metrics, FormatMetric("mem_bytes_sys", "gauge", "The current memory usage in bytes that the process has obtained from the os", m.Sys))

	return strings.Join(metrics, "\n")
}
