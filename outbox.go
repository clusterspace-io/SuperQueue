package main

type Outbox struct {
	SQ          *SuperQueue
	MessageChan chan *QueueItem
}

func NewOutbox(sq *SuperQueue, bufferLength int64) *Outbox {
	ob := &Outbox{
		SQ:          sq,
		MessageChan: make(chan *QueueItem, bufferLength),
	}
	return ob
}

// Adds an item to the outbox in memory.
func (o *Outbox) Add(item *QueueItem) {
	o.MessageChan <- item
}

// Pops an item from the outbox in memory.
func (o *Outbox) Pop() *QueueItem {
	select {
	case item := <-o.MessageChan:
		return item
	default:
		return nil
	}
}
