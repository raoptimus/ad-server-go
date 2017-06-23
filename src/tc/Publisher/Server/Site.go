package main

import (
	"errors"
	"sync"
	"tc/openrtbex"
)

type (
	Site struct {
		sync.RWMutex
		Id              int                `bson:"_id"`
		Approved        bool               `bson:"Approved"`
		CategoryId      openrtbex.Category `bson:"CategoryId"`
		Host            string             `bson:"Host"`
		IsActive        bool               `bson:"IsActive"`
		IsAdult         bool               `bson:"IsAdult"`
		IsPremium       bool               `bson:"IsPremium"`
		IsShowWatermark bool               `bson:"IsShowWatermark"`
		UserId          int                `bson:"UserId"`

		User    *User
		AdCodes map[int]*AdCode
	}
)

func (s *Site) GetAdCode(id int) (adCode *AdCode, err error) {
	s.RLock()
	adCode, ok := s.AdCodes[id]
	s.RUnlock()

	if ok {
		return
	}

	err = errors.New("adCode not found in site")
	return
}

func (s *Site) GetCategoryListRTB() []string {
	list := make([]string, 0)

	if s.CategoryId == openrtbex.CategoryAdult && s.IsAdult {
		list = append(list, openrtbex.CategoryAdult.ToFormatRTB())
	} else {
		list = append(list, openrtbex.CategoryOther.ToFormatRTB())

		if s.IsAdult {
			list = append(list, openrtbex.CategoryAdult.ToFormatRTB())
		}
	}

	if s.CategoryId != openrtbex.CategoryAdult && s.CategoryId != openrtbex.CategoryOther {
		list = append(list, s.CategoryId.ToFormatRTB())
	}

	return list
}
