package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
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

//download profiles:
//go tool pprof http://localhost:6062/debug/pprof/profile   # 30-second CPU profile
//go tool pprof http://localhost:6062/debug/pprof/heap      # heap profile
//go tool pprof http://localhost:6062/debug/pprof/block     # goroutine blocking profile
//
//see profile in browser:
//http://localhost:6062/debug/pprof/
//
func listenViewProfile() {
	var netprofile = flag.Bool(
		"netprofile",
		true,
		"record profile; see http://localhost:6062/debug/prof",
	)

	var cpuprofile = flag.String(
		"cpuprofile",
		"",
		"write cpu profile to file",
	)

	flag.Parse()

	if *netprofile {
		go func() {
			log.Println(http.ListenAndServe(":6062", nil))
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
