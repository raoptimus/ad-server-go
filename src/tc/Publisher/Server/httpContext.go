package main

import (
	"encoding/hex"
	"errors"
	"github.com/bsm/openrtb"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"tc/detect"
	"tc/openrtbex"
	"tc/tchttp"
)

type HttpContext struct {
	*tchttp.HttpContext
	AdCode *AdCode
	AdZone *AdZone
	Site   *Site

	//width ad's image
	AdW int
	//height ad's image
	AdH int
	//limit ads
	AdL     int
	MinCpm  float32
	BlockId int
	Page    string
	Ref     string
	Host    string

	Req     *openrtb.Request
	Win     BidList
	Loss    BidList
	Auction *Auction
	Rc4     *Rc4Bin

	Device   openrtbex.Device
	Os       openrtbex.Os
	Browser  openrtbex.Browser
	Operator openrtbex.Operator

	GeoId   int
	Country *detect.Country
	City    *detect.City
	Region  *detect.Region
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) (h *HttpContext, err error) {
	h = &HttpContext{}
	h.HttpContext, err = tchttp.NewHttpContext(w, r)
	return
}

func (s *HttpContext) LoadData(device *detect.DeviceDetector, geo *detect.GeoDetector,
	operator *detect.OperatorDetector, rc4 *Rc4Bin) error {
	q := s.R.URL.Query()
	x := q.Get("x")

	if x == "" {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("X is empty")
	}

	b, err := hex.DecodeString(x)
	if err != nil {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("X is not encode")
	}

	xDec := string(b)
	codeId := 0

	if strings.Contains(xDec, "AdCodeId") { //old logical
		v, err := url.ParseQuery(xDec)
		if err != nil {
			s.ErrCode = http.StatusMethodNotAllowed
			return err
		}

		codeId, err = strconv.Atoi(v.Get("AdCodeId"))
		if err != nil {
			s.ErrCode = http.StatusMethodNotAllowed
			return err
		}
	} else {
		xp := strings.Split(xDec, "|")
		if len(xp) < 2 {
			s.ErrCode = http.StatusMethodNotAllowed
			return errors.New("X is invalided")
		}

		codeId, err = strconv.Atoi(xp[1])
		if err != nil {
			s.ErrCode = http.StatusMethodNotAllowed
			return errors.New("CodeId is not int")
		}
	}

	s.Site, err = StoreContext.Sites.GetSiteByAdCode(codeId)
	if err != nil {
		return err
	}

	s.AdCode, err = s.Site.GetAdCode(codeId)
	if err != nil {
		return err
	}

	s.AdZone = s.AdCode.AdZone
	if !s.AdCode.Enabled || !s.AdZone.IsActive {
		s.ErrCode = http.StatusForbidden
		return errors.New("AdZone disabled")
	}

	s.AdL, s.AdW, s.AdH, err = s.AdCode.GetParams()
	if err != nil {
		return err
	}

	s.BlockId, _ = strconv.Atoi(q.Get("blockid"))
	s.Page = q.Get("page")
	s.Ref = s.W.Header().Get("Referer")
	if s.Page == "" {
		s.Page = s.Ref
	}

	s.Device, s.Os, s.Browser = device.Detect(s.Ua)
	s.Operator = operator.Detect(s.Ip)
	s.GeoId, s.Country, s.City, s.Region = geo.Detect(s.Ip)
	s.MinCpm = s.AdZone.GetMinCpm(s.Region)
	s.Rc4 = rc4
	s.Host = s.R.Host

	return nil
}
