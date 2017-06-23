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

func main() {
	cpu := runtime.NumCPU()
	log.Println("CPU:", cpu)
	runtime.GOMAXPROCS(cpu)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go listenProfile()
	NewApiController(NewController())

	log.Println("Server ready")

	s := <-c
	log.Println("Got signal:", s)
	log.Println("Shuting down...")
}

func listenProfile() {
	var netprofile = flag.Bool(
		"netprofile",
		false,
		"record profile; see http://server:6000/debug/prof",
	)
	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memProfile = flag.String("memprofile", "", "write memory profile to this file")
	flag.IntVar(&measureDefaultRate, "mdr", measureDefaultRate, "measure default rate")
	flag.Parse()

	if !*netprofile {
		return
	}

	go func() {
		log.Println(http.ListenAndServe(":6000", nil))
	}()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
	}

	//	pprof.StopCPUProfile()

	if *memProfile != "" {
		f, err := os.Create(*memProfile)

		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()

		return
	}
}
