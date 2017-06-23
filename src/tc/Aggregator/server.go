package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"tc/stat"
	"time"
)

type Server struct {
	server     net.Listener
	merger     *merger
	aggregator *aggregator
	runDump    chan os.Signal
}

func newServer() *Server {
	return &Server{
		merger:     newMerger(),
		aggregator: newAggregator(),
		runDump:    make(chan os.Signal, 1),
	}
}

//separated for testing

//run the server
func (s *Server) run() {
	go s.listenAndServe()
	go s.dumper()
}

func (s *Server) listenAndServe() {
	var err error
	s.server, err = net.Listen(config.net, config.laddr)
	if err != nil {
		panic(err)
	}

	rpc.RegisterName("Aggregator", s)
	log.Println("Server ready")

	for {
		conn, err := s.server.Accept()

		if err != nil {
			log.Println(err)
			return
		}

		go func(conn net.Conn) {
			defer catchPanic()
			defer conn.Close()
			rpc.ServeConn(conn)
		}(conn)
	}
}

func (s *Server) Merge(m stat.Map, _ *int) error {
	s.merger.push(&m)
	return nil
}

func (s *Server) dumper() {
	ticker := time.Tick(config.dumpInterval)
	for {
		select {
		case <-ticker:
		case <-s.runDump:
		}
		s.dump()
	}
}

func (s *Server) dump() {
	merger := s.merger

	log.Printf("Merger queue size: %d\n", len(merger.queue))

	wait := time.Second
	for len(merger.queue) == cap(merger.queue) {
		time.Sleep(wait)
		if wait < time.Minute {
			wait *= 2
		}
	}

	s.merger = newMerger()
	s.aggregator.dump(merger)
}
