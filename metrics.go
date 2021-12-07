package main

var (
	InFlightMessages      int64 = 0
	TotalInFlightMessages int64 = 0
	QueuedMessages        int64 = 0
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
)
