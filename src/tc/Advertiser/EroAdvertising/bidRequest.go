package main

import (
	"fmt"
	"github.com/bsm/openrtb"
	"net/url"
	"tc/openrtbex"
)

const (
	CategoryAdult string = "IAB25-3"
)

type BidRequest struct {
	Ua       *string `json:"ua"`
	Ip       *string `json:"ip"`
	Language *string `json:"in"`
	AdCount  *int    `json:"numads"`
	Ref      *string `json:"doc"`
	SiteId   int     `json:"siteid"`
	SpaceId  *int    `json:"spaceid"`
}

type AdCodeType int

const (
	InVideoPauseRoll      AdCodeType = 0
	Teasers               AdCodeType = 1
	InVideoOverlay        AdCodeType = 2
	InVideoPostRoll       AdCodeType = 4
	InVideoPreRoll        AdCodeType = 5
	InEmbedOverlay        AdCodeType = 6
	InHtml5VideoPauseRoll AdCodeType = 7
	InHtml5VideoOverlay   AdCodeType = 8
	Banners300x250        AdCodeType = 9
	InEmbedPreRoll        AdCodeType = 10
	Popunder              AdCodeType = 11

	MobileBanners300x250 AdCodeType = 20
	MobileBanners300x100 AdCodeType = 21
	MobileBanners300x50  AdCodeType = 22
	MobilePopunder       AdCodeType = 23
)

const (
	bidRequestSiteId = 68259
)

func NewBidRequest(req *openrtb.Request) *BidRequest {
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	bidReq := &BidRequest{}
	bidReq.Ua = req.Device.Ua
	bidReq.Ip = req.Device.Ip
	bidReq.Language = req.Device.Language
	adCount := req.Imp[0].Ext["impExt"].(openrtbex.ImpExt).Count
	bidReq.AdCount = &adCount
	bidReq.Ref = req.Site.Ref
	bidReq.SiteId = bidRequestSiteId
	typeId := reqExt.CodeTypeId

	var spaceId int

	switch typeId {
	case openrtbex.AdCodeTypeInVideoPreRoll, openrtbex.AdCodeTypeInVideoPostRoll, openrtbex.AdCodeTypeMobileBanners300x250, openrtbex.AdCodeTypeBanners300x250, openrtbex.AdCodeTypeInEmbedPreRoll:
		{
			spaceId = 275309
		}

	case openrtbex.AdCodeTypeMobileBanners300x100:
		{
			spaceId = 275310
		}

	case openrtbex.AdCodeTypeMobileBanners300x50:
		{
			spaceId = 275311
		}

	case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInEmbedOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
		{
			spaceId = 275312
		}

	default:
		{
			spaceId = 275309
		}
	}

	bidReq.SpaceId = &spaceId
	return bidReq
}

func (s *BidRequest) ToUrlValues() *url.Values {
	ref := ""
	lang := ""

	if s.Ref != nil {
		ref = *s.Ref
	}

	if s.Language != nil {
		lang = *s.Language
	}

	return &url.Values{
		"ua":      {*s.Ua},
		"ip":      {*s.Ip},
		"in":      {lang},
		"numads":  {fmt.Sprintf("%d", *s.AdCount)},
		"doc":     {ref},
		"siteid":  {fmt.Sprintf("%d", s.SiteId)},
		"spaceid": {fmt.Sprintf("%d", *s.SpaceId)},
	}
}
