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
	"tc/openrtbex"
	"time"
)

//start the program
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Starting server...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go listenViewProfile()
	NewController()

	log.Println("Server is ready.")

	s := <-c

	log.Println("Got signal:", s)
	log.Println("Shutting down...")
}

func skipHttpError(err error) bool {
	return strings.Contains(err.Error(), "i/o timeout") ||
		strings.Contains(err.Error(), "request canceled") ||
		strings.Contains(err.Error(), "closed network connection") ||
		strings.Contains(err.Error(), "failure in name resolution")
}

func listenServe() {
	server, err := net.Listen("tcp", "127.0.0.1:8082")

	if err != nil {
		panic(err)
	}

	s := new(Controller)
	rpc.Register(s)
	openrtbex.GobRegisterExt()

	for {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				log.Println("PANIC:", r)
				return
			}
		}()

		conn, err := server.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					debug.PrintStack()
					log.Println("PANIC:", r)
					return
				}
			}()
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(600 * time.Millisecond))
			rpc.ServeConn(conn)
		}(conn)
	}
}

//download profiles:
//go tool pprof http://localhost:6061/debug/pprof/profile   # 30-second CPU profile
//go tool pprof http://localhost:6061/debug/pprof/heap      # heap profile
//go tool pprof http://localhost:6061/debug/pprof/block     # goroutine blocking profile
//
//see profile in browser:
//http://localhost:6061/debug/pprof/
//
func listenViewProfile() {
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
			log.Println(http.ListenAndServe(":6061", nil))
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
