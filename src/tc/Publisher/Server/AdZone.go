package main

import (
	"fmt"
	"tc/detect"
	"tc/openrtbex"
)

const (
	DefaultMinCpm float32 = 0.0
)

type AdZone struct {
	Id                        int                  `bson:"_id"`
	SiteId                    int                  `bson:"SiteId"`
	UserId                    int                  `bson:"UserId"`
	IsActive                  bool                 `bson:"IsActive"`
	Player                    openrtbex.PlayerType `bson:"Player"`
	TypeId                    openrtbex.AdZoneType `bson:"TypeId"`
	MinPricePerClickByCountry map[string]float32   `bson:"MinPricePerClickByCountry"`
	Style                     string               `bson:"Style"`

	a *AdStyle
}

func (s AdZone) GetStyle() *AdStyle {
	if s.a != nil {
		return s.a
	}

	s.a = NewAdStyle(s.Style)
	return s.a
}

func (s *AdZone) GetMinCpm(r *detect.Region) float32 {
	id := fmt.Sprintf("%d", r.Id)
	minCpm, ok := s.MinPricePerClickByCountry[id]

	if ok {
		return minCpm
	}

	return DefaultMinCpm
}
