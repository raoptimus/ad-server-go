package stat

import (
	"testing"
)

func TestEventAddingRaw(t *testing.T) {
	e := NewEvent(nil)
	if e.ClickCount != 0 {
		t.Fatalf("Got ClickCount = %d expected 0", e.ClickCount)
	}
	if e.UniqueClickCount != 0 {
		t.Fatalf("Got UniqueClickCount = %d expected 0", e.UniqueClickCount)
	}
	if e.DoubleClickCount != 0 {
		t.Fatalf("Got DoubleClickCount = %d expected 0", e.DoubleClickCount)
	}
	if e.BadClickCount != 0 {
		t.Fatalf("Got BadClickCount = %d expected 0", e.BadClickCount)
	}
	if e.ShowCount != 0 {
		t.Fatalf("Got ShowCount = %d expected 0", e.ShowCount)
	}
	if e.UniqueShowCount != 0 {
		t.Fatalf("Got UniqueShowCount = %d expected 0", e.UniqueShowCount)
	}

	e.AddRaw(NewRawEvent(ClickEvent))
	if e.ClickCount != 1 {
		t.Fatalf("Got ClickCount = %d expected 1", e.ClickCount)
	}

	e.AddRaw(NewRawEvent(UniqueClickEvent))
	if e.UniqueClickCount != 1 {
		t.Fatalf("Got UniqueClickCount = %d expected 1", e.UniqueClickCount)
	}
	if e.ClickCount != 2 {
		t.Fatalf("Got ClickCount = %d expected 2", e.ClickCount)
	}

	e.AddRaw(NewRawEvent(DoubleClickEvent))
	if e.DoubleClickCount != 1 {
		t.Fatalf("Got DoubleClickCount = %d expected 1", e.DoubleClickCount)
	}

	e.AddRaw(NewRawEvent(BadClickEvent))
	if e.DoubleClickCount != 1 {
		t.Fatalf("Got BadClickCount = %d expected 1", e.BadClickCount)
	}

	e.AddRaw(NewRawEvent(ShowEvent))
	if e.ShowCount != 1 {
		t.Fatalf("Got ShowCount = %d expected 1", e.ShowCount)
	}

	e.AddRaw(NewRawEvent(UniqueShowEvent))
	if e.UniqueShowCount != 1 {
		t.Fatalf("Got UniqueShowCount = %d expected 1", e.UniqueShowCount)
	}
	if e.ShowCount != 2 {
		t.Fatalf("Got ShowCount = %d expected 2", e.ShowCount)
	}

	etalon := NewEvent(nil)
	etalon.ClickCount = 2
	etalon.UniqueClickCount = 1
	etalon.DoubleClickCount = 1
	etalon.BadClickCount = 1
	etalon.ShowCount = 2
	etalon.UniqueShowCount = 1

	if *e != *etalon {
		t.Fatalf("Got %+v = expected %+v", e, etalon)
	}
}
