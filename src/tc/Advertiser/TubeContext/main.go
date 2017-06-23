package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

const DefaultTimeout = 600 * time.Millisecond

func main() {
	pid := os.Getpid()
	tmpDir := os.TempDir()
	tmpFileName := tmpDir + "/TubeContext.pid"
	err := ioutil.WriteFile(tmpFileName, []byte(strconv.Itoa(pid)), os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Printf("Write pid %v to file %v", pid, tmpFileName)
	//	os.Exit(1)

	cpu := runtime.NumCPU()
	log.Println("CPU:", cpu)
	runtime.GOMAXPROCS(cpu)

	l, err := time.LoadLocation("Etc/GMT-3")

	if err != nil {
		panic(err)
	}

	time.Local = l
	log.Println("Time:", time.Now())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	listenViewProfile()
	NewCommandController()
	NewController()

	log.Println("Server is ready")

	s := <-c
	log.Println("Got signal:", s)
	log.Println("Shuting down...")
}

func listenViewProfile() {
	var netprofile = flag.Bool(
		"netprofile",
		true,
		"record profile; see http://server:6060/debug/prof",
	)
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	var memprofile = flag.String("memprofile", "", "write memory profile to this file")
	flag.IntVar(&measureDefaultRate, "mdr", measureDefaultRate, "measure default rate")
	flag.Parse()

	if *netprofile {
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
	}

	//	pprof.StopCPUProfile()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)

		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()

		return
	}
}
