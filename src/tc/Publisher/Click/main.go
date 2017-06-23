package main

import (
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

type Id int

func (id Id) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func main() {
	pid := os.Getpid()
	tmpDir := os.TempDir()
	tmpFileName := tmpDir + "/Click.pid"
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

	NewController()

	s := <-c
	log.Println("Got signal:", s)
	log.Println("Shuting down...")
}
