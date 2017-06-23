package main

import (
	"github.com/bsm/openrtb"
	"log"
	"strconv"
	"sync"
	"tc/openrtbex"
	"time"
)

type (
	AdFilter struct {
		storage *adFilterStorage
	}

	adFilterStorage struct {
		sync.RWMutex
		cache *Cache
	}
)

func NewAdFilter() *AdFilter {
	return &AdFilter{
		storage: &adFilterStorage{
			cache: NewCache(time.Duration(2 * time.Minute)),
		},
	}
}

func (s *AdFilter) getId(reqExt *openrtbex.RequestExt, siteExt *openrtbex.SiteExt) string {
	return strconv.Itoa(siteExt.Id) + "|" +
		strconv.Itoa(reqExt.GeoId) + "|" +
		strconv.Itoa(reqExt.CityId) + "|" +
		strconv.Itoa(reqExt.CodeTypeId.ToInt()) + "|" +
		strconv.Itoa(reqExt.Device.ToInt()) + "|" +
		strconv.Itoa(reqExt.Operator.ToInt()) + "|" +
		strconv.Itoa(reqExt.Os.ToInt()) + "|" +
		strconv.Itoa(reqExt.Browser.ToInt())
}

func (s *AdFilter) Filter(req *openrtb.Request) []*Ad {
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	siteExt := req.Site.Ext["siteExt"].(openrtbex.SiteExt)
	filterId := s.getId(&reqExt, &siteExt)

	adList, exists, expired := s.storage.Get(filterId)

	if exists {
		if expired {
			go s.cuCache(filterId, req, &reqExt)
		}

		return adList
	} else {
		return s.cuCache(filterId, req, &reqExt)
	}
}

func (s *AdFilter) cuCache(filterId string, req *openrtb.Request, reqExt *openrtbex.RequestExt) []*Ad {
	defer do(measure("filter-adFilter"))
	isDebug := reqExt.IsDebug
	campList := StoreContext.Campaigns.GetAll()
	adList := make([]*Ad, 0)

	for _, camp := range campList {
		err := camp.Filter(req)

		if err != nil {
			if isDebug {
				log.Println(err)
			}
			continue
		}

		if camp.TypeId == openrtbex.CampaignTypePopunder {
			ad := &Ad{
				Campaign: camp,
			}

			adList = append(adList, ad)
			continue
		}

		for _, ad := range camp.Ads {
			err := ad.Filter(req)

			if err != nil {
				if isDebug {
					log.Println(err)
				}
				continue
			}

			adList = append(adList, ad)
		}
	}

	//todo aiList.PreCache()
	s.storage.Set(filterId, adList)
	//	_, preAdList := aiList.PreCache()
	//	self.storage.Set(*uk, preAdList)
	//
	//	list := aiFilteredList.Sort(req.IsNew).Uniq(req.AdLimit)

	//	return &list

	return adList
}

func (self *adFilterStorage) Get(uk string) (data []*Ad, exists bool, expired bool) {
	d, exists, expired := self.cache.Get(uk, "adFilter")

	if exists {
		data = d.([]*Ad)
		return data, exists, expired
	}

	return nil, exists, expired
}

func (self *adFilterStorage) Set(uk string, data []*Ad) {
	self.cache.Put(uk, data)
}

func (self *AdFilter) Len() int {
	return self.storage.cache.Len()
}
