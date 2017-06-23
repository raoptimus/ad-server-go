package main

import (
	"errors"
	"fmt"
	"tc/openrtbex"
)

type AdCode struct {
	Id              int                  `bson:"_id"`
	AdZoneId        int                  `bson:"AdZoneId"`
	TypeId          openrtbex.AdCodeType `bson:"TypeId"`
	SiteId          int                  `bson:"SiteId"`
	Enabled         bool                 `bson:"Enabled"`
	IsShowWatermark bool                 `bson:"IsShowWatermark"`

	Site   *Site
	AdZone *AdZone
}

func (s *AdCode) GetParams() (limit, width, height int, err error) {
	switch s.TypeId {
	case openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInHtml5VideoPauseRoll:
		{
			return 2, 250, 250, nil
		}
	case openrtbex.AdCodeTypeTeasers, openrtbex.AdCodeTypeBanners300x250:
		{
			style := s.AdZone.GetStyle()

			if s.TypeId == openrtbex.AdCodeTypeBanners300x250 {
				return style.BlockCount, 300, 250, nil
			} else {
				return style.BlockCount, 250, 250, nil
			}
		}
	case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInEmbedOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
		{
			return 1, 960, 80, nil
		}
	case openrtbex.AdCodeTypeInVideoPostRoll, openrtbex.AdCodeTypeInVideoPreRoll, openrtbex.AdCodeTypeInEmbedPreRoll:
		{
			return 1, 300, 250, nil
		}
	case openrtbex.AdCodeTypeMobileBanners300x250:
		{
			return 1, 300, 250, nil
		}
	case openrtbex.AdCodeTypeMobileBanners300x100:
		{
			return 1, 300, 100, nil
		}
	case openrtbex.AdCodeTypeMobileBanners300x50:
		{
			return 1, 300, 50, nil
		}
	case openrtbex.AdCodeTypePopunder, openrtbex.AdCodeTypeMobilePopunder:
		{
			return 1, 0, 0, nil
		}
	default:
		return 0, 0, 0, errors.New(fmt.Sprintf("Not found typeId %d", s.TypeId))
	}
}

func (s *AdCode) IsCrypt() bool {
	switch s.TypeId {
	case openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInVideoPreRoll, openrtbex.AdCodeTypeInVideoPostRoll:
		{
			return true
		}
	}

	return false
}
