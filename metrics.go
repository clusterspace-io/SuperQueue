package main

import (
	"SuperQueue/logger"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync/atomic"
)

var (
	InFlightMessages      int64 = 0
	TotalInFlightMessages int64 = 0
	QueuedMessages        int64 = 0
	QueueMessageSize      int64 = 0
	QueueMaxLen           int64 = 0
	TotalQueuedMessages   int64 = 0
	DelayedMessages       int64 = 0
	TimedoutMessages      int64 = 0
	AckedMessages         int64 = 0
	NackedMessages        int64 = 0
	PostRecordRequests    int64 = 0
	PostRecordLatency     int64 = 0
	GetRecordRequests     int64 = 0
	GetRecordLatency      int64 = 0
	AckMisses             int64 = 0
	NackMisses            int64 = 0
	EmptyQueueResponses   int64 = 0
	FullQueueResponses    int64 = 0
	TotalRequests         int64 = 0
	HTTP500s              int64 = 0
	HTTP400s              int64 = 0
	AckLatency            int64 = 0
	NackLatency           int64 = 0

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

func GetMetrics() string {
	metrics := make([]string, 30)
	metrics = append(metrics, FormatMetric("in_flight_messages", "gauge", "The current number of in-flight messages", atomic.LoadInt64(&InFlightMessages)))
	metrics = append(metrics, FormatMetric("total_in_flight_messages", "counter", "The total number of in-flight messages ever sent", atomic.LoadInt64(&TotalInFlightMessages)))
	metrics = append(metrics, FormatMetric("queued_messages", "gauge", "The current number of queued messages", atomic.LoadInt64(&QueuedMessages)))
	metrics = append(metrics, FormatMetric("queued_messages_size", "counter", "The sum of all queued message body byte lengths", atomic.LoadInt64(&QueueMessageSize)))
	metrics = append(metrics, FormatMetric("total_queued_messages", "counter", "The sum of all queued messages since process start", atomic.LoadInt64(&TotalQueuedMessages)))
	metrics = append(metrics, FormatMetric("queue_max_len", "gauge", "The max number of queued messages allowed", atomic.LoadInt64(&QueueMaxLen)))
	metrics = append(metrics, FormatMetric("delayed_messages", "gauge", "The current number of delayed messages", atomic.LoadInt64(&DelayedMessages)))

	metrics = append(metrics, FormatMetric("timedout_messages", "counter", "The total number of timedout messages", atomic.LoadInt64(&TimedoutMessages)))

	metrics = append(metrics, FormatMetric("acked_messages", "counter", "The total number of acknowledged messages", atomic.LoadInt64(&AckedMessages)))

	metrics = append(metrics, FormatMetric("nacked_messages", "counter", "The total number of negatively messages", atomic.LoadInt64(&NackedMessages)))

	metrics = append(metrics, FormatMetric("post_record_reqs", "counter", "The total number of POST /record requests", atomic.LoadInt64(&PostRecordRequests)))

	metrics = append(metrics, FormatMetric("post_record_latency", "counter", "The sum of POST /record latency", atomic.LoadInt64(&PostRecordLatency)))

	metrics = append(metrics, FormatMetric("get_record_reqs", "counter", "The total number of GET /record requests", atomic.LoadInt64(&GetRecordRequests)))

	metrics = append(metrics, FormatMetric("get_record_latency", "counter", "The sum of GET /record latencies", atomic.LoadInt64(&GetRecordLatency)))

	metrics = append(metrics, FormatMetric("ack_misses", "counter", "The total number of ack requests that fail to ack a message", atomic.LoadInt64(&AckMisses)))

	metrics = append(metrics, FormatMetric("nack_misses", "counter", "The total number of nack requests that fail to nack a message", atomic.LoadInt64(&NackMisses)))

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
