package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"log"
	"reflect"
	"strconv"
	"tc/detect"
	"tc/openrtbex"
	"time"
)

const MinBalance = 20
const AllAllowed = -1
const MeasureCount = 1

const (
	WhiteList = 1
	BlackList = 2
)

type ArrayInt []interface{}

func (s ArrayInt) In(id int) bool {
	var err error

	for _, v := range s {
		d := -1
		switch v.(type) {
		case string:
			d, err = strconv.Atoi(v.(string))

			if err != nil {
				log.Printf("Cannot convert %v to int", v)
				return false
			}
		case int:
			d = v.(int)
		case int64:
			d = int(v.(int64))
		case float64:
			d = int(v.(float64))
		default:
			log.Printf("Unknown type for int: %v %v", v, reflect.TypeOf(v))
			return false
		}

		if d == id {
			return true
		}
	}

	return false
}

type (
	Campaign struct {
		Id                   int                    `bson:"_id"`
		Title                string                 `bson:"Title"`
		WithGaUtm            bool                   `bson:"WithGaUtm"`
		AddDomainOutput      bool                   `bson:"AddDomainOutput"`
		URL                  string                 `bson:"URL"`
		UserId               int                    `bson:"UserId"`
		FriendId             int                    `bson:"FriendId"`
		IsAdult              bool                   `bson:"IsAdult"`
		IsAllowRaw           bool                   `bson:"IsRaw"`
		IsMobile             bool                   `bson:"IsMobile"`
		CategoryId           openrtbex.Category     `bson:"CategoryId"`
		PaymentType          openrtbex.PaymentType  `bson:"PaymentType"`
		HourSince            int                    `bson:"HourSince"`
		HourTill             int                    `bson:"HourTill"`
		ForAllWeekDays       bool                   `bson:"ForAllWeekDays"`
		DaysOfWeek           ArrayInt               `bson:"DaysOfWeek"`
		GeoAccesses          map[string]*GeoAccess  `bson:"GeoAccesses"`
		GettingBadClicks     bool                   `bson:"GettingBadClicks"`
		ApplyUrlForAds       bool                   `bson:"ApplyUrlForAds"`
		AllowTypes           ArrayInt               `bson:"AllowTypes"`
		AllowDevices         ArrayInt               `bson:"AllowDevices"`
		AllowBrowsers        ArrayInt               `bson:"AllowBrowsers"`
		AllowOs              ArrayInt               `bson:"AllowOs"`
		AllowMobileOperators ArrayInt               `bson:"AllowMobileOperators"`
		IsWebmaster          bool                   `bson:"IsWebmaster"`
		TypeId               openrtbex.CampaignType `bson:"TypeId"`
		IsActive             bool                   `bson:"IsActive"`
		UseList              int                    `bson:"UseList"`
		SiteFilterList       SiteFilterList         `bson:"SiteFilterList"`
		Percent              int                    `bson:"Percent"`
		IsLimited            bool                   `bson:"IsLimited"`
		Deleted              bool

		User *User       `json:"-"`
		Ads  map[int]*Ad `json:"-"`
	}
	CampaignList []*Campaign

	SiteFilter struct {
		Id int `bson:"_id"`
		// Host    string `bson:"Host"`
		IsAllow bool `bson:"IsAllow"`
	}

	GeoAccess struct {
		IsAllow bool        `bson:"IsAllow"`
		Price   json.Number `bson:"CostPerClick"`
	}

	SiteFilterList []*SiteFilter
)

func (c *Campaign) AdLen() int {
	return len(c.Ads)
}

func (sf SiteFilterList) Find(id int) *SiteFilter {
	for i := 0; i < len(sf); i++ {
		if sf[i].Id == id {
			return sf[i]
		}
	}

	return nil
}

func (c *Campaign) geoKey(countryId int, cityId int) string {
	return fmt.Sprintf("%d|%d|%d", c.Id, countryId, cityId)
}

func (c *Campaign) WeekdayAllowed() bool {
	if c.ForAllWeekDays || len(c.DaysOfWeek) == 7 {
		return true
	}

	weekday := int(time.Now().Weekday())

	if weekday == 0 {
		weekday = 7
	}

	return c.DaysOfWeek.In(weekday)
}

func (c *Campaign) HourAllowed() bool {
	hour := time.Now().Hour()
	since := c.HourSince
	till := c.HourTill

	if since == 0 && till == 0 {
		return true
	}

	if till > since {
		return hour >= since && hour < till
	} else {
		return (hour >= since && hour < 23) || (hour >= 0 && hour < till)
	}
}

func (c *Campaign) GeoIdAllowed(id int) bool {
	if id < 0 || id >= len(detect.GeoGroups) {
		log.Printf("Cannot find geoid %d (max id is %d)", id, len(detect.GeoGroups))
		return false
	}

	geo := detect.GeoGroups[id]
	key := c.geoKey(geo.CountryId, geo.CityId)
	access, ok := c.GeoAccesses[key]

	if !ok {
		return false
	}

	return access.IsAllow
}

func (c *Campaign) GeoPrice(geoId int) float32 {
	geo := detect.GeoGroups[geoId]
	key := c.geoKey(geo.CountryId, geo.CityId)
	// access is always defined
	// because we filter campaigns which have no correct access in GeoIdAllowed
	access, ok := c.GeoAccesses[key]

	if !ok {
		return float32(0)
	}

	price, err := access.Price.Float64()

	if err != nil {
		log.Println(err)
	}

	if price == 0.0 {
		log.Println("Error: price is zero", "key", key, "geo", access)
	}

	return float32(price)
}

func (c *Campaign) IsBase() bool {
	return c.UserId == 6 //user basic
}

func (c *Campaign) SiteAllowed(id int) bool {
	if id == AllAllowed {
		return true
	}
	switch c.UseList {
	case WhiteList:
		{
			f := c.SiteFilterList.Find(id)
			return f != nil && f.IsAllow
		}
	case BlackList:
		{
			f := c.SiteFilterList.Find(id)
			return f == nil || f.IsAllow
		}
	}

	return true
}

func (c *Campaign) DeviceAllowed(d openrtbex.Device) bool {
	if d == AllAllowed {
		return true
	}
	if d == openrtbex.DeviceTablet && c.AllowDevices.In(openrtbex.DeviceDesktop.ToInt()) {
		//можно показать дестопную рекламу на планшете
		return true
	}

	return c.AllowDevices.In(d.ToInt())
}

func (c *Campaign) OsAllowed(o openrtbex.Os) bool {
	if o == AllAllowed {
		return true
	}
	return c.AllowOs.In(o.ToInt())
}

func (c *Campaign) BrowserAllowed(o openrtbex.Browser) bool {
	if o == AllAllowed {
		return true
	}
	return c.AllowBrowsers.In(o.ToInt())
}

func (c *Campaign) TypeAllowed(id openrtbex.AdCodeType) bool {
	return c.AllowTypes.In(id.ToInt())
}

func (c *Campaign) OperatorAllowed(op openrtbex.Operator) bool {
	if op == AllAllowed {
		return true
	}
	return c.AllowMobileOperators.In(op.ToInt())
}

func (c *Campaign) CategoryAllowed(siteCategoryId openrtbex.Category, siteIsAdult bool) bool {

	if c.CategoryId == siteCategoryId {
		return true
	}

	if c.CategoryId != openrtbex.CategoryAdult && siteCategoryId != openrtbex.CategoryAdult {
		return true
	}

	if c.CategoryId == openrtbex.CategoryAdult && siteIsAdult {
		return true
	}

	if siteCategoryId == openrtbex.CategoryAdult && c.IsAdult {
		return true
	}

	return false
}

func (c *Campaign) IsDebit() bool {
	return c.IsWebmaster || c.IsBase() || c.User.Balance >= MinBalance
}

func (c *Campaign) Filter(req *openrtb.Request) error {
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	siteExt := req.Site.Ext["siteExt"].(openrtbex.SiteExt)
	reason := ""

	switch {
	case c.Deleted:
		reason = "it was deleted"
	case c.IsLimited:
		reason = "Limited"
	case c.TypeId == openrtbex.CampaignTypePopunder && c.URL == "":
		reason = "Url is empty for popunder campaign"
	case c.IsWebmaster && int(c.UserId) != siteExt.UserId:
		reason = fmt.Sprintf("it's a webmaster %d campaign, but site's userId is %d", c.UserId, siteExt.UserId)
	case !c.CategoryAllowed(siteExt.CategoryId, siteExt.IsAllowAdult):
		reason = fmt.Sprintf("\n\t\tCampaign|Site\n"+
			"CategoryId:\t%d\t|%d\n"+
			"IsAdult:\t%t\t|%t\n"+
			"IsAllowAdult:\t-\t|%t",
			c.CategoryId, siteExt.CategoryId,
			c.IsAdult, c.IsAdult,
			siteExt.IsAllowAdult)
	case !c.TypeAllowed(reqExt.CodeTypeId):
		reason = fmt.Sprintf("it allows only %v types, but request's typeId is %d", c.AllowTypes, reqExt.CodeTypeId)
	case !c.DeviceAllowed(reqExt.Device):
		reason = fmt.Sprintf("it's DeviceAllowed is %v, but request's device is %d", c.AllowDevices, reqExt.Device)
	case !c.OsAllowed(reqExt.Os):
		reason = fmt.Sprintf("it's AllowOs is %v, but request's Os is %d", c.AllowOs, reqExt.Os)
	case !c.BrowserAllowed(reqExt.Browser):
		reason = fmt.Sprintf("it's AllowBrowsers is %v, but request's Browser is %d", c.AllowBrowsers, reqExt.Browser)
	case !c.OperatorAllowed(reqExt.Operator):
		reason = fmt.Sprintf("it's AllowMobileOperators is %v, but request's Operator is %d", c.AllowMobileOperators, reqExt.Operator)
	case !c.HourAllowed():
		reason = fmt.Sprintf("HourAllowed: %d >= %d < %d", c.HourSince, time.Now().Hour(), c.HourTill)
	case !c.WeekdayAllowed():
		reason = "WeekdayAllowed"
	case !c.SiteAllowed(siteExt.Id):
		reason = fmt.Sprintf("Site %d filtered", siteExt.Id)
	case !c.GeoIdAllowed(reqExt.GeoId):
		reason = "GeoIdAllowed"
	case !c.IsDebit():
		reason = fmt.Sprintf("User has no money %v < %v", c.User.Balance, MinBalance)
		//this use in adIndex.Filter this not for campaign filter
		//	case !c.IsAllowRaw && where.Session().WasClicked(c.Id):
		//		reason = "Already clicked"
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Campaign %d is denied because %+v", c.Id, reason))
}
