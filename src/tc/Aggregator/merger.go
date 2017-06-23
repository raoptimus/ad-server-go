package main

import (
	"sync"
	"tc/stat"
)

type merger struct {
	data  *stat.Map
	wg    sync.WaitGroup
	queue chan *stat.Map
	done  chan bool
}

func newMerger() *merger {
	m := &merger{
		data:  stat.NewMap(),
		queue: make(chan *stat.Map, config.queueSize),
		done:  make(chan bool),
	}
	go m.merge()
	return m
}

func (m *merger) push(sm *stat.Map) {
	m.wg.Add(1)
	m.queue <- sm
}

func (m *merger) merge() {
	for data := range m.queue {
		m.data.Merge(data)
		m.wg.Done()
	}
	m.done <- true
}

//wait for queued merges and finish merging goroutine
func (m *merger) finish() {
	m.wg.Wait()
	close(m.queue)
	<-m.done
}
