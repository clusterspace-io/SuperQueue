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

// Adds an item to the outbox in memory. Returns whether it was full.
func (o *Outbox) Add(item *QueueItem) bool {
	// TODO: Handle when full
	select {
	case o.MessageChan <- item:
		return false
	default:
		return true
	}
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
