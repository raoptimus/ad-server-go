package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
	"tc/data"
)

type (
	SiteStorage struct {
		sync.RWMutex
		items map[int]*Site
		fk    *fkAdCodeSite
	}
	fkAdCodeSite struct {
		sync.RWMutex
		items map[int]int
	}
)

func NewSiteStorage() *SiteStorage {
	store := &SiteStorage{
		items: make(map[int]*Site),
		fk: &fkAdCodeSite{
			items: make(map[int]int),
		},
	}
	store.LoadReloadAll()
	return store
}
func (s *SiteStorage) addFK(codeId int, siteId int) {
	s.fk.RLock()
	_, ok := s.fk.items[codeId]
	s.fk.RUnlock()

	if ok {
		return
	}

	s.fk.Lock()
	s.fk.items[codeId] = siteId
	s.fk.Unlock()
}

func (s *SiteStorage) GetSiteByAdCode(codeId int) (site *Site, err error) {
	defer do(measure("cache-get-site-by-code"))
	s.fk.RLock()
	siteId, ok := s.fk.items[codeId]
	s.fk.RUnlock()

	if ok {
		return s.Get(siteId)
	}

	s.loadAdCode(codeId)

	s.fk.RLock()
	siteId, ok = s.fk.items[codeId]
	s.fk.RUnlock()

	if ok {
		return s.Get(siteId)
	}

	return nil, errors.New(fmt.Sprintf("siteint not found by adCodeint %d", codeId))
}

func (s *SiteStorage) Get(id int) (site *Site, err error) {
	defer do(measure("cache-get-site"))

	s.RLock()
	site, ok := s.items[id]
	s.RUnlock()

	if !ok {
		return nil, errors.New(fmt.Sprintf("Not found %d site", id))
	}

	return site, nil
}

func (s *SiteStorage) GetAll() []*Site {
	copy := make([]*Site, 0)

	for _, site := range s.items {
		copy = append(copy, site)
	}

	return copy
}

func (s *SiteStorage) Set(id int, site *Site) {
	s.RLock()
	site, ok := s.items[id]
	s.RUnlock()

	if !ok {
		s.Lock()
		defer s.Unlock()

		s.items[id] = site
	} else {
		*site = *site
	}
}

func (s *SiteStorage) LoadById(id int) {
	s.load(id)
}

func (s *SiteStorage) load(siteId int) {
	defer do(trace("find-sites"))

	var result []*Site

	where := bson.M{}

	if siteId > 0 {
		where["_id"] = siteId
	}

	where["Deleted"] = false
	where["Approved"] = true
	where["IsVerified"] = true

	query := data.DataContext.Sites.Find(where)
	err := query.All(&result)

	if err != nil {
		log.Println("Preload sites error:", err)
		data.DataContext.Sites.Database.Session.Refresh()
	} else {
		newCount := 0

		s.Lock()
		defer s.Unlock()

		for _, site := range result {
			oldSite, ok := s.items[site.Id]

			if ok {
				site.User = oldSite.User
				site.AdCodes = oldSite.AdCodes
				*oldSite = *site
			} else {
				site.User = StoreContext.Users.Get(site.UserId)
				site.AdCodes = make(map[int]*AdCode)
				s.items[site.Id] = site
				newCount++
			}
		}

		log.Printf("Found %d new sites of %d, \n", newCount, len(result))
	}
}

func (s *SiteStorage) loadAdCode(codeId int) {
	var adCode *AdCode
	err := data.DataContext.AdCodes.Find(bson.M{"_id": codeId}).One(&adCode)

	if err != nil {
		log.Printf("Load adCode %d error %v", codeId, err)
		data.DataContext.AdCodes.Database.Session.Refresh()
		return

	}

	site, err := s.Get(adCode.SiteId)

	if err != nil {
		return
	}

	adCode.Site = site

	site.RLock()
	oldAdCode, ok := site.AdCodes[adCode.Id]
	site.RUnlock()

	if ok {
		adCode.AdZone = oldAdCode.AdZone
		*oldAdCode = *adCode

		s.addFK(adCode.Id, site.Id)

		return
	}

	var adZone *AdZone
	err = data.DataContext.AdZones.Find(bson.M{"_id": adCode.AdZoneId}).One(&adZone)

	if err != nil {
		log.Printf("Load adZone %d error %v", adCode.AdZoneId, err)
		data.DataContext.AdZones.Database.Session.Refresh()
		return
	}

	adCode.AdZone = adZone

	site.Lock()
	site.AdCodes[adCode.Id] = adCode
	site.Unlock()

	s.addFK(adCode.Id, site.Id)
}

func (s *SiteStorage) LoadReloadAll() {
	s.load(0)
}

func (s *SiteStorage) Delete(id int) {
	copy := make(map[int]*Site)

	for id2, site2 := range s.items {
		if id2 == id {
			s.fk.Lock()
			for adCodeId, _ := range site2.AdCodes {
				delete(s.fk.items, adCodeId)
			}
			s.fk.Unlock()

			continue
		}

		copy[id2] = site2
	}

	s.items = copy
}

func (s *SiteStorage) Len() int {
	return len(s.items)
}
