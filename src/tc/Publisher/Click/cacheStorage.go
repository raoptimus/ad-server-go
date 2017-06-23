package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
	"tc/data"
)

type (
	CacheStorage struct {
		adCodeList *cacheStorageAdCode
		siteList   *cacheStorageSite
		adZoneList *cacheStorageAdZone
		adList     *cacheStorageAd
		campList   *cacheStorageCampaign
	}
	cacheStorageAdCode struct {
		sync.RWMutex
		items map[int]*AdCode
	}
	cacheStorageSite struct {
		sync.RWMutex
		items map[int]*Site
	}
	cacheStorageAdZone struct {
		sync.RWMutex
		items map[int]*AdZone
	}
	cacheStorageAd struct {
		sync.RWMutex
		items map[int]*Ad
	}
	cacheStorageCampaign struct {
		sync.RWMutex
		items map[int]*Campaign
	}
)

type (
	Data struct {
		Ad       *Ad
		Campaign *Campaign
		AdZone   *AdZone
		AdCode   *AdCode
		Site     *Site
	}
)

func NewCacheStorage() *CacheStorage {
	return &CacheStorage{
		adCodeList: &cacheStorageAdCode{
			items: make(map[int]*AdCode),
		},
		siteList: &cacheStorageSite{
			items: make(map[int]*Site),
		},
		adZoneList: &cacheStorageAdZone{
			items: make(map[int]*AdZone),
		},
		adList: &cacheStorageAd{
			items: make(map[int]*Ad),
		},
		campList: &cacheStorageCampaign{
			items: make(map[int]*Campaign),
		},
	}
}

func (s *CacheStorage) GetData(adCodeId, adId, campaignId int) (d *Data, err error) {
	var (
		adCode *AdCode
		adZone *AdZone
		site   *Site
		ad     *Ad
		camp   *Campaign
	)
	adCode, err = s.getAdCode(adCodeId)
	if err != nil {
		return
	}

	if adCode.AdZone == nil {
		if adCode.AdZoneId == 0 {
			b, _ := json.Marshal(adCode)
			log.Println(string(b))
			err = errors.New("AdZoneId is empty")
			return
		}
		adCode.AdZone, err = s.getAdZone(adCode.AdZoneId)
		if err != nil {
			return
		}
	}
	adZone = adCode.AdZone

	if adCode.Site == nil {
		adCode.Site, err = s.getSite(adCode.SiteId)
		if err != nil {
			return
		}
	}
	site = adCode.Site

	if adId > 0 {
		ad, err = s.getAd(adId)
		if err != nil {
			return
		}
		campaignId = ad.CampaignId
	} else {
		ad = &Ad{
			Id:     0,
			CityId: 0,
		}
	}

	if ad.Campaign == nil {
		if campaignId == 0 {
			err = errors.New("CampaignId is empty")
			return
		}
		ad.Campaign, err = s.getCampaign(campaignId)
		if err != nil {
			return
		}
	}

	camp = ad.Campaign

	d = &Data{
		AdCode:   adCode,
		AdZone:   adZone,
		Site:     site,
		Ad:       ad,
		Campaign: camp,
	}

	return
}

func (s *CacheStorage) getAdCode(id int) (adCode *AdCode, err error) {
	s.adCodeList.RLock()
	adCode, ok := s.adCodeList.items[id]
	s.adCodeList.RUnlock()

	if ok {
		return
	}

	err = data.DataContext.AdCodes.Find(bson.M{"_id": id}).One(&adCode)
	if err != nil {
		err = errors.New(fmt.Sprintf("Load adCode (%d) from mdb is failed", id))
		return
	}

	s.adCodeList.Lock()
	s.adCodeList.items[id] = adCode
	s.adCodeList.Unlock()

	return adCode, nil
}

func (s *CacheStorage) getSite(id int) (site *Site, err error) {
	s.siteList.RLock()
	site, ok := s.siteList.items[id]
	s.siteList.RUnlock()

	if ok {
		return site, nil
	}

	err = data.DataContext.Sites.Find(bson.M{"_id": id}).One(&site)

	if err != nil {
		err = errors.New(fmt.Sprintf("Load site (%d) from mdb is failed", id))
		return
	}

	s.siteList.Lock()
	s.siteList.items[id] = site
	s.siteList.Unlock()

	return site, nil
}

func (s *CacheStorage) getAdZone(id int) (adZone *AdZone, err error) {
	s.adZoneList.RLock()
	adZone, ok := s.adZoneList.items[id]
	s.adZoneList.RUnlock()

	if ok {
		return
	}

	err = data.DataContext.AdZones.Find(bson.M{"_id": id}).One(&adZone)

	if err != nil {
		err = errors.New(fmt.Sprintf("Load adZone (%d) from mdb is failed", id))
		return
	}

	s.adZoneList.Lock()
	s.adZoneList.items[id] = adZone
	s.adZoneList.Unlock()

	return adZone, nil
}

func (s *CacheStorage) getAd(id int) (ad *Ad, err error) {
	s.adList.RLock()
	ad, ok := s.adList.items[id]
	s.adList.RUnlock()

	if ok {
		return
	}

	err = data.DataContext.Ads.Find(bson.M{"_id": id}).One(&ad)

	if err != nil {
		err = errors.New(fmt.Sprintf("Load ad (%d) from mdb is failed", id))
		return
	}

	s.adList.Lock()
	s.adList.items[id] = ad
	s.adList.Unlock()
	return
}

func (s *CacheStorage) getCampaign(id int) (camp *Campaign, err error) {
	s.campList.RLock()
	camp, ok := s.campList.items[id]
	s.campList.RUnlock()

	if ok {
		return
	}

	err = data.DataContext.Campaigns.Find(bson.M{"_id": id}).One(&camp)
	if err != nil {
		err = errors.New(fmt.Sprintf("Load campaign (%d) from mdb is failed", id))
		return
	}

	s.campList.Lock()
	s.campList.items[id] = camp
	s.campList.Unlock()
	return
}
