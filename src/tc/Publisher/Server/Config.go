package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type Config struct {
	MongoServerMain        string
	MongoServerAdIndex     string
	MongoServerSession     string
	SocketRequestTimeoutMs int
}

var (
	config     *Config
	configLock = new(sync.RWMutex)
)

func loadConfig(name string) {
	file, err := ioutil.ReadFile(name + ".json")

	if err != nil {
		log.Println("open config: ", err)
		os.Exit(1)
	}

	temp := new(Config)

	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		os.Exit(1)
	}

	configLock.Lock()
	config = temp
	configLock.Unlock()
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()

	return config
}

/** TODO for reload
func init() {
  loadConfig(true)
  s := make(chan os.Signal, 1)
  signal.Notify(s, syscall.SIGUSR2)
  go func() {
    for {
      <-s
      loadConfig(false)
      log.Println("Reloaded")
    }
  }()
}
*/
