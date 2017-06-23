package main

import (
	"log"
	"sync"
	"tc/stat"
)

type (
	aggregator struct {
		tables []Table
		queue  chan []Table
	}
)

func newAggregator() *aggregator {
	a := &aggregator{
		tables: makeTables(),
		queue:  make(chan []Table, 100),
	}
	go a.saver()
	return a
}

func (a *aggregator) dump(m *merger) {
	m.finish()
	full := m.data.Slice()
	stats := full
	log.Printf("Dumping %d events\b", len(stats))
	log.Printf("Saver queue size: %d\n", len(a.queue))

	for len(stats) > 0 {
		n := min(len(stats), config.dumpBatchSize)
		a.queue <- a.aggregate(stats[:n])
		a.tables = makeTables()
		stats = stats[n:]
	}
	saveFull(full)
}

func (a *aggregator) saver() {
	for tables := range a.queue {
		saveTables(tables)
	}
}

func (a *aggregator) aggregate(stats stat.Slice) []Table {
	var wg sync.WaitGroup
	wg.Add(len(a.tables))
	for _, t := range a.tables {
		go func(t Table) {
			t.Add(stats)
			wg.Done()
		}(t)
	}
	wg.Wait()
	return a.tables
}
