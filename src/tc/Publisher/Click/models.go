package main

import (
	"strconv"
	"strings"
	"tc/openrtbex"
)

type (
	AdCode struct {
		Id       int                  `bson:"_id"`
		AdZoneId int                  `bson:"AdZoneId"`
		TypeId   openrtbex.AdCodeType `bson:"TypeId"`
		SiteId   int                  `bson:"SiteId"`

		Site   *Site   `bson:"-"`
		AdZone *AdZone `bson:"-"`
	}
	AdZone struct {
		Id     int                  `bson:"_id"`
		TypeId openrtbex.AdZoneType `bson:"TypeId"`
	}
	Site struct {
		Id           int                `bson:"_id"`
		CategoryId   openrtbex.Category `bson:"CategoryId"`
		IsAllowAdult bool               `bson:"IsAdult"`
		UserId       int                `bson:"UserId"`
	}
)

type (
	Ad struct {
		Id         int    `bson:"_id"`
		CityId     int    `bson:"CityId"`
		Url        string `bson:"URL"`
		CampaignId int    `bson:"CampaignId"`

		Campaign *Campaign `bson:"-"`
	}
	Campaign struct {
		Id               int                    `bson:"_id"`
		Title            string                 `bson:"Title"`
		IsAllowRaw       bool                   `bson:"IsAllowRaw"`
		AddDomainOutput  bool                   `bson:"AddDomainOutput"`
		GettingBadClicks bool                   `bson:"GettingBadClicks"`
		WithGaUtm        bool                   `bson:"WithGaUtm"`
		UserId           int                    `bson:"UserId"`
		PaymentType      openrtbex.PaymentType  `bson:"PaymentType"`
		FriendId         int                    `bson:"FriendId"`
		TypeId           openrtbex.CampaignType `bson:"TypeId"`
	}
)

//todo
func (s *Ad) ClickUrl(hc *HttpContext) string {
	if s.Id == 0 {
		return s.Url
	}

	u := s.Url
	if s.Campaign.WithGaUtm {
		s.addGaUtmc(&u)
	}
	s.replaceMetaTags(&u, hc)
	return u
}

func (s *Campaign) IsFree() bool {
	return s.UserId == 6
}

func (s *Campaign) IsWm(siteUserId int) bool {
	return s.UserId == siteUserId
}

func (ad *Ad) addGaUtmc(url *string) {
	var d string
	if strings.Contains(*url, "?") {
		d = "&"
	} else {
		d = "?"
	}

	pos := strings.Index(*url, "#")

	if pos != -1 {
		*url = (*url)[0:pos] + d + "utm_source=tubecontext&utm_medium=%site%&utm_campaign=%campaign%" + (*url)[pos:]
	} else {
		*url += d + "utm_source=tubecontext&utm_medium=%site%&utm_campaign=%campaign%"
	}
}

func (s *Ad) replaceMetaTags(url *string, hc *HttpContext) {
	*url = strings.Replace(*url, "#campaign#", "%campaign%", -1)
	*url = strings.Replace(*url, "%campaign%", strconv.Itoa(s.Campaign.Id), -1)
	*url = strings.Replace(*url, "#site#", "%site%", -1)
	*url = strings.Replace(*url, "%site%", strconv.Itoa(hc.Data.Site.Id), -1)
	*url = strings.Replace(*url, "#ad#", "%ad%", -1)
	*url = strings.Replace(*url, "%ad%", strconv.Itoa(s.Id), -1)
	*url = strings.Replace(*url, "#device#", "%device%", -1)
	*url = strings.Replace(*url, "%device%", hc.Device.ToString(), -1)
	*url = strings.Replace(*url, "#type#", "%type%", -1)
	*url = strings.Replace(*url, "%type%", hc.Data.AdZone.TypeId.ToString(), -1)
	*url = strings.Replace(*url, "#os#", "%os%", -1)
	*url = strings.Replace(*url, "%os%", hc.Os.ToString(), -1)
	*url = strings.Replace(*url, "#browser#", "%browser%", -1)
	*url = strings.Replace(*url, "%browser%", hc.Browser.ToString(), -1)
	*url = strings.Replace(*url, "#operator#", "%operator%", -1)
	*url = strings.Replace(*url, "%operator%", hc.Operator.ToString(), -1)
}
