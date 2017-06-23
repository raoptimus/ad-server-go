package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"tc/tc"
	"time"
)

//start the program
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Starting server...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go listenProfile()
	go listenServe()

	log.Println("Server is ready.")

	s := <-c

	log.Println("Got signal:", s)
	log.Println("Shutting down...")
}

func listenServe() {
	server, err := net.Listen("tcp", "127.0.0.1:9082")

	if err != nil {
		panic(err)
	}

	s := new(Controller)
	rpc.Register(s)
	tc.GobRegisterExt()

	for {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				log.Println("PANIC:", r)
				return
			}
		}()

		conn, err := server.Accept()
		defer conn.Close()
		conn.SetDeadline(time.Now().Add(600 * time.Millisecond))

		if err != nil {
			log.Fatal(err)
		}

		rpc.ServeConn(conn)
	}
}

//
//func listenServe() {
//	l, err := net.Listen("tcp", "127.0.0.1:9082")
//
//	if err != nil {
//		panic(err)
//	}
//
//	defer l.Close()
//
//	for {
//		// Wait for a connection.
//		conn, err := l.Accept()
//
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		go func(conn net.Conn) {
//			defer conn.Close()
//			conn.SetDeadline(time.Now().Add(600 * time.Millisecond))
//
//			defer func() {
//				if r := recover(); r != nil {
//					debug.PrintStack()
//					log.Println("PANIC:", r)
//					return
//				}
//			}()
//
//			var tcReq *openrtb.Request
//			dec := json.NewDecoder(conn)
//
//			if err := dec.Decode(&tcReq); err != nil {
//				log.Println("Request decoder error:", err)
//				return
//			}
//
//			tcResp := &openrtb.Response{}
//			err := new(Controller).SendRequest(tcReq, tcResp)
//
//			if err != nil {
//				if *tcReq.Device.Ip != "..." && skipHttpError(err) {
//					return
//				}
//
//				log.Println("Request encode error:", err)
//				return
//			}
//
//			enc := json.NewEncoder(conn)
//			enc.Encode(tcResp)
//		}(conn)
//	}
//}
//
////got bid winners
//func listenConfirmRequests() {
//	l, err := net.Listen("tcp", "127.0.0.1:9083")
//
//	if err != nil {
//		panic(err)
//	}
//
//	defer l.Close()
//
//	for {
//		// Wait for a connection.
//		conn, err := l.Accept()
//
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		go func(conn net.Conn) {
//			defer conn.Close()
//			conn.SetDeadline(time.Now().Add(300 * time.Millisecond))
//
//			var confirm tc.Confirm
//			dec := json.NewDecoder(conn)
//
//			if err := dec.Decode(&confirm); err != nil {
//				log.Println("Request decoder error:", err)
//			}
//
//			if len(confirm.Results) <= 0 {
//				return
//			}
//
//			go new(Controller).SendConfirm(&confirm)
//		}(conn)
//	}
//}

func skipHttpError(err error) bool {
	return strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "request canceled") ||
		strings.Contains(err.Error(), "closed network connection") ||
		strings.Contains(err.Error(), "failure in name resolution")
}

//download profiles:
//go tool pprof http://localhost:6061/debug/pprof/profile   # 30-second CPU profile
//go tool pprof http://localhost:6061/debug/pprof/heap      # heap profile
//go tool pprof http://localhost:6061/debug/pprof/block     # goroutine blocking profile
//
//see profile in browser:
//http://localhost:6061/debug/pprof/
//
func listenProfile() {
	var netprofile = flag.Bool(
		"netprofile",
		true,
		"record profile; see http://localhost:6061/debug/prof",
	)

	var cpuprofile = flag.String(
		"cpuprofile",
		"",
		"write cpu profile to file",
	)

	flag.Parse()

	if *netprofile {
		go func() {
			log.Println(http.ListenAndServe(":7061", nil))
		}()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)

		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
}
