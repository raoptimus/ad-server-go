package bootstrap

import (
	"log"
	"os"
	"os/signal"
)

func Run(main func()) {
	Init()
	go main()
	Wait()
}

func Wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	log.Println("Got signal:", <-c)
	log.Println("Shutting down...")
}
