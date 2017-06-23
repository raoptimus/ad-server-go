package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type (
	config struct {
		sync.RWMutex
		data *configData
	}
	configData map[string]interface{}
)

var cfg *config

func MongoServerMain() string {
	return String("MongoServerMain", "")
}

func MongoServerAdIndex() string {
	return String("MongoServerAdIndex", "")
}

func ServerName() string {
	return *flag.String("server", String("ServerName", ""), "server name")
}

func PgSqlServer() string {
	return String("PgSqlServer", "")
}

func EnableProfile() bool {
	return Bool("enableProfile", false)
}

func String(name, value string) string {
	v := cfg.get(name)
	if v == nil {
		return value
	}
	value = fmt.Sprintf("%v", v)
	return value
}

func Bool(name string, value bool) bool {
	v := cfg.get(name)
	if v == nil {
		return value
	}
	s := fmt.Sprintf("%v", v)
	if s == "1" || s == "true" {
		value = true
	} else if s == "0" || s == "false" {
		value = false
	}

	return value
}

func Int(name string, value int) int {
	v := cfg.get(name)
	if v == nil {
		return value
	}
	s := fmt.Sprintf("%v", v)
	n, err := strconv.Atoi(s)
	if err != nil {
		value = n
	}
	return value
}

func (s *config) get(name string) interface{} {
	s.RLock()
	defer s.RUnlock()

	v, ok := (*s.data)[name]
	if ok {
		return v
	}

	return nil
}

func loadConfig(cfgFileName *string) (*configData, error) {
	file, err := ioutil.ReadFile(*cfgFileName)
	if err != nil {
		return nil, errors.New("Error open config file '" + *cfgFileName + "': " + err.Error())
	}

	var d configData
	if err = json.Unmarshal(file, &d); err != nil {
		return nil, errors.New("Error parse config file '" + *cfgFileName + "': " + err.Error())
	}

	fmt.Println("Config is loaded", *cfgFileName)
	return &d, nil
}

func init() {
	var cfgFileName string
	flag.StringVar(&cfgFileName, "config", "config/release.json", "Configuration file options")
	flag.Parse()

	d, err := loadConfig(&cfgFileName)
	if err != nil {
		log.Fatalln(err)
	}
	cfg = &config{
		data: d,
	}
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGUSR2)
	go func() {
		for {
			<-s
			d, err = loadConfig(&cfgFileName)
			if err != nil {
				log.Println(err)
				return
			}
			cfg.Lock()
			defer cfg.Unlock()

			cfg.data = d
		}
	}()
}
