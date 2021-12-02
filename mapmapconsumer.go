package main

import (
	"time"
)

type MapMapConsumer struct {
	ticker       time.Ticker
	endChan      chan struct{}
	MapMap       *MapMap
	ConsumerFunc func(int64, map[ItemID]*QueueItem)
	lastConsume  int64
}

// Gracefully shutdowns the consumer
func (m *MapMapConsumer) Stop() {
	m.endChan <- struct{}{}
}

// Starts consuming messages
func (m *MapMapConsumer) Start() {
	m.ticker = *time.NewTicker(time.Duration(m.MapMap.Interval) * time.Millisecond)
	m.lastConsume = time.Now().UnixMilli()
	go func() {
		for {
			select {
			case n := <-m.ticker.C:
				nowTime := n.UnixMilli()
				// fmt.Println("ticking", m.lastConsume, nowTime, m.MapMap.CalculateBucket(m.lastConsume), m.MapMap.CalculateBucket(nowTime))
				m.MapMap.ConsumeRange(m.lastConsume, nowTime, m.ConsumerFunc)
				m.lastConsume = nowTime + m.MapMap.Interval
			case <-m.endChan:
				// Exit
				return
			}
		}
	}()
}
