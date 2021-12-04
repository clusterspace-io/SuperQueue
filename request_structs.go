package main

type PostRecordRequest struct {
	Payload string  `json:"payload" validate:"required"`
	DelayMS float64 `json:"delay_ms"`
}

type PostRecordResponse struct {
	ID string `json:"id"`
}
