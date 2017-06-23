package main

import (
	"sync"
	"sync/atomic"
	"tc/store"
	"time"
)

type (
	SessionStorage struct {
		sync.RWMutex
		*store.RedisStorage
		fk map[string]*Session
	}
)

const DELAY_DELETE_FK = 3 * time.Second

func NewSessionStorage(dbName, network, address string, life time.Duration) *SessionStorage {
	return &SessionStorage{
		fk:           make(map[string]*Session),
		RedisStorage: store.NewRedisStorage(dbName, network, address, life),
	}
}

func (s *SessionStorage) get(key, addr string) (ss *Session) {
	s.RLock()
	ss, ok := s.fk[key]
	s.RUnlock()

	if ok {
		s.releaseLock(ss)
		return ss
	}

	s.Lock()
	defer s.Unlock()

	ss, ok = s.fk[key]

	if ok {
		s.releaseLock(ss)
		return ss
	}

	ss = &Session{}
	err := s.RedisStorage.Get(key, ss)

	if err != nil {
		ss = NewSession(key, addr)
	} else {
		ss.reset()
	}

	ss.releaseTimer = time.AfterFunc(DELAY_DELETE_FK, func() {
		s.deleteFk(ss)
	})
	s.fk[key] = ss
	s.releaseLock(ss)
	return ss
}

func (s *SessionStorage) releaseLock(ss *Session) {
	ss.releaseTimer.Reset(DELAY_DELETE_FK)
	atomic.AddInt32(ss.releaseCount, 1)
}

func (s *SessionStorage) releaseUnlock(ss *Session) {
	ss.AdShown.RLock()
	modified := ss.AdShown.modified
	ss.AdShown.RUnlock()

	if modified {
		s.RedisStorage.Set(ss.Id, ss)

		ss.AdShown.Lock()
		ss.AdShown.modified = false
		ss.AdShown.Unlock()
	}

	atomic.AddInt32(ss.releaseCount, -1)
}

func (s *SessionStorage) deleteFk(ss *Session) {
	if !atomic.CompareAndSwapInt32(ss.releaseCount, 0, 0) {
		return
	}

	s.Lock()
	defer s.Unlock()
	delete(s.fk, ss.Id)
}
