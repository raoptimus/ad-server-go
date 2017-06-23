package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
	"tc/data"
)

type (
	Storage struct {
		db  *sqlx.DB
		ads item
	}
	item struct {
		sync.RWMutex
		items map[int]*item
		ctr   *float64
	}
)

func NewStorage() *Storage {
	s := &Storage{
		db: data.DataContext.PgSqlDb,
		ads: item{
			items: make(map[int]*item, 0),
		},
	}
	go s.loadAll()
	return s
}

func (s *Storage) GetCtrByAd(adId int) (ctr float64) {
	s.ads.RLock()
	defer s.ads.RUnlock()
	itm, ok := s.ads.items[adId]

	if ok {
		ctr = *itm.ctr
	}
	return
}

func (s *Storage) loadCtrDailyAds() {
	recovery()
	fmt.Println("Loading CTR from statsByAds")
	rows, err := s.db.Query(`SELECT "AdId", "ShowCount", "ClickCount" FROM tc."statsByAds" WHERE "ShowCount" > 0 AND "ClickCount" > 0 ORDER BY "ForDate" DESC`)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()
	firstAdId := 0
	var adId, showCount, clickCount, newCount, updateCount, deleteCount int

	for rows.Next() {
		err := rows.Scan(&adId, &showCount, &clickCount)
		if err != nil {
			log.Println(err)
			return
		}
		if firstAdId == adId {
			break
		}
		if firstAdId == 0 {
			firstAdId = adId
		}

		var ctr float64
		if showCount > 0 {
			ctr = (float64(clickCount) / float64(showCount)) * float64(100)
		} else {
			ctr = 0.0
		}

		s.ads.RLock()
		itm, ok := s.ads.items[adId]
		s.ads.RUnlock()

		if !ok {
			if ctr == 0.0 {
				continue
			}
			itm = &item{
				ctr: &ctr,
			}
			s.ads.Lock()
			s.ads.items[adId] = itm
			s.ads.Unlock()
			newCount++
			continue
		}
		if ctr == 0.0 {
			s.ads.Lock()
			delete(s.ads.items, adId)
			s.ads.Unlock()
			deleteCount++
			continue
		}
		if *itm.ctr != ctr {
			itm.Lock()
			*itm.ctr = ctr
			itm.Unlock()
			updateCount++
		}
	}

	if rows.Err() != nil {
		log.Println(rows.Err())
	}

	fmt.Sprintf("CTR from statsByAds is loaded, new: %d, update: %d, delete: %d", newCount, updateCount, deleteCount)
}

func (s *Storage) loadAll() {
	go s.loadCtrDailyAds()
}
