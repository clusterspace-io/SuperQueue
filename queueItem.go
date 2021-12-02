package main

type ItemID string

type QueueItem struct {
	ID   ItemID
	Body []byte
}
