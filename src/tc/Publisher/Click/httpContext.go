package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"errors"
	"hash/crc32"
	"net/http"
	"net/url"
	"strconv"
	"tc/detect"
	"tc/openrtbex"
	"tc/tchttp"
	"time"
)

type HttpContext struct {
	*tchttp.HttpContext
	X struct {
		Value    string
		Url      string
		Price    float32
		AdCodeId int
		AdId     int
		Creation time.Time
		Ttl      time.Duration
		CampId   int
		Hash     string
	}
	Hash string
	Data *Data

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
	if err != nil {
		return
	}
	err = h.setX()
	return
}

func (s *HttpContext) LoadData(device *detect.DeviceDetector, geo *detect.GeoDetector, operator *detect.OperatorDetector) error {
	s.Device, s.Os, s.Browser = device.Detect(s.Ua)
	s.Operator = operator.Detect(s.Ip)
	s.GeoId, s.Country, s.City, s.Region = geo.Detect(s.Ip)
	d, err := StoreContext.CacheStorage.GetData(s.X.AdCodeId, s.X.AdId, s.X.CampId)
	if err != nil {
		return err
	}
	if d.Ad.Id == 0 {
		d.Ad.Url = s.X.Url
	}
	s.Data = d
	return nil
}

func (s *HttpContext) setX() error {
	r := s.R
	r.ParseForm()
	xVal := r.Form.Get("m")
	if xVal == "" {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("X is empty")
	}
	b, err := base64.URLEncoding.DecodeString(xVal)
	if err != nil {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("X is not decode: " + err.Error())
	}
	rf := flate.NewReader(bytes.NewBuffer(b))
	var buf bytes.Buffer
	_, err = buf.ReadFrom(rf)
	rf.Close()
	if err != nil {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("X is not uncompress 2: " + err.Error() + buf.String())
	}
	xDec := buf.String()
	v, err := url.ParseQuery(xDec)
	if err != nil {
		s.ErrCode = http.StatusMethodNotAllowed
		return err
	}
	adCodeId, _ := strconv.Atoi(v.Get("c"))
	adId, _ := strconv.Atoi(v.Get("a"))
	campId, _ := strconv.Atoi(v.Get("ca"))
	sCreation := v.Get("ct")
	ttl := v.Get("tt")
	url := v.Get("u")
	price, err := strconv.ParseFloat(v.Get("p"), 64)
	//current hash
	h := crc32.NewIEEE()
	_, err = h.Write([]byte(s.Ip + s.Ua + sCreation + ttl))
	if err != nil {
		s.ErrCode = http.StatusMethodNotAllowed
		return err
	}
	hash := strconv.Itoa(int(h.Sum32()))

	//validate
	if (adId == 0 && campId == 0) || adCodeId == 0 || hash == "" || (url == "" && adId == 0) {
		s.ErrCode = http.StatusMethodNotAllowed
		return errors.New("Data from X is incorrect")
	}

	s.X.Hash = hash
	s.X.Value = xVal
	s.X.Url = url
	s.X.Price = float32(price)
	s.X.AdCodeId = adCodeId
	s.X.AdId = adId
	s.X.CampId = campId

	//hash from X
	s.Hash = v.Get("h")
	if s.Hash == s.X.Hash {
		s.X.Ttl, _ = time.ParseDuration(ttl)
		tm, _ := strconv.Atoi(sCreation)
		s.X.Creation = time.Unix(int64(tm), 0).UTC()
	}

	return nil
}
