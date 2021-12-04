package main

import (
	"SuperQueue/logger"
	"net/http"
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

	Server.registerRoutes()

	logger.Info("Starting Host API on port 8080")
	Server.Echo.Logger.Fatal(Server.Echo.Start(":8080"))
}

func (s *HTTPServer) registerRoutes() {
	s.Echo.POST("/record", Post_Record)
	s.Echo.GET("/record", Get_Record)
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
	body := new(PostRecordRequest)
	if err := ValidateRequest(c, body); err != nil {
		logger.Debug("Validation failed ", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	itemID := ksuid.New().String()

	var delayTime *time.Time

	if body.DelayMS > 1000 {
		// nil trick
		dt := time.Now().Add(time.Millisecond * time.Duration(body.DelayMS))
		delayTime = &dt
	}

	SQ.Enqueue(&QueueItem{
		ID:                     itemID,
		Payload:                []byte(body.Payload),
		CreatedAt:              time.Now(),
		Bucket:                 "fake-bucket",
		ExpireAt:               time.Now().Add(4 * time.Hour),
		InFlightTimeoutSeconds: 30,
		BackoffMinMS:           300,
		BackoffMultiplier:      2,
		Version:                0,
	}, delayTime)

	return c.JSON(200, PostRecordResponse{
		ID:        itemID,
		EnqueueAt: delayTime,
	})
}

func Get_Record(c echo.Context) error {
	item, err := SQ.Dequeue()
	if err != nil {
		return c.String(500, "Failed to dequeue item")
	}
	// Empty
	if item == nil {
		return c.String(http.StatusNoContent, "Empty")
	}
	return c.JSON(200, map[string]string{
		"id":      item.ID,
		"payload": string(item.Payload),
	})
}
