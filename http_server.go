package main

import (
	"SuperQueue/logger"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	_ "net/http/pprof"

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
	Server.Echo.Start(":" + GetEnvOrDefault("HTTP_PORT", "8080"))
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

	if os.Getenv("TEST_MODE") == "true" {
		logger.Warn("TEST_MODE true, enabling debug routes")
		d := Server.Echo.Group("/debug")
		d.GET("/vars", wrapStdHandler)
		d.GET("/pprof/heap", wrapStdHandler)
		d.GET("/pprof/goroutine", wrapStdHandler)
		d.GET("/pprof/block", wrapStdHandler)
		d.GET("/pprof/threadcreate", wrapStdHandler)
		d.GET("/pprof/cmdline", wrapStdHandler)
		d.GET("/pprof/profile", wrapStdHandler)
		d.GET("/pprof/symbol", wrapStdHandler)
		d.GET("/pprof/trace", wrapStdHandler)
	}
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

// Wrapper for all stdlib /debug/* handlers
func wrapStdHandler(c echo.Context) error {
	w, r := c.Response().Writer, c.Request()
	if h, p := http.DefaultServeMux.Handler(r); len(p) != 0 {
		h.ServeHTTP(w, r)
		return nil
	}
	return echo.NewHTTPError(http.StatusNotFound)
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
	return c.String(200, GetMetrics())
}
