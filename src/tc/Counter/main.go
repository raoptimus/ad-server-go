package main

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"tc/bootstrap"
	"tc/openrtbex"
	"tc/stat"
	"time"
)

var config = struct {
	net           string
	laddr         string
	rawQueueSize  int
	sendQueueSize int
	mongo         string
}{
	"tcp",
	"127.0.0.1:9990",
	1000,
	10,
	"d1444.webazilla.com/serverInfo",
}

type Server struct {
	mongo             *mgo.Database
	client            *rpc.Client
	server            net.Listener
	counter           *counter
	sendQueue         chan *stat.Map
	serverEventsQueue chan *stat.Map
}

func NewServer() *Server {
	ss, err := mgo.Dial(config.mongo)
	if err != nil {
		panic(err)
	}
	return &Server{
		counter:           newCounter(),
		sendQueue:         make(chan *stat.Map, config.sendQueueSize),
		serverEventsQueue: make(chan *stat.Map, config.sendQueueSize),
		mongo:             ss.DB(""),
	}
}

func (s *Server) AddRaw(raw stat.RawEvent, unused *int) error {
	s.counter.push(raw)
	return nil
}

func (s *Server) listenAndServe() {
	var err error
	s.server, err = net.Listen(config.net, config.laddr)

	if err != nil {
		panic(err)
	}

	rpc.RegisterName("Counter", s)
	go s.listenAndServeJson()
	log.Println("Server ready")

	for {
		conn, err := s.server.Accept()

		if err != nil {
			log.Println("Accept error:", err)
		}

		go func(conn net.Conn) {
			defer catchPanic()
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(6000 * time.Millisecond))
			rpc.ServeConn(conn)
		}(conn)
	}
}
func (s *Server) listenAndServeJson() {
	server, err := net.Listen(config.net, ":9991")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Println("Accept error:", err)
		}

		go func(conn net.Conn) {
			defer catchPanic()
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(6000 * time.Millisecond))
			jsonrpc.ServeConn(conn)
		}(conn)
	}
}

func (s *Server) dumper() {
	go s.sender()
	go s.serverEventsSaver()
	dump := make(chan os.Signal)
	signal.Notify(dump, syscall.SIGUSR1)
	ticker := time.Tick(1 * time.Minute)
	for {
		select {
		case <-ticker:
		case <-dump:
		}
		s.dump()
	}
}

func (s *Server) dump() {
	wait := time.Second
	for len(s.sendQueue) > cap(s.sendQueue)/2 {
		time.Sleep(wait)
		if wait < time.Minute {
			wait *= 2
		}
	}
	c := s.counter
	s.counter = newCounter()
	c.stop()
	log.Printf("Counted %d raw events\n", c.counted)
	s.sendQueue <- c.events
	s.serverEventsQueue <- c.serverEvents
}

func (s *Server) sender() {
	for m := range s.sendQueue {
		s.send(m)
	}
}

func (s *Server) serverEventsSaver() {
	for m := range s.serverEventsQueue {
		s.saveServerEvents(m)
	}
}

func (s *Server) send(m *stat.Map) {
	log.Printf("Sending %d events to the Aggregator\n", m.Len())
	wait := time.Second
	for {
		if s.client == nil {
			s.connect()
		}
		err := s.client.Call("Aggregator.Merge", m, nil)
		if err == nil {
			break
		}
		log.Println("Call error:", err)
		s.client.Close()
		s.client = nil
		time.Sleep(wait)
		if wait < time.Minute {
			wait *= 2
		}
	}
	log.Println("Sending done")
}

func (s *Server) connect() {
	var err error
	wait := time.Second
	for {
		log.Println("Connecting to aggregator...")
		s.client, err = rpc.Dial("tcp", "127.0.0.1:8880")
		if err == nil {
			log.Println("...connected")
			return
		}
		log.Println("Cannot connect to aggregator:", err)
		time.Sleep(wait)
		if wait < time.Minute {
			wait *= 2
		}
	}
}

func (s *Server) saveServerEvents(m *stat.Map) {
	type tuple struct {
		field string
		value int64
	}

	for _, e := range m.Slice() {
		set := bson.M{
			"Server":      e.Server,
			"RequestDate": e.Date,
		}
		showCount := e.ShowCount + e.FreeShowCount + e.WmShowCount
		uniqueShowCount := e.UniqueShowCount + e.UniqueFreeShowCount + e.UniqueWmShowCount
		totalClickCount := e.ClickCount + e.FreeClickCount + e.WmClickCount + e.BadClickCount + e.DoubleClickCount
		uniqueClickCount := e.UniqueClickCount + e.UniqueFreeClickCount + e.UniqueWmClickCount

		inc := bson.M{
			"ReqFeedCount": e.RequestCount,
		}

		fkey := func(field string) string {
			return fmt.Sprintf("Friends.%d.%s", e.Campaign.BrokerId, field)
		}
		tkey := func(field string) string {
			typeId := e.Site.AdZone.AdCode.Type.ToInt()
			return fmt.Sprintf("Types.%d.%s", typeId, field)
		}
		tfkey := func(field string) string {
			return fkey(tkey(field))
		}

		add := func(tuples ...tuple) {
			for _, tuple := range tuples {
				if tuple.value == 0 {
					continue
				}
				inc[tuple.field] = tuple.value
				inc[fkey(tuple.field)] = tuple.value
				inc[tkey(tuple.field)] = tuple.value
				inc[tfkey(tuple.field)] = tuple.value
			}
		}
		add(
			tuple{"ShowCount", showCount},
			tuple{"UniqueShowCount", uniqueShowCount},
			tuple{"TotalClickCount", totalClickCount},
			tuple{"UniqueClickCount", uniqueClickCount},
			tuple{"BadClickCount", e.BadClickCount},
			tuple{"FreeClickCount", e.FreeClickCount},
			tuple{"DoubleClickCount", e.DoubleClickCount},
			tuple{"SelfClickCount", e.WmClickCount},
		)

		if e.PaymentType == openrtbex.PaymentTypeCPM {
			add(
				tuple{"PayShowCount", e.ShowCount},
				tuple{"NoPayClickCount", e.ClickCount},
			)
		} else {
			add(tuple{"ClickCount", e.ClickCount})
		}

		//TODO: use e.IsPaid()?
		sumCost := stat.Float(e.Price + e.Tax + e.Referrals)
		inc["SumCost"] = sumCost
		inc[fkey("SumCost")] = sumCost
		inc[tkey("SumCost")] = sumCost
		inc[tfkey("SumCost")] = sumCost

		b := bson.M{}
		b["$setOnInsert"] = set
		b["$inc"] = inc

		timeBlock := e.Date.Local().Truncate(5 * time.Minute).Unix()
		key := fmt.Sprintf("%s|%d", e.Server, timeBlock)
		id := fmt.Sprintf("%x", md5.Sum([]byte(key)))

		_, err := s.mongo.C("traffic_byservers").UpsertId(id, b)
		if err != nil {
			s.mongo.Session.Refresh()
			log.Println(err)
		}
	}
}

func (s *Server) shutdown() {
	s.server.Close()
}

func main() {
	bootstrap.Init()

	server := NewServer()
	go server.listenAndServe()
	go server.dumper()
	defer server.shutdown()

	bootstrap.Wait()
}

func catchPanic() {
	if r := recover(); r != nil {
		debug.PrintStack()
		log.Println("PANIC:", r)
		return
	}
}
