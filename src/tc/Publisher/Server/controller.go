package main

import (
	"errors"
	"github.com/bsm/openrtb"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"tc/detect"
	"tc/openrtbex"
	"tc/subscribe"
)

const RC4_KEY = "RA$7xPZR"

type (
	Controller struct {
		Subscriber *subscribe.Subscribers
		Device     *detect.DeviceDetector
		Geo        *detect.GeoDetector
		Operator   *detect.OperatorDetector
		Rc4bin     *Rc4Bin
	}
)

func NewController() *Controller {
	openrtbex.GobRegisterExt()
	subs := subscribe.NewSubscribers()
	subs.Subscribe(&subscribe.Subscriber{
		Name:        "TeaserNet",
		Network:     "tcp",
		Address:     "127.0.0.1:8082",
		LossConfirm: false,
		BrokerId:    1,
	})
	//	subs.subscriber.Subscribe(&subscribe.Subscriber{
	//		Name:        "RTBWin",
	//		Network:     "tcp",
	//		Address:     "127.0.0.1:9082",
	//		LossConfirm: false,
	//    BrokerId:    7,
	//	})
	//	subs.subscriber.Subscribe(&subscribe.Subscriber{
	//		Name:        "TubeContext",
	//		Network:     "tcp",
	//		Address:     "127.0.0.1:8081",
	//		LossConfirm: true,
	//		BrokerId:    0,
	//	})
	subs.Subscribe(&subscribe.Subscriber{
		Name:        "EroAdvertising",
		Network:     "tcp",
		Address:     "127.0.0.1:8084",
		LossConfirm: true,
		BrokerId:    5,
	})

	c := &Controller{
		Subscriber: subs,
	}
	c.Rc4bin, _ = NewRc4Bin(RC4_KEY)
	c.Device = detect.NewDeviceDetector()
	c.Geo = detect.NewGeoDetector()
	c.Operator = detect.NewOperatorDetector()

	go c.listen()
	return c
}

func (s *Controller) serve(w http.ResponseWriter, r *http.Request) {
	defer do(measure("serve"))
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}
	}()

	//haProxy life check
	if r.Method == "HEAD" {
		return
	}

	hc, err := NewHttpContext(w, r)
	if err != nil {
		hc.WriteError(err, "")
		return
	}
	err = hc.LoadData(s.Device, s.Geo, s.Operator, s.Rc4bin)
	if err != nil {
		hc.WriteError(err, "")
		return
	}

	at := 1
	siteId := strconv.Itoa(hc.Site.Id)
	reqId := strconv.Itoa(hc.AdCode.Id)

	req := &openrtb.Request{
		Id: &reqId,
		At: &at, // 1 - аукцион по первой цене
		Site: &openrtb.Site{
			Id:     &siteId,
			Domain: &hc.Site.Host,
			Cat:    hc.Site.GetCategoryListRTB(),
			Page:   &hc.Page,
			Ref:    &hc.Ref,
			Ext: openrtb.Extensions{
				"siteExt": openrtbex.SiteExt{
					Id:           hc.Site.Id,
					UserId:       hc.Site.UserId,
					CategoryId:   hc.Site.CategoryId,
					IsAllowAdult: hc.Site.IsAdult,
					IsPremium:    hc.Site.IsPremium,
				},
			},
		},
		User: &openrtb.User{
			Id: &hc.SessionId,
		},
		Device: &openrtb.Device{
			Ua:         &hc.Ua,
			Ip:         &hc.Ip,
			Language:   &hc.Lang,
			Devicetype: hc.Device.ToOpenRTB(),
		},
		Imp: []openrtb.Impression{
			openrtb.Impression{
				Id:       &reqId,
				Bidfloor: &hc.MinCpm,
				Banner: &openrtb.Banner{
					W: &hc.AdW,
					H: &hc.AdH,
				},
				Ext: openrtb.Extensions{
					"impExt": openrtbex.ImpExt{
						Count: hc.AdL,
					},
				},
			},
		},
		Ext: openrtb.Extensions{
			"requestExt": openrtbex.RequestExt{
				Device:     hc.Device,
				Os:         hc.Os,
				Browser:    hc.Browser,
				CodeId:     hc.AdCode.Id,
				ZoneId:     hc.AdZone.Id,
				CodeTypeId: hc.AdCode.TypeId,
				PlayerType: hc.AdZone.Player,
				ZoneType:   hc.AdZone.TypeId,
				GeoId:      hc.GeoId,
				CountryId:  hc.Country.Id,
				CityId:     hc.City.Id,
				IsDebug:    hc.IsDebug,
				Operator:   hc.Operator,
			},
		},
	}

	a := NewAuction(s.Subscriber)
	win, _ := a.Play(req, hc.AdL, at)
	a.Confirm(!hc.IsDebug, req)

	if !hc.IsDebug && len(win) <= 0 {
		hc.ErrCode = http.StatusNoContent
		hc.WriteError(errors.New("Not content; Win="+strconv.Itoa(len(win))), "")
		return
	}

	hc.Req = req
	hc.Auction = a
	hc.Win = win

	NewFeedResult(hc).write()
}

func (m *Controller) faviconServe(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", 204)
}

func (m *Controller) crossDomainServe(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/xml")
	io.WriteString(w, `<cross-domain-policy><allow-access-from domain="*"/></cross-domain-policy>`)
}

func (m *Controller) listen() {
	http.HandleFunc("/crossdomain.xml", m.crossDomainServe)
	http.HandleFunc("/favicon.ico", m.faviconServe)
	http.HandleFunc("/", m.serve)
	log.Fatal(http.ListenAndServe(":10081", nil))
}
