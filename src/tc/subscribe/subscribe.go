package subscribe

import (
	"errors"
	"sync"
	"tc/openrtbex"
)

//todo func check, stop if not unavailable
type Subscriber struct {
	Name string
	//"tcp" or "http"
	Network string
	//"127.0.0.1:8080" or "http://domain/..."
	Address string
	//the requirement to send the confirmation of loss
	LossConfirm bool
	AllowTypes  []openrtbex.AdCodeTypeNew
	//FriendId
	BrokerId int
}

type Subscribers struct {
	sync.RWMutex
	items map[string]*Subscriber
}

func NewSubscribers() *Subscribers {
	return &Subscribers{
		items: make(map[string]*Subscriber),
	}
}

func (s *Subscribers) Valid(sub *Subscriber) error {
	if sub.Name == "" {
		return errors.New("Name is empty")
	}
	if sub.Network == "" {
		return errors.New("Network is empty")
	}
	if sub.Address == "" {
		return errors.New("Network is empty")
	}
	if sub.BrokerId < 0 {
		return errors.New("BrokerId is empty")
	}

	s.RLock()
	defer s.RUnlock()

	for name, item := range s.items {
		if item.BrokerId == sub.BrokerId {
			return errors.New("Sub by BrokerId already exists")
		}
		if name == sub.Name {
			return errors.New("Sub by Name already exists")
		}
		if item.Network+item.Address == sub.Network+sub.Address {
			return errors.New("Sub by Network already exists")
		}
	}

	return nil
}

func (s *Subscribers) FindByBroker(brokerId int) *Subscriber {
	s.RLock()
	defer s.RUnlock()

	for _, item := range s.items {
		if item.BrokerId == brokerId {
			return item
		}
	}

	return nil
}

func (s *Subscribers) Subscribe(sub *Subscriber) {
	s.Lock()
	s.items[sub.Name] = sub
	s.Unlock()
}

func (s *Subscribers) UnSubscribe(name string) {
	s.RLock()
	_, ok := s.items[name]
	s.RUnlock()

	if ok {
		return
	}

	s.Lock()
	delete(s.items, name)
	s.Unlock()
}

func (s *Subscribers) GetAll() []*Subscriber {
	s.RLock()
	defer s.RUnlock()

	list := make([]*Subscriber, 0)
	for _, sub := range s.items {
		list = append(list, sub)
	}
	return list
}
