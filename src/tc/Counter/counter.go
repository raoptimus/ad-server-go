package main

import (
	"sync"
	"sync/atomic"
	"tc/stat"
	"time"
)

type counter struct {
	counted      int64
	events       *stat.Map
	serverEvents *stat.Map
	rawQueue     chan stat.RawEvent
	wg           sync.WaitGroup
}

func newCounter() *counter {
	c := &counter{
		events:       stat.NewMap(),
		serverEvents: stat.NewMap(),
		rawQueue:     make(chan stat.RawEvent, config.rawQueueSize),
	}
	go c.count()
	return c
}

func (c *counter) push(raw stat.RawEvent) {
	atomic.AddInt64(&c.counted, 1)
	c.wg.Add(1)
	c.rawQueue <- raw
}

func (c *counter) count() {
	for raw := range c.rawQueue {
		raw.Date = raw.Date.Local().Truncate(5 * time.Minute)
		c.serverEvents.AddRaw(&raw)

		raw.Date = raw.Date.Local().Truncate(time.Hour)
		c.events.AddRaw(&raw)

		c.wg.Done()
	}
}

func (c *counter) stop() {
	c.wg.Wait()
	close(c.rawQueue)
}
