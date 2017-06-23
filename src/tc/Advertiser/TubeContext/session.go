package main

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Session struct {
		sync.RWMutex
		Id            string              `json:"_id"`
		Expires       time.Time           `json:"expires"`
		Ip            string              `json:"ip"`
		TotalRequests *int32              `json:"totalRequests"`
		CampClicked   *SessionCampClicked `json:"campaignClicked"`
		AdShown       *SessionAdShown     `json:"adShown"`

		adLock       *sessionAdLock `json:"-"`
		releaseCount *int32         `json:"-"`
		releaseTimer *time.Timer    `json:"-"`
	}

	sessionAdLock struct {
		sync.RWMutex
		items map[string]time.Time `json:"items"`
	}

	SessionCampClicked struct {
		sync.RWMutex
		items []int `json:"items"`
	}

	SessionAdShown struct {
		sync.RWMutex
		ShowCount map[string]*int32 `json:"showCount"`
	}
)

const LIFE_AD = time.Second * 10
const LIFE_SESSION = time.Hour * 24
const AD_NEW_PER_REQ = 5

func NewSession(key, addr string) *Session {
	n := int32(0)

	ss := &Session{
		Id:            key,
		Ip:            addr,
		Expires:       time.Now().Add(LIFE_SESSION),
		TotalRequests: &n,
		CampClicked: &SessionCampClicked{
			items: make([]int, 0),
		},
		AdShown: &SessionAdShown{
			ShowCount: make(map[string]*int32),
		},
	}

	ss.reset()
	return ss
}

func (s *Session) reset() {
	rn := int32(0)
	s.releaseCount = &rn
	s.adLock = &sessionAdLock{
		items: make(map[string]time.Time),
	}
}

func (s *Session) getAdShowCount(adId string) int {
	s.AdShown.RLock()
	defer s.AdShown.RUnlock()

	shows, ok := s.AdShown.ShowCount[adId]

	if !ok {
		return 0
	}

	return int(atomic.LoadInt32(shows))
}

func (s *Session) incAdShowCount(adId string) {
	s.AdShown.RLock()
	shows, ok := s.AdShown.ShowCount[adId]
	s.AdShown.RUnlock()

	if !ok {
		i := int32(0)
		shows = &i
		s.AdShown.Lock()
		s.AdShown.ShowCount[adId] = shows
		s.AdShown.Unlock()
	}

	atomic.AddInt32(shows, 1)
	runtime.Gosched()
}

func (s *Session) isNewRequest() bool {
	t := atomic.LoadInt32(s.TotalRequests) + 1
	return t%AD_NEW_PER_REQ == 0
}

func (s *Session) incRequest() {
	atomic.AddInt32(s.TotalRequests, 1)
	runtime.Gosched()
}

func (s *Session) lockAd(adId string) bool {
	s.adLock.RLock()
	expires, ok := s.adLock.items[adId]
	s.adLock.RUnlock()

	if !ok || expires.Before(time.Now()) {
		s.adLock.Lock()
		s.adLock.items[adId] = time.Now().Add(LIFE_AD)
		s.adLock.Unlock()
		return true
	}

	return false
}

func (s *Session) campWasClicked(campId int) bool {
	s.CampClicked.RLock()
	defer s.CampClicked.RUnlock()

	for _, id := range s.CampClicked.items {
		if id == campId {
			return true
		}
	}

	return false
}

func (s *Session) campClick(campId int) {
	s.CampClicked.Lock()
	defer s.CampClicked.Unlock()

	s.CampClicked.items = append(s.CampClicked.items, campId)
}
