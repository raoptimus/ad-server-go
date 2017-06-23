package stat

import (
	"bytes"
	"encoding/gob"
	"sync"
)

type (
	Map struct {
		sync.RWMutex
		data eventsMap
	}
	Slice     []*Event
	eventsMap map[EventKey]*Event
)

func NewMap() *Map {
	return &Map{
		data: make(eventsMap),
	}
}

//AddRaw adds raw event to map
//In case of unique event it will increment both uniq and non uniq counters
func (m *Map) AddRaw(raw *RawEvent) {
	key := raw.Key()
	e := m.Get(key)
	if e != nil {
		e.AddRaw(raw)
		return
	}

	e = NewEvent(raw)
	m.Lock()
	m.data[key] = e
	m.Unlock()
}

func (m *Map) Get(key EventKey) *Event {
	m.RLock()
	e, _ := m.data[key]
	m.RUnlock()
	return e
}

func (m *Map) Len() int {
	m.RLock()
	l := len(m.data)
	m.RUnlock()
	return l
}

func (m *Map) Merge(with *Map) {
	with.RLock()
	for _, e := range with.data {
		m.Add(e)
	}
	with.RUnlock()
}

//Add adds e counters to map
func (m *Map) Add(e *Event) {
	key := e.Key()
	existing := m.Get(key)
	if existing != nil {
		existing.Add(e)
		return
	}
	m.Lock()
	m.data[key] = e
	m.Unlock()
}

func (m *Map) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	m.RLock()
	enc.Encode(m.data)
	m.RUnlock()
	return b.Bytes(), nil
}

func (m *Map) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	m.Lock()
	defer m.Unlock()
	return dec.Decode(&m.data)
}

func (m *Map) Slice() Slice {
	m.Lock()
	slice := make(Slice, 0, len(m.data))
	for _, e := range m.data {
		slice = append(slice, e)
	}
	m.Unlock()
	return slice
}
