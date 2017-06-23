package main

/*
 * tests run server themselves
 * benchmarks require running server. Call like this:
 * go test -run Dial -bench . -cpu=1,4,8
 */

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/rpc"
	"os"
	"tc/bootstrap"
	"tc/stat"
	"testing"
	"time"
)

var client *rpc.Client
var server *Server

func init() {
	bootstrap.LoadTimeZone()
}

func TestMain(t *testing.T) {
	server = newServer()
	go server.listenAndServe()
}

func TestDial(t *testing.T) {
	client = dial(t.Fatal)
}

func dial(fatal func(args ...interface{})) *rpc.Client {
	client, err := rpc.Dial(config.net, config.laddr)
	if err != nil {
		fatal("dialing:", err)
	}
	return client
}

func TestMerging(t *testing.T) {
	const clicks = 13
	m := stat.NewMap()
	e := randomEvent()
	e.ClickCount = clicks

	m.Add(e)
	err := callMerge(m)
	if err != nil {
		t.Fatal("Got error: ", err)
	}
	server.merger.wg.Wait()
	found := server.merger.data.Get(e.Key())
	if found == nil {
		t.Fatalf("Event not found after merge")
	}
	if e.ClickCount != clicks {
		t.Fatalf("1) Bad ClickCount: expected %d, got %d", clicks, e.ClickCount)
	}

	m.Add(e)
	err = callMerge(m)
	if err != nil {
		t.Fatal("Got error: ", err)
	}
	server.merger.wg.Wait()
	if e.ClickCount != 2*clicks {
		t.Fatalf("2) Bad ClickCount: expected %d, got %d", 2*clicks, e.ClickCount)
	}
}

func TestDumping(t *testing.T) {
	server.dump()
	if server.merger.data.Len() != 0 {
		t.Fatal("Merger is not empty after dump")
	}
}

func TestAggregation(t *testing.T) {
	m := stat.NewMap()
	for i := 0; i < 100; i++ {
		raw := stat.NewRawEvent(stat.UniqueShowEvent)
		raw.Campaign.Id = i % 7
		raw.Campaign.Ad.Id = i % 19
		raw.Site.Id = i % 13
		m.AddRaw(raw)
	}
	merger := newMerger()
	merger.push(m)
	merger.finish()

	a := newAggregator()
	tables := a.aggregate(merger.data.Slice())
	for _, tbl := range tables {
		n := 0
		switch tbl.Name() {
		case "statsByAds":
			n = 19
		case "statsBySites":
			n = 13
		case "statsByCampaigns":
			n = 7
		default:
			continue
		}
		rows := len(tbl.Rows())
		if rows != n {
			t.Fatalf("Exptected %d rows in %s, got %d\n", n, tbl.Name(), rows)
		}
	}
}

func TestMergingFile(t *testing.T) {
	server := newServer()
	go server.dumper()
	file, err := os.Open("../stat/raw-test.data")
	check(err)
	dec := json.NewDecoder(file)
	i := 0
	m := stat.NewMap()
	for {
		var raw stat.RawEvent
		err := dec.Decode(&raw)
		if err == io.EOF {
			break
		}
		check(err)
		raw.Date = raw.Date.Local().Truncate(time.Hour)
		m.AddRaw(&raw)
		i++
	}
	{
		var shows int64
		for _, e := range m.Slice() {
			if e.Site.Id == 1771 {
				shows += e.ShowCount
			}
		}
		const S = 12346
		if shows != S {
			t.Fatalf("Expected %d ShowCount for site %d, got %d\n", S, 1771, shows)
		}
	}

	const M = 61512
	if i != M {
		t.Fatalf("Expected %d raw records, got %d\n", M, i)
	}
	server.Merge(*m, nil)
	server.aggregator.dump(server.merger)
	{
		var shows int64
		for _, e := range server.merger.data.Slice() {
			if e.Site.Id == 1771 {
				shows += e.ShowCount
			}
		}
		const S = 12346
		if shows != S {
			t.Fatalf("Expected %d ShowCount for site %d, got %d\n", S, 1771, shows)
		}
	}

	for len(server.aggregator.queue) > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(1000 * time.Millisecond)
	print("check postgres for valid data\n")
}

func BenchmarkMerging1000(b *testing.B) {
	benchmarkMerging(b, 1000)
}
func BenchmarkMerging10000(b *testing.B) {
	benchmarkMerging(b, 10000)
}
func BenchmarkMerging100000(b *testing.B) {
	benchmarkMerging(b, 100000)
}

func benchmarkMerging(b *testing.B, n int) {
	m := randomMap(n)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := callMerge(m)
			if err != nil {
				b.Fatal("Got error:", err)
			}
		}
	})
}

func callMerge(m *stat.Map) error {
	return client.Call("Aggregator.Merge", m, nil)
}

func randomEvent() *stat.Event {
	r := rand.Int() + 1
	t := stat.EventType(1 + r%stat.EventsNum)
	e := stat.NewEvent(stat.NewRawEvent(t))
	e.Campaign.Id = r % 1000
	e.Campaign.UserId = r % 100
	e.Campaign.Ad.Id = r % 2500
	e.Site.Id = r % 1500
	e.Geo.Id = r % 6000
	e.Geo.CountryId = r % 10
	return e
}

func randomMap(n int) *stat.Map {
	m := stat.NewMap()
	for i := 0; i < n; i++ {
		m.Add(randomEvent())
	}
	return m
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
