package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
	"tc/data"
	"tc/openrtbex"
	"time"
)

type (
	CampaignStorage struct {
		sync.RWMutex
		items  map[int]*Campaign
		adCamp map[int]int
	}
)

func NewCampaignStorage() *CampaignStorage {
	store := &CampaignStorage{
		items: make(map[int]*Campaign),
	}
	store.LoadReloadAll()
	return store
}

func (self *CampaignStorage) Get(id int) (camp *Campaign, err error) {
	defer do(measure("cache-get-camp"))

	self.RLock()
	camp, ok := self.items[id]
	self.RUnlock()

	if !ok {
		return nil, errors.New(fmt.Sprintf("Not found %v campaign", id))
	}

	return camp, nil
}

func (self *CampaignStorage) GetAll() []*Campaign {
	items := make([]*Campaign, 0)

	for _, camp := range self.items {
		items = append(items, camp)
	}

	return items
}

func (self *CampaignStorage) Set(id int, camp *Campaign) {
	self.RLock()
	camp, ok := self.items[id]
	self.RUnlock()

	if !ok {
		self.Lock()
		defer self.Unlock()

		self.items[id] = camp
	} else {
		*self.items[id] = *camp
	}
}

func (self *CampaignStorage) LoadById(id int) {
	self.load(id, 0)
}

func (self *CampaignStorage) LoadByUserId(userId int) {
	self.load(0, userId)
}

func (self *CampaignStorage) load(campId int, userId int) {
	defer do(trace("find-campaigns"))

	var result []*Campaign

	where := bson.M{"IsActive": true, "PreDeleted": false, "FriendId": 0}

	if campId > 0 {
		where["_id"] = campId
	}

	if userId > 0 {
		where["UserId"] = userId
	}

	query := data.DataContext.Campaigns.Find(where)
	err := query.All(&result)

	if err != nil {
		log.Println("preload campaigns error:", err)
		data.DataContext.Campaigns.Database.Session.Refresh()
	} else {
		newCount := 0

		self.Lock()
		defer self.Unlock()

		for _, camp := range result {
			oldCamp, ok := self.items[camp.Id]

			if ok {
				camp.User = oldCamp.User
			} else {
				camp.User = StoreContext.Users.Get(camp.UserId)
			}

			//todo
			//			if !camp.IsDebit() {
			//				continue
			//			}

			if ok {
				camp.Ads = oldCamp.Ads
				*oldCamp = *camp

				if camp.ApplyUrlForAds {
					for _, ad := range camp.Ads {
						if ad.URL == camp.URL {
							continue
						}
						ad.URL = camp.URL
					}
				}
			} else {
				camp.Ads = make(map[int]*Ad)
				self.loadAds(camp)
				self.items[camp.Id] = camp
				newCount++
			}
		}

		log.Printf("Found %d new campaigns of %d, \n", newCount, len(result))
	}
}

func (self *CampaignStorage) loadAds(camp *Campaign) {
	if camp.TypeId == openrtbex.CampaignTypePopunder {
		ad := &Ad{
			Campaign:    camp,
			CampaignId:  camp.Id,
			UserId:      camp.UserId,
			IsWebmaster: camp.IsWebmaster,
		}

		ads := make(map[int]*Ad, 0)
		ads[int(0)] = ad
		camp.Ads = ads

		return
	}

	copy := make(map[int]*Ad, 0)

	for id, ad := range camp.Ads {
		copy[id] = ad
	}

	where := bson.M{"CampaignId": camp.Id, "IsActive": true, "PreDeleted": false, "Approved": true}
	query := data.DataContext.Ads.Find(where)
	iter := query.Iter()

	ad := &Ad{}

	for iter.Next(ad) {
		ad.Campaign = camp
		ad.CtrStorage = NewCtrStorage(20 * time.Minute)
		_, ok := copy[ad.Id]

		if !ok {
			copy[ad.Id] = ad
		} else {
			*copy[ad.Id] = *ad
		}

		ad = &Ad{}
	}

	if err := iter.Err(); err != nil {
		log.Panicln(err)
	}

	camp.Ads = copy
}

func (self *CampaignStorage) LoadReloadAll() {
	self.load(0, 0)
}

func (self *CampaignStorage) UnlimitAll() {
	self.RLock()
	defer self.RUnlock()

	for _, camp := range self.items {
		camp.IsLimited = false
	}
}

func (self *CampaignStorage) Limit(id int) {
	camp, err := self.Get(id)

	if err != nil {
		return
	}

	camp.IsLimited = true
}

func (self *CampaignStorage) Delete(id int) {
	camp, err := self.Get(id)

	if err == nil {
		camp.Deleted = true //for other refs

		copy := make(map[int]*Campaign)

		for id2, camp2 := range self.items {
			if id2 == id {
				continue
			}

			copy[id2] = camp2
		}

		self.items = copy
	}
}

func (self *CampaignStorage) GetAd(id int) (ad *Ad, err error) {
	err = data.DataContext.Ads.Find(bson.M{"_id": id}).One(&ad)

	if err != nil {
		return nil, err
	}

	camp, err := self.Get(ad.CampaignId)

	if err != nil {
		return nil, err
	}

	ad, ok := camp.Ads[id]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Ad %d not found in campaign %d", camp.Id, id))
	}

	return ad, nil
}

func (self *CampaignStorage) LoadReloadAd(id int) {
	var ad *Ad
	where := bson.M{"_id": id, "IsActive": true, "PreDeleted": false, "Approved": true}
	err := data.DataContext.Ads.Find(where).One(&ad)

	if err != nil {
		log.Println("preload ad error:", err)
		return
	}

	camp, err := self.Get(ad.CampaignId)

	if err != nil {
		log.Println("preload ad error:", err)
		return
	}

	copy := make(map[int]*Ad, 0)

	for oid, oad := range camp.Ads {
		copy[oid] = oad
	}

	ad.Campaign = camp
	oad, ok := copy[ad.Id]

	if !ok {
		ad.CtrStorage = NewCtrStorage(20 * time.Minute)
		copy[ad.Id] = ad
	} else {
		ad.CtrStorage = oad.CtrStorage
		*copy[ad.Id] = *ad
	}

	camp.Ads = copy
}

func (self *CampaignStorage) DeleteAd(id int) {
	var ad *Ad
	where := bson.M{"_id": id}
	err := data.DataContext.Ads.Find(where).One(&ad)

	if err != nil {
		log.Println("delete ad error:", err)
		return
	}

	camp, err := self.Get(ad.CampaignId)

	if err != nil {
		log.Println("delete ad error:", err)
		return
	}

	adInCamp, ok := camp.Ads[id]

	if ok {
		adInCamp.Deleted = true //for other refs

		copy := make(map[int]*Ad, 0)

		for _, oad := range camp.Ads {
			if id == oad.Id {
				continue
			}

			copy[oad.Id] = oad
		}

		camp.Ads = copy
	}
}

func (self *CampaignStorage) Len() int {
	return len(self.items)
}
