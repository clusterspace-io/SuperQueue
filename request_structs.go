package main

import "time"

type PostRecordRequest struct {
	Payload string  `json:"payload" validate:"required"`
	DelayMS float64 `json:"delay_ms"`
}

type PostRecordResponse struct {
	ID        string     `json:"id"`
	EnqueueAt *time.Time `json:"enqueue_at"`
}

type NackRecordRequest struct {
	DelayMS *float64 `json:"delay_ms"`
}
