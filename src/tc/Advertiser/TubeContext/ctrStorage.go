package main

import (
	"crypto/md5"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"sync"
	"tc/data"
	"tc/openrtbex"
	"time"
)

const DEFAULT_CTR = 0.01234567

type (
	CtrStorage struct {
		sync.RWMutex
		items    ctrStorageItems
		duration time.Duration
		queue    *ctrStorageQueue
	}
	ctrStorageItem struct {
		expires time.Time
		data    *BidCtr
	}
	ctrStorageItems map[string]*ctrStorageItem
	ctrStorageQueue struct {
		sync.RWMutex
		queue  map[string]bool
		worker bool
	}
)

func NewCtrStorage(duration time.Duration) *CtrStorage {
	c := CtrStorage{
		duration: duration,
		items:    make(ctrStorageItems),
		queue: &ctrStorageQueue{
			queue:  make(map[string]bool),
			worker: false,
		},
	}
	return &c
}
func (c *CtrStorage) PutAll(list map[string]*BidCtr) {
	c.Lock()
	defer c.Unlock()

	for key, data := range list {
		c.items[key] = &ctrStorageItem{
			expires: time.Now().Add(c.duration),
			data:    data,
		}
	}
}

func (c *CtrStorage) Put(key string, data *BidCtr) {
	c.Lock()
	defer c.Unlock()

	c.items[key] = &ctrStorageItem{
		expires: time.Now().Add(c.duration),
		data:    data,
	}
}

func (s *CtrStorage) GetOneOrDefault(ad *Ad, reqExt *openrtbex.RequestExt) *BidCtr {
	cityId := 0

	if ad.CityId > 0 {
		cityId = ad.CityId
	}

	var oldDeviceId int
	var operatorId int

	switch reqExt.Device {
	case openrtbex.DeviceMobile:
		{
			switch reqExt.Os {
			case openrtbex.OsIOs:
				oldDeviceId = 11
			case openrtbex.OsAndroid:
				oldDeviceId = 12
			case openrtbex.OsWindows:
				oldDeviceId = 13
			case openrtbex.OsSymbian:
				oldDeviceId = 14
			case openrtbex.OsBlackBerry:
				oldDeviceId = 15
			default:
				oldDeviceId = 10
			}

			operatorId = reqExt.Operator.ToInt()
		}
	case openrtbex.DeviceTablet:
		{
			switch reqExt.Os {
			case openrtbex.OsIOs:
				oldDeviceId = 1
			case openrtbex.OsAndroid:
				oldDeviceId = 2
			default:
				oldDeviceId = 0
			}

			operatorId = reqExt.Operator.ToInt()
		}
	default:
		{
			oldDeviceId = -1
			operatorId = -1
		}
	}

	k := strconv.Itoa(ad.Id) + "|" +
		strconv.Itoa(reqExt.GeoId) + "|" +
		strconv.Itoa(cityId) + "|" +
		strconv.Itoa(reqExt.CodeTypeId.ToInt()) + "|" +
		strconv.Itoa(oldDeviceId) + "|" +
		strconv.Itoa(operatorId)

	crypt := md5.New()
	crypt.Write([]byte(k))
	id := hex.EncodeToString(crypt.Sum(nil))
	return s.get(id)
}

func (c *CtrStorage) get(id string) *BidCtr {
	defer do(measure("cache-get-BidCtr"))
	c.RLock()
	item, exists := c.items[id]
	c.RUnlock()
	var data *BidCtr

	if exists {
		expired := item.expires.Before(time.Now())
		data = item.data

		if expired {
			item.expires.Add(time.Since(item.expires)).Add(c.duration)
			go c.queuePush(id)
		}
	} else {
		data = c.getDefault(id)
		go c.queuePush(id)
	}

	return data
}

func (c *CtrStorage) getDefault(id string) *BidCtr {
	return &BidCtr{
		Id:          id,
		Ctr:         DEFAULT_CTR,
		IsNew:       true,
		StatBlockId: 0,
	}
}

func (c *CtrStorage) queuePush(key string) {
	c.queue.RLock()

	if c.queue.queue[key] {
		c.queue.RUnlock()
		return
	}

	c.queue.RUnlock()

	c.queue.Lock()
	c.queue.queue[key] = true
	c.queue.Unlock()

	if !c.queue.worker {
		c.queue.worker = true
		go c.loader()
	}
}

func (c *CtrStorage) loader() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("CtrStorage work panic:", r)
			return
		}
	}()

	defer func() {
		c.queue.Lock()
		c.queue.worker = false
		c.queue.Unlock()
	}()

	found := false

	for {
		items := make(ctrStorageItems)
		batch := make(ctrStorageItems)

		//copy
		c.RLock()
		for id, item := range c.items {
			items[id] = item
		}
		c.RUnlock()

		c.queue.RLock()
		for id, _ := range c.queue.queue {
			batch[id] = &ctrStorageItem{
				expires: time.Now().Add(c.duration),
			}

			if len(batch) >= 100 {
				break
			}
		}
		c.queue.RUnlock()

		for id, item := range batch {
			item.data = c.findOneOrDef(id)
			items[id] = item
		}

		c.Lock()
		c.items = items
		c.Unlock()

		c.queue.Lock()

		for id, _ := range batch {
			delete(c.queue.queue, id)
		}

		c.queue.Unlock()

		if len(batch) > 0 {
			found = true
		} else {
			if !found {
				break //finish
			}

			time.Sleep(1000 * time.Millisecond)
			found = false
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (c *CtrStorage) findOneOrDef(id string) *BidCtr {
	var ctr BidCtr

	query := data.DataContext.AdIndex.FindId(id).
		Select(bson.M{"_id": 1, "Ctr": 1, "IsNew": 1, "StatBlockId": 1})
	err := query.One(&ctr)

	if err != nil {
		data.DataContext.AdIndex.Database.Session.Refresh()
		return c.getDefault(id)
	}

	return &ctr
}

func (c *CtrStorage) Delete(key string) {
	defer do(measure("cache-delete"))
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

func (c *CtrStorage) Len() int {
	return len(c.items)
}
