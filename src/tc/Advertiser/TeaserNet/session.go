package main

import (
	"strconv"
	"sync"
	"time"
)

type (
	Session struct {
		Id      string          `json:"_id"`
		Expires time.Time       `json:"expires"`
		Ip      string          `json:"ip"`
		AdShown *SessionAdShown `json:"adShown"`

		releaseCount *int32      `json:"-"`
		releaseTimer *time.Timer `json:"-"`
	}

	SessionAdShown struct {
		sync.RWMutex
		Items []int `json:"adIds"`

		modified bool `json:"-"`
	}
)

const LIFE_AD = time.Second * 10
const LIFE_SESSION = time.Hour * 24
const AD_NEW_PER_REQ = 5

func NewSession(key, addr string) *Session {
	ss := &Session{
		Id:      key,
		Ip:      addr,
		Expires: time.Now().Add(LIFE_SESSION),
		AdShown: &SessionAdShown{
			Items: make([]int, 0),
		},
	}

	ss.reset()
	return ss
}

func (s *Session) reset() {
	rn := int32(0)
	s.releaseCount = &rn
	s.AdShown.modified = false
}

func (s *Session) getAdShownList() []int {
	s.AdShown.RLock()
	defer s.AdShown.RUnlock()

	copy := s.AdShown.Items
	return copy
}

func (s *Session) setAdAsShown(adId *string) {
	id, err := strconv.Atoi(*adId)

	if err != nil {
		return
	}

	s.AdShown.RLock()
	copy := s.AdShown.Items
	s.AdShown.RUnlock()

	for _, adId := range copy {
		if adId == id {
			return
		}
	}

	s.AdShown.Lock()
	defer s.AdShown.Unlock()

	s.AdShown.Items = append(s.AdShown.Items, id)
	s.AdShown.modified = true
}
