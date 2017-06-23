package store

import (
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type RedisStorage struct {
	pool   *redis.Pool
	prefix string
	life   time.Duration
}

const (
	MAX_IDLE     = 1024
	IDLE_TIMEOUT = 240 * time.Second
)

//network - "tcp", address - ":6379"
//network - "unix", address - "go fmt"
//life - time.Duration(56700*time.Second)
func NewRedisStorage(dbName, network, address string, life time.Duration) *RedisStorage {
	return &RedisStorage{
		pool: &redis.Pool{
			MaxIdle:     MAX_IDLE,
			IdleTimeout: IDLE_TIMEOUT,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout(
					network,
					address,
					time.Duration(time.Millisecond*50),
					time.Duration(time.Millisecond*100),
					time.Duration(time.Millisecond*100),
				)
				if err != nil {
					return nil, errors.New("Redis dial error: " + err.Error())
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
		prefix: dbName + ".",
		life:   life,
	}
}

func (s *RedisStorage) Get(k string, v interface{}) (err error) {
	conn := s.pool.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", s.prefix+k))

	if err != nil {
		return errors.New("Get data from redis error: " + err.Error())
	}

	err = json.Unmarshal(reply, &v)

	if err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) Set(k string, v interface{}) {
	conn := s.pool.Get()
	defer conn.Close()

	b, err := json.Marshal(v)

	if err != nil {
		log.Println(err)
		return
	}
	conn.Send("SET", s.prefix+k, string(b))
	conn.Send("EXPIRE", s.prefix+k, time.Now().Add(s.life))
	conn.Flush()
}
