package main

import (
	"SuperQueue/logger"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/ksuid"
)

type HTTPServer struct {
	Echo *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

var (
	Server *HTTPServer
)

func StartHTTPServer() {
	echoInstance := echo.New()
	Server = &HTTPServer{
		Echo: echoInstance,
	}
	Server.Echo.HideBanner = true
	Server.Echo.Use(middleware.Logger())
	Server.Echo.Validator = &CustomValidator{validator: validator.New()}

	// Count requests
	Server.Echo.Use(IncrementCounter)
	Server.registerRoutes()

	logger.Info("Starting SuperQueue on port ", GetEnvOrDefault("HTTP_PORT", "8080"))
	Server.Echo.Logger.Fatal(Server.Echo.Start(":" + GetEnvOrDefault("HTTP_PORT", "8080")))
}

func IncrementCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		atomic.AddInt64(&TotalRequests, 1)
		return next(c)
	}
}

func PostRecordLatencyCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		atomic.AddInt64(&PostRecordLatency, int64(time.Since(start)))
		return nil
	}
}

func GetRecordLatencyCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		atomic.AddInt64(&GetRecordLatency, int64(time.Since(start)))
		return nil
	}
}

func AckRecordLatencyCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		atomic.AddInt64(&AckLatency, int64(time.Since(start)))
		return nil
	}
}

func NackRecordLatencyCounter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		atomic.AddInt64(&NackLatency, int64(time.Since(start)))
		return nil
	}
}

func (s *HTTPServer) registerRoutes() {
	s.Echo.GET("/hc", func(c echo.Context) error {
		return c.String(200, "y")
	})

	s.Echo.POST("/record", Post_Record, PostRecordLatencyCounter)
	s.Echo.GET("/record", Get_Record, GetRecordLatencyCounter)

	s.Echo.POST("/ack/:recordID", Post_AckRecord, AckRecordLatencyCounter)
	s.Echo.POST("/nack/:recordID", Post_NackRecord, NackRecordLatencyCounter)

	s.Echo.GET("/metrics", Get_Metrics)
}

func ValidateRequest(c echo.Context, s interface{}) error {
	if err := c.Bind(s); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(s); err != nil {
		return err
	}
	return nil
}

func Post_Record(c echo.Context) error {
	if atomic.LoadInt64(&QueuedMessages)+atomic.LoadInt64(&DelayedMessages)+atomic.LoadInt64(&InFlightMessages)+1 > QueueMaxLen {
		// We could exceed the max length if we do this
		return c.String(409, "Could exceed queue max length")
	}
	defer atomic.AddInt64(&PostRecordRequests, 1)
	bodyBytes, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		logger.Error("Failed to read body bytes:")
		logger.Error(err)
	}
	c.Request().Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	body := new(PostRecordRequest)
	if err := ValidateRequest(c, body); err != nil {
		logger.Debug("Validation failed ", err)
		atomic.AddInt64(&HTTP400s, 1)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	itemID := ksuid.New().String()

	var delayTime *time.Time

	if body.DelayMS > 1000 {
		// nil trick
		dt := time.Now().Add(time.Millisecond * time.Duration(body.DelayMS))
		delayTime = &dt
	}

	err = SQ.Enqueue(&QueueItem{
		ID:                     itemID,
		Payload:                []byte(body.Payload),
		CreatedAt:              time.Now(),
		StorageBucket:          "fake-bucket",
		ExpireAt:               time.Now().Add(4 * time.Hour),
		InFlightTimeoutSeconds: 30,
		Version:                0,
	}, delayTime)
	if err != nil {
		logger.Error("Failed to enqueue!")
		logger.Error(err)
		return c.String(500, err.Error())
	}

	atomic.AddInt64(&QueueMessageSize, int64(len(bodyBytes)))
	return c.String(http.StatusCreated, "")
}

func Get_Record(c echo.Context) error {
	defer atomic.AddInt64(&GetRecordRequests, 1)
	item, err := SQ.Dequeue()
	if err != nil {
		atomic.AddInt64(&HTTP500s, 1)
		return c.String(500, "Failed to dequeue record")
	}
	// Empty
	if item == nil {
		atomic.AddInt64(&EmptyQueueResponses, 1)
		return c.String(http.StatusNoContent, "Empty")
	}
	return c.JSON(200, map[string]interface{}{
		"id":       SQ.Partition + "_" + item.ID,
		"payload":  string(item.Payload),
		"attempts": item.Attempts,
	})
}

func Post_AckRecord(c echo.Context) error {
	recordID := c.Param("recordID")
	if recordID == "" {
		atomic.AddInt64(&HTTP400s, 1)
		return c.String(400, "No record ID given")
	}

	itemID := strings.Split(recordID, "_")[1]
	if itemID == "" {
		return c.String(http.StatusBadRequest, "Bad record ID given")
	}

	SQ.InFlightMapLock.Lock()
	item, exists := (*SQ.InFlightItems)[itemID]
	SQ.InFlightMapLock.Unlock()
	// Check if record exists
	if !exists {
		atomic.AddInt64(&AckMisses, 1)
		return c.String(404, "Record not found")
	}

	// Ack the record
	err := item.AckItem(SQ)
	if err != nil {
		atomic.AddInt64(&HTTP500s, 1)
		return c.String(500, "Failed to ack record")
	}
	return c.String(200, "")
}

func Post_NackRecord(c echo.Context) error {
	recordID := c.Param("recordID")
	if recordID == "" {
		atomic.AddInt64(&HTTP400s, 1)
		return c.String(400, "No record ID given")
	}

	itemID := strings.Split(recordID, "_")[1]
	if itemID == "" {
		return c.String(http.StatusBadRequest, "Bad record ID given")
	}

	item, exists := (*SQ.InFlightItems)[itemID]
	// Check if record exists
	if !exists {
		atomic.AddInt64(&NackMisses, 1)
		return c.String(404, "Record not found")
	}

	body := new(NackRecordRequest)
	if err := ValidateRequest(c, body); err != nil {
		logger.Debug("Validation failed ", err)
		atomic.AddInt64(&HTTP400s, 1)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	var delayMS *int64
	if body.DelayMS != nil {
		tmp := int64(*body.DelayMS)
		delayMS = &tmp
	}

	// Ack the record
	err := item.NackItem(SQ, delayMS)
	if err != nil {
		logger.Error(err)
		atomic.AddInt64(&HTTP500s, 1)
		return c.String(500, "Failed to ack record")
	}
	return c.String(200, "")
}

func Get_Metrics(c echo.Context) error {
	finalString := ""
	finalString += fmt.Sprintf("#TYPE in_flight_messages gauge\n#HELP in_flight_messages The current number of in-flight messages\nin_flight_messages %d", atomic.LoadInt64(&InFlightMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE total_in_flight_messages counter\n#HELP total_in_flight_messages The total number of in-flight messages\ntotal_in_flight_messages %d", atomic.LoadInt64(&TotalInFlightMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE queued_messages gauge\n#HELP queued_messages The total number of queued messages\nqueued_messages %d", atomic.LoadInt64(&QueuedMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE queued_messages_size gauge\n#HELP queued_messages_size The total number of bytes of the queued messages\nqueued_messages_size %d", atomic.LoadInt64(&QueueMessageSize)) + "\n"
	finalString += fmt.Sprintf("#TYPE queue_max_len gauge\n#HELP queue_max_len The max number of queued messages allowed\nqueue_max_len %d", atomic.LoadInt64(&QueueMaxLen)) + "\n"
	finalString += fmt.Sprintf("#TYPE total_queued_messages counter\n#HELP total_queued_messages The total number of queued messages\ntotal_queued_messages %d", atomic.LoadInt64(&TotalQueuedMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE delayed_messages gauge\n#HELP delayed_messages The current number of delayed messages\ndelayed_messages %d", atomic.LoadInt64(&DelayedMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE timedout_messages counter\n#HELP timedout_messages The total number of timedout messages\ntimedout_messages %d", atomic.LoadInt64(&TimedoutMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE acked_messages counter\n#HELP acked_messages The total number of acknowledged messages\nacked_messages %d", atomic.LoadInt64(&AckedMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE nacked_messages counter\n#HELP nacked_messages The total number of negatively messages\nnacked_messages %d", atomic.LoadInt64(&NackedMessages)) + "\n"
	finalString += fmt.Sprintf("#TYPE post_record_reqs counter\n#HELP post_record_reqs The total number of POST /record requests\npost_record_reqs %d", atomic.LoadInt64(&PostRecordRequests)) + "\n"
	finalString += fmt.Sprintf("#TYPE post_record_latency counter\n#HELP post_record_latency The sum of POST /record latency\npost_record_latency %d", atomic.LoadInt64(&PostRecordLatency)) + "\n"
	finalString += fmt.Sprintf("#TYPE get_record_reqs counter\n#HELP get_record_reqs The total number of GET /record requests\nget_record_reqs %d", atomic.LoadInt64(&GetRecordRequests)) + "\n"
	finalString += fmt.Sprintf("#TYPE get_record_latency counter\n#HELP get_record_latency The sum of GET /record latencies\nget_record_latency %d", atomic.LoadInt64(&GetRecordLatency)) + "\n"
	finalString += fmt.Sprintf("#TYPE ack_misses counter\n#HELP ack_misses The total number of ack requests that fail to ack a message\nack_misses %d", atomic.LoadInt64(&AckMisses)) + "\n"
	finalString += fmt.Sprintf("#TYPE nack_misses counter\n#HELP nack_misses The total number of nack requests that fail to nack a message\nnack_misses %d", atomic.LoadInt64(&NackMisses)) + "\n"
	finalString += fmt.Sprintf("#TYPE empty_queue_responses counter\n#HELP empty_queue_responses The total number of GET /record requests that result in an empty queue response\nempty_queue_responses %d", atomic.LoadInt64(&EmptyQueueResponses)) + "\n"
	finalString += fmt.Sprintf("#TYPE full_queue_responses counter\n#HELP full_queue_responses The total number of GET /record requests that result in a full queue response\nfull_queue_responses %d", atomic.LoadInt64(&FullQueueResponses)) + "\n"
	finalString += fmt.Sprintf("#TYPE http_total_requests counter\n#HELP http_total_requests The total number of http requests processed returning any code\nhttp_total_requests %d", atomic.LoadInt64(&TotalRequests)) + "\n"
	finalString += fmt.Sprintf("#TYPE http_500s counter\n#HELP http_500s The total number of returned 500 http responses\nhttp_500s %d", atomic.LoadInt64(&HTTP500s)) + "\n"
	finalString += fmt.Sprintf("#TYPE http_400s counter\n#HELP http_400s The total number of returned 400 http responses\nhttp_400s %d", atomic.LoadInt64(&HTTP400s)) + "\n"
	finalString += fmt.Sprintf("#TYPE ack_latency counter\n#HELP ack_latency The sum of POST /ack latency. Use both ack_misses and acked_messages to calculate latency per request\nack_latency %d", atomic.LoadInt64(&AckLatency)) + "\n"
	finalString += fmt.Sprintf("#TYPE nack_latency counter\n#HELP nack_latency The sum of POST /nack latency. Use both nack_misses and nacked_messages to calculate latency per request\nnack_latency %d", atomic.LoadInt64(&NackLatency)) + "\n"

	// Resource metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	finalString += fmt.Sprintf("#TYPE mem_bytes_heap gauge\n#HELP mem_bytes_heap The current memory usage in bytes of the process' heap\nmem_bytes_heap %d", m.Alloc) + "\n"
	finalString += fmt.Sprintf("#TYPE mem_bytes_sys gauge\n#HELP mem_bytes_sys The current memory usage in bytes that the process has obtained from the os\nmem_bytes_sys %d", m.Sys) + "\n"

	return c.String(200, finalString)
}
