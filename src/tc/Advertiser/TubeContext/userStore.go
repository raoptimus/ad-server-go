package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
	"tc/data"
	"time"
)

type (
	UserStorage struct {
		sync.RWMutex
		items map[int]*User
	}
)

func NewUserStorage() *UserStorage {
	store := &UserStorage{
		items: make(map[int]*User),
	}
	store.update()
	return store
}

func (self *UserStorage) Updater() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("UserStorage work panic:", r)
		}
	}()

	for {
		time.Sleep(1 * time.Minute)
		self.update()
	}
}

func (self *UserStorage) update() {
	defer do(trace("find-users"))

	var result []*User
	query := data.DataContext.Users.Find(nil).Select(bson.M{"Balance": 1})
	err := query.All(&result)

	if err != nil {
		log.Println("preload users error:", err)
		data.DataContext.Users.Database.Session.Refresh()
	} else {
		newCount := 0

		self.Lock()
		defer self.Unlock()

		for _, user := range result {
			_, ok := self.items[user.Id]

			if !ok {
				self.items[user.Id] = user
				newCount++
			} else {
				*self.items[user.Id] = *user
			}
		}

		log.Printf("Found %d new users of %d\n", newCount, len(result))
	}
}

func (self *UserStorage) Get(id int) *User {
	defer do(measure("cache-get-User"))

	self.RLock()
	user, ok := self.items[id]
	self.RUnlock()

	if !ok {
		user = NewDefaultUser(id)

		self.Lock()
		self.items[id] = user
		self.Unlock()
	}

	return user
}

func (self *UserStorage) Len() int {
	return len(self.items)
}
