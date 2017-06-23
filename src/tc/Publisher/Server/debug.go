package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

var measureDefaultRate = 100000

type doCallback func()

//usage: defer un(func() unCalback)
func do(c doCallback) {
	c()
}

//usage: defer un(trace("custom-label"))
func trace(s string) doCallback {
	start := time.Now()
	log.Println("trace-start:", s)
	return func() {
		log.Printf("trace-end: %s; Elapsed time: %v\n\n", s, time.Now().Sub(start))
	}
}

func long(s string, min time.Duration) doCallback {
	start := time.Now()
	return func() {
		d := time.Now().Sub(start)
		if d > min {
			log.Printf("Long call: %s: %v > %v\n", s, d, min)
		}
	}
}

type measurer struct {
	total      time.Duration
	min        time.Duration
	max        time.Duration
	n          int
	printEvery int
}

var measures map[string]*measurer

func measureReset(s string) {
	log.Println("measureReset:", s)
	delete(measures, s)
}

func measureResetAll() {
	for s := range measures {
		measureReset(s)
	}
}

//usage: defer un(measure("custom-label"))
func measure(s string) doCallback {
	return measureEvery(s, measureDefaultRate)
}

//usage: defer un(measureEvery("custom-label", 100))
func measureEvery(s string, printEvery int) doCallback {
	if measures == nil {
		measures = make(map[string]*measurer)
	}

	m, ok := measures[s]

	if !ok {
		m = &measurer{n: 1, printEvery: printEvery}
		measures[s] = m
	}

	start := time.Now()

	return func() {
		d := time.Now().Sub(start)

		switch {
		case d > m.max:
			m.max = d
		case d < m.min || m.min == 0:
			m.min = d
		}

		m.total += d
		m.n++

		if m.n%printEvery == 0 && m.max >= time.Duration(30*time.Millisecond) {
			log.Printf(
				"[%s] avg: %v; min: %v; max: %v; calls: %d\n",
				s,
				m.total/time.Duration(m.n),
				m.min,
				m.max,
				m.n,
			)

			//reset
			m = &measurer{n: 1, printEvery: printEvery}
			measures[s] = m
		}
	}
}

// resets measure counters
// without query args resets all
// otherwise resets listed
// usage: http://localhost:9090/measure/reset?filter-ads&session-start
func measureResetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(q) == 0 {
		measureResetAll()
	}
	for s := range q {
		measureReset(s)
	}
}

// usage: http://localhost:9090/measure/set?rate=10000
func measureSetHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if rate, err := strconv.Atoi(q.Get("rate")); err == nil {
		measureDefaultRate = rate
	} else {
		log.Println("measureSetHandler error:", err)
	}
}

// usage: http://localhost:9090/measure/print
//todo print to ResponseWriter
func measurePrintHandler(w http.ResponseWriter, r *http.Request) {
	for s, m := range measures {
		log.Printf(
			"[%s] avg: %v; min: %v; max: %v; calls: %d\n",
			s,
			m.total/time.Duration(m.n),
			m.min,
			m.max,
			m.n,
		)
	}

	return
}

func LogPrintln(v ...interface{}) {
	//TODO if not debug mode then return;
	log.Println(v)
}
