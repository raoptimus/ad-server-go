package stat

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
)

type testFields map[string]int64

var m = NewMap()

func testMapAddRaw(t *testing.T, etype EventType, fields testFields) {
	raw := &RawEvent{Type: etype}
	m.AddRaw(raw)
	stat := m.Get(raw.Key())
	if stat == nil {
		t.Fatal("Stat not found")
	}

	v := reflect.ValueOf(stat).Elem()
	for field, expected := range fields {
		got := v.FieldByName(field).Int()
		if got != expected {
			t.Fatalf("Got %s = %d; expected = %d", field, got, expected)
		}
	}
}

func TestMapAddRaw(t *testing.T) {
	testMapAddRaw(t, ClickEvent, testFields{"ClickCount": 1})
	testMapAddRaw(t, ClickEvent, testFields{"ClickCount": 2})
	testMapAddRaw(t, UniqueClickEvent, testFields{"UniqueClickCount": 1, "ClickCount": 3})
	testMapAddRaw(t, UniqueClickEvent, testFields{"UniqueClickCount": 2, "ClickCount": 4})
	testMapAddRaw(t, DoubleClickEvent, testFields{"DoubleClickCount": 1, "ClickCount": 4})
	testMapAddRaw(t, DoubleClickEvent, testFields{"DoubleClickCount": 2, "ClickCount": 4})
	testMapAddRaw(t, BadClickEvent, testFields{"BadClickCount": 1, "ClickCount": 4})
	testMapAddRaw(t, BadClickEvent, testFields{"BadClickCount": 2, "ClickCount": 4})
	testMapAddRaw(t, ShowEvent, testFields{"ShowCount": 1})
	testMapAddRaw(t, ShowEvent, testFields{"ShowCount": 2})
	testMapAddRaw(t, UniqueShowEvent, testFields{"UniqueShowCount": 1, "ShowCount": 3})
	testMapAddRaw(t, UniqueShowEvent, testFields{"UniqueShowCount": 2, "ShowCount": 4})

}

func TestEncoding(t *testing.T) {
	m := NewMap()
	b, err := m.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	decoded := NewMap()
	if err := decoded.UnmarshalBinary(b); err != nil {
		t.Fatal(err)
	}
	if m.Len() != decoded.Len() {
		t.Fatal("Len is not equal")
	}
	for _, e := range decoded.data {
		stat := m.Get(e.Key())
		if stat == nil {
			t.Fatal(e.Key(), " not found")
		}
	}
}

func TestMerging(t *testing.T) {
	m := NewMap()
	addendum := NewMap()
	N := rand.Int63() % 100
	for i := 0; i < int(N); i++ {
		addendum.AddRaw(&RawEvent{Type: ClickEvent})
		addendum.AddRaw(&RawEvent{Type: UniqueClickEvent})
		addendum.AddRaw(&RawEvent{Type: DoubleClickEvent})
		addendum.AddRaw(&RawEvent{Type: BadClickEvent})
		addendum.AddRaw(&RawEvent{Type: ShowEvent})
		addendum.AddRaw(&RawEvent{Type: UniqueShowEvent})
	}

	m.Merge(addendum)
	stat := m.Get(EventKey{})
	if stat == nil {
		t.Fatal("Stat not found after merging")
	}

	switch {
	case stat.ClickCount != 2*N: //+uniq
		t.Fatal("Merge failed: clicks")
	case stat.UniqueClickCount != N:
		t.Fatal("Merge failed: unique clicks")
	case stat.DoubleClickCount != N:
		t.Fatal("Merge failed: double clicks")
	case stat.BadClickCount != N:
		t.Fatal("Merge failed: bad clicks")
	case stat.ShowCount != 2*N:
		t.Fatal("Merge failed: shows")
	case stat.UniqueShowCount != N:
		t.Fatal("Merge failed: shows")
	}

}

func TestMergingFile(t *testing.T) {
	file, err := os.Open("raw-test.data")
	check(err)
	dec := json.NewDecoder(file)
	m := NewMap()
	i := 0
	for {
		var raw RawEvent
		err := dec.Decode(&raw)
		if err == io.EOF {
			break
		}
		check(err)
		raw.Date = raw.Date.Local().Truncate(time.Hour)
		m.AddRaw(&raw)
		i++
	}
	const M = 61512
	if i != M {
		t.Fatalf("Expected %d raw records, got %d\n", M, i)
	}
	const N = 23952
	if m.Len() != N {
		t.Fatalf("Expected len=%d got %d\n", N, m.Len())
	}
}

func BenchmarkAdding(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.AddRaw(randomRawEvent())
		}
	})
}

// var with = randomMap()

func BenchmarkMerging(b *testing.B) {
	with := randomMap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Merge(with)
	}
	// b.Log("Map size:", m.Len())
}

func BenchmarkMergingParalell(b *testing.B) {
	with := randomMap()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Merge(with)
		}
	})
	// b.Log("Map size:", m.Len())
}

func randomRawEvent() *RawEvent {
	t := EventType(1 + rand.Intn(EventsNum))
	raw := NewRawEvent(t)
	raw.Site.Id = rand.Int()
	raw.Campaign.Ad.Id = rand.Int()
	raw.Campaign.Id = rand.Int()
	return raw
}

func randomMap() *Map {
	m := NewMap()
	for i := 0; i < 100000; i++ {
		m.AddRaw(randomRawEvent())
	}
	return m
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
