package stat

import (
	"log"
	"sync/atomic"
	tc "tc/openrtbex"
	"time"
)

const (
	// shows
	ShowEvent EventType = 1 + iota
	UniqueShowEvent

	FreeShowEvent
	UniqueFreeShowEvent

	WmShowEvent
	UniqueWmShowEvent
	// clicks
	ClickEvent
	UniqueClickEvent

	FreeClickEvent
	UniqueFreeClickEvent

	WmClickEvent
	UniqueWmClickEvent

	// nonunique
	DoubleClickEvent
	BadClickEvent

	EventsNum = iota - 1
)

const precision = 1e6

type (
	EventType int
	EventKey  struct {
		Site struct {
			Id        int
			UserId    int
			RefUserId int
			Category  tc.Category

			AdZone struct {
				Id   int
				Type tc.AdZoneType

				AdCode struct {
					Id   int
					Type tc.AdCodeType
				}
			}
		}
		Campaign struct {
			Id          int
			UserId      int
			BrokerId    int
			Type        tc.CampaignType
			IsWebmaster bool
			IsFree      bool

			Ad struct {
				Id int
			}
		}
		Geo struct {
			Id        int
			CountryId int
		}

		Date     time.Time
		Server   string
		Operator tc.Operator
		Os       tc.Os
		Device   tc.Device
		Browser  tc.Browser
	}
	RawEvent struct {
		EventKey    `bson:",inline"`
		PaymentType tc.PaymentType
		Type        EventType
		FeedRequest bool
		Money       struct {
			Price     int64
			Tax       int64
			Referrals int64
		} `table:"-"`
	}
	RawEvents []RawEvent

	// `Counters' must be in the beginning of the struct which embeds it
	// see BUGS section in `godoc sync/atomic
	Counters struct {
		// shows
		ShowCount       int64
		UniqueShowCount int64

		FreeShowCount       int64
		UniqueFreeShowCount int64

		WmShowCount       int64
		UniqueWmShowCount int64
		// clicks
		ClickCount       int64
		UniqueClickCount int64

		FreeClickCount       int64
		UniqueFreeClickCount int64

		WmClickCount       int64
		UniqueWmClickCount int64
		//nonuniq
		DoubleClickCount int64
		BadClickCount    int64

		RequestCount int64 `table:"-"`

		//money
		Price     int64 `table:"-"`
		Tax       int64 `table:"-"`
		Referrals int64 `table:"-"`
	}
	Event struct {
		Counters `bson:",inline"`
		RawEvent `bson:",inline"` //must be last due to alignment
	}
)

func Float(i int64) float64 {
	return float64(i) / precision
}

func (t EventType) isFree() bool {
	switch t {
	case FreeShowEvent, UniqueFreeShowEvent, FreeClickEvent, UniqueFreeClickEvent:
		return true
	}
	return false
}

func (t EventType) isWebmaster() bool {
	switch t {
	case WmShowEvent, UniqueWmShowEvent, WmClickEvent, UniqueWmClickEvent:
		return true
	}
	return false
}

func NewRawEvent(t EventType) *RawEvent {
	raw := &RawEvent{Type: t}
	raw.Date = time.Now().Local().Truncate(time.Hour)
	return raw
}

func (raw *RawEvent) Key() EventKey {
	return raw.EventKey
}

func (raw *RawEvent) IsBad() bool {
	switch raw.Type {
	case BadClickEvent, DoubleClickEvent:
		return true
	}
	return false
}

func NewEvent(raw *RawEvent) *Event {
	e := &Event{}
	if raw != nil {
		e.RawEvent = *raw
		e.AddRaw(raw)
	}
	e.Campaign.IsFree = e.Type.isFree()
	e.Campaign.IsWebmaster = e.Type.isWebmaster()
	return e
}

//In case of unique event it will increment both uniq and non uniq counters
func (e *Event) AddRaw(raw *RawEvent) {
	e.add(raw.Type, 1)

	switch raw.Type {
	case UniqueClickEvent: // clicks
		e.add(ClickEvent, 1)
	case UniqueFreeClickEvent:
		e.add(FreeClickEvent, 1)
	case UniqueWmClickEvent:
		e.add(WmClickEvent, 1)
	case UniqueShowEvent: //shows
		e.add(ShowEvent, 1)
	case UniqueFreeShowEvent:
		e.add(FreeShowEvent, 1)
	case UniqueWmShowEvent:
		e.add(WmShowEvent, 1)
	}

	if raw.FeedRequest {
		atomic.AddInt64(&e.RequestCount, 1)
	}
	if raw.IsBad() {
		return
	}
	atomic.AddInt64(&e.Counters.Price, raw.Money.Price)
	atomic.AddInt64(&e.Counters.Tax, raw.Money.Tax)
	atomic.AddInt64(&e.Counters.Referrals, raw.Money.Referrals)
}

func (e *Event) Add(what *Event) {
	e.Counters.Add(&what.Counters)
}

func (e *Event) IsPaid() bool {
	return !e.Campaign.IsWebmaster && !e.Campaign.IsFree
}

func (c *Counters) Add(what *Counters) {
	// shows
	c.add(ShowEvent, what.ShowCount)
	c.add(UniqueShowEvent, what.UniqueShowCount)

	c.add(FreeShowEvent, what.FreeShowCount)
	c.add(UniqueFreeShowEvent, what.UniqueFreeShowCount)

	c.add(WmShowEvent, what.WmShowCount)
	c.add(UniqueWmShowEvent, what.UniqueWmShowCount)

	// clicks
	c.add(ClickEvent, what.ClickCount)
	c.add(UniqueClickEvent, what.UniqueClickCount)

	c.add(FreeClickEvent, what.FreeClickCount)
	c.add(UniqueFreeClickEvent, what.UniqueFreeClickCount)

	c.add(WmClickEvent, what.WmClickCount)
	c.add(UniqueWmClickEvent, what.UniqueWmClickCount)

	//nonuniq
	c.add(DoubleClickEvent, what.DoubleClickCount)
	c.add(BadClickEvent, what.BadClickCount)

	atomic.AddInt64(&c.RequestCount, what.RequestCount)
	atomic.AddInt64(&c.Price, what.Price)
	atomic.AddInt64(&c.Tax, what.Tax)
	atomic.AddInt64(&c.Referrals, what.Referrals)
}

func (c *Counters) add(t EventType, n int64) {
	if n == 0 {
		return
	}
	if n < 0 {
		log.Println("Negative n")
		return
	}
	// counter := strings.TrimRight(reflect.TypeOf(t).Name(), "Event") + "Count"
	// i := reflect.ValueOf(c).FieldByName(counter).Addr().Interface().(*int64)
	var i *int64
	switch t {
	// shows
	case ShowEvent:
		i = &c.ShowCount
	case UniqueShowEvent:
		i = &c.UniqueShowCount

	case FreeShowEvent:
		i = &c.FreeShowCount
	case UniqueFreeShowEvent:
		i = &c.UniqueFreeShowCount

	case WmShowEvent:
		i = &c.WmShowCount
	case UniqueWmShowEvent:
		i = &c.UniqueWmShowCount

		// clicks
	case ClickEvent:
		i = &c.ClickCount
	case UniqueClickEvent:
		i = &c.UniqueClickCount

	case FreeClickEvent:
		i = &c.FreeClickCount
	case UniqueFreeClickEvent:
		i = &c.UniqueFreeClickCount

	case WmClickEvent:
		i = &c.WmClickCount
	case UniqueWmClickEvent:
		i = &c.UniqueWmClickCount
		// nonuniq
	case DoubleClickEvent:
		i = &c.DoubleClickCount
	case BadClickEvent:
		i = &c.BadClickCount
	default:
		log.Println("Unknown event type:", t)
		return
	}
	atomic.AddInt64(i, n)
}
