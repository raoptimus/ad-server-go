package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/rpc"
	"os"
	"tc/stat"
	"testing"
)

func TestCall(t *testing.T) {
	t.Log("WARNING: restart Counter before testing")
	client, err := rpc.Dial(config.net, config.laddr)
	defer client.Close()
	if err != nil {
		t.Fatal("dialing:", err)
	}

	test := func(n int) {
		raw := stat.NewRawEvent(stat.ShowEvent)

		err = client.Call("Counter.AddRaw", raw, nil)
		if err != nil {
			t.Fatal("Count error:", err)
		}

		reply := stat.NewEvent(raw)
		err = client.Call("Counter.Get", raw.Key(), reply)
		if err != nil {
			t.Fatal("Get error:", err)
		}

		if reply == nil {
			t.Fatal("Stat not found")
		}
		if int(reply.ShowCount) != n {
			t.Logf("%+v\n", reply)
			t.Fatal("Show count != ", n)
		}
	}
	test(1)
	test(2)

}

func TestMergingFile(t *testing.T) {
	server := NewServer()

	file, err := os.Open("../stat/raw-test.data")
	check(err)
	dec := json.NewDecoder(file)
	i := 0
	for {
		var raw stat.RawEvent
		err := dec.Decode(&raw)
		if err == io.EOF {
			break
		}
		check(err)
		server.AddRaw(raw, nil)
		i++
	}
	const M = 61512
	if i != M {
		t.Fatalf("Expected %d raw records, got %d\n", M, i)
	}
	server.counter.stop()
	const N = 23952
	l := server.counter.events.Len()
	if l != N {
		t.Fatalf("Expected len=%d got %d\n", N, l)
	}
}

func BenchmarkAdding(b *testing.B) {
	client, err := rpc.Dial(config.net, config.laddr)
	if err != nil {
		b.Fatal("dialing:", err)
	}
	defer client.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			raw := &stat.RawEvent{Type: stat.ShowEvent}
			raw.Site.Id = rand.Int()
			err = client.Call("Counter.Add", raw, nil)
			if err != nil {
				b.Fatal("Get error:", err)
			}
		}
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
