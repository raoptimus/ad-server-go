package main

import (
	"log"
	"runtime/debug"
	"tc/bootstrap"
)

func main() {
	bootstrap.Run(func() {
		storage := NewStorage()
		NewApiController(storage)
	})
}

func recovery() {
	if r := recover(); r != nil {
		debug.PrintStack()
		log.Println("PANIC:", r)
		return
	}
}
