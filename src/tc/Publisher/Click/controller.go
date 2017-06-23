package main

import (
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"runtime/debug"
	"strconv"
	"sync"
	"tc/detect"
	"tc/openrtbex"
	"time"
)

const DEF_REDIRECT_URL_ERROR = "https://google.com"
const DEF_REDIRECT_URL_ADULT = "http://....tv"
const DEF_REDIRECT_URL_NO_ADULT = "http://kinorama.tv"
const SERVER_CONFIRM_NETWORK = "tcp"
const SERVER_CONFIRM_ADDR = ":8888"
const SERVER_CONFIRM_METHOD = "Controller.ClickConfirm"

type Controller struct {
	sync.RWMutex
	device   *detect.DeviceDetector
	geo      *detect.GeoDetector
	operator *detect.OperatorDetector
	confirm  *rpc.Client
}

func NewController() *Controller {
	c := &Controller{}
	c.device = detect.NewDeviceDetector()
	c.geo = detect.NewGeoDetector()
	c.operator = detect.NewOperatorDetector()
	openrtbex.GobRegisterExt()
	c.reconnectConfirm(false)

	log.Println("Server ready")

	go c.listen()
	return c
}

func (s *Controller) listen() {
	http.HandleFunc("/", s.serve)
	log.Fatal(http.ListenAndServe(":10082", nil))
}

func (s *Controller) serve(w http.ResponseWriter, r *http.Request) {
	redirectUrl := DEF_REDIRECT_URL_ERROR
	defer func() {
		if rex := recover(); rex != nil {
			debug.PrintStack()
			log.Println("PANIC:", rex)
			http.Redirect(w, r, redirectUrl, 302)
			return
		}
	}()

	hc, err := NewHttpContext(w, r)
	if err != nil {
		log.Println(err)
		hc.WriteError(err, redirectUrl)
		return
	}

	err = hc.LoadData(s.device, s.geo, s.operator)
	if err != nil {
		log.Println(err)
		hc.WriteError(err, redirectUrl)
		return
	}

	var bad error
	switch {
	case hc.Hash != hc.X.Hash:
		{
			bad = errors.New("Hash is incorrect")
		}
	case hc.X.Creation.Add(hc.X.Ttl).Before(time.Now().UTC()):
		{
			bad = errors.New("Click is expired")
		}
	default:
		bad = nil
	}
	if bad == nil {
		bad = hc.ValidLang()
		if bad == nil {
			bad = hc.ValidUa()
		}
	}

	redirectUrl = hc.Data.Ad.ClickUrl(hc)
	if bad != nil && !hc.Data.Campaign.GettingBadClicks {
		if hc.Data.Site.CategoryId == openrtbex.CategoryAdult {
			redirectUrl = DEF_REDIRECT_URL_ADULT
		} else {
			redirectUrl = DEF_REDIRECT_URL_NO_ADULT
		}
	}

	if r.Method == "GET" {
		w.Header().Add("Content-Type", "text/html;charset=utf-8")
		io.WriteString(w, "<html><header></header><body>")
		io.WriteString(w, fmt.Sprintf("<form method=\"post\" name=\"frm\" action=\"%s\">", "."))
		io.WriteString(w, fmt.Sprintf("<input type=\"hidden\" name=\"m\" value=\"%s\" />", hc.X.Value))
		if hc.IsDebug {
			io.WriteString(w, fmt.Sprintf("<input type=\"hidden\" name=\"debug\" value=\"%s\" />", hc.Debug))
		}
		io.WriteString(w, "<script>document.forms[\"frm\"].submit();</script>")
		io.WriteString(w, fmt.Sprintf("<noscript><meta http-equiv=\"refresh\" content=\"1;url=%s\" /></noscript>", redirectUrl))
		io.WriteString(w, "</body></html>")
		return
	}

	bidType := openrtbex.BidTypePay
	if hc.Data.Campaign.IsFree() {
		bidType = openrtbex.BidTypeFree
	} else if hc.Data.Campaign.IsWm(hc.Data.Site.Id) {
		bidType = openrtbex.BidTypeWm
	}

	badErr := ""
	status := openrtbex.ConfirmStatusWin
	if bad != nil {
		badErr = bad.Error()
		status = openrtbex.ConfirmStatusError
		if !hc.IsDebug {
			bad = nil
			badErr = ""
		}
	}

	adId := strconv.Itoa(hc.Data.Ad.Id)
	bid := openrtb.Bid{
		Id: &adId,
		Ext: openrtb.Extensions{
			"bidExt": openrtbex.BidExt{
				PaymentType:  hc.Data.Campaign.PaymentType,
				Type:         bidType,
				Price:        hc.X.Price,
				CampaignId:   hc.Data.Campaign.Id,
				AdId:         hc.Data.Ad.Id,
				UserId:       hc.Data.Campaign.UserId,
				BrokerId:     hc.Data.Campaign.FriendId,
				CampaignType: hc.Data.Campaign.TypeId,
			},
			"bidStatusExt": openrtbex.BidStatusExt{
				Status: status,
				Error:  badErr,
			},
		},
	}

	respId := strconv.Itoa(hc.Data.AdCode.Id)
	resp := openrtb.Response{
		Id: &respId,
		Seatbid: []openrtb.Seatbid{
			openrtb.Seatbid{
				Bid: []openrtb.Bid{bid},
			},
		},
		Ext: openrtb.Extensions{
			"responseExt": openrtbex.ResponseExt{
				UserId:      &hc.SessionId,
				DeviceIp:    &hc.Ip,
				ConfirmType: openrtbex.ConfirmTypeClick,
			},
		},
	}
	log.Println(resp)
	go s.sendConfirm(&resp)

	hc.RedirectJs(redirectUrl)
}

func (s *Controller) reconnectConfirm(forced bool) bool {
	if !forced {
		s.RLock()
		if s.confirm != nil {
			s.RUnlock()
			return true
		}
		s.RUnlock()
	}

	s.Lock()
	defer s.Unlock()

	client, err := rpc.Dial(SERVER_CONFIRM_NETWORK, SERVER_CONFIRM_ADDR)
	if err != nil {
		log.Println("Dial to server confirm error:", err)
		return false
	}
	s.confirm = client
	return true
}

func (s *Controller) sendConfirm(resp *openrtb.Response) {
	if !s.reconnectConfirm(false) {
		//todo two attempt
		return
	}

	err := s.confirm.Call(SERVER_CONFIRM_METHOD, resp, nil)
	if err != nil && err == rpc.ErrShutdown {
		if !s.reconnectConfirm(true) {
			return
		}
		err = s.confirm.Call(SERVER_CONFIRM_METHOD, resp, nil)
	}

	if err != nil {
		log.Println("Call.", SERVER_CONFIRM_METHOD, err)
	}
}
