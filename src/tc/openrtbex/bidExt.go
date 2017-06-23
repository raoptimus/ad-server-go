package openrtbex

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"hash/crc32"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const SHOW_URL = "http://......com"

type BidExt struct {
	PaymentType  PaymentType  `json:"ptype"`
	Title        string       `json:"title,omitempty"`
	Width        int          `json:"w,omitempty"`
	Height       int          `json:"h,omitempty"`
	Image        string       `json:"image"`
	Url          string       `json:"url"`
	Price        float32      `json:"price"`
	Priority     int          `json:"priority"`
	Ctr          float32      `json:"ctr"`
	Type         BidType      `json:"type"`
	CampaignId   int          `json:"campId"`
	StatBlockId  int          `json:"statId"`
	AdId         int          `json:"adId"`
	UserId       int          `json:"userId"`
	BrokerId     int          `json:"brokerId"`
	CampaignType CampaignType `json:"ctype"`

	//confirm

}

func (s *BidExt) ClickUrl(req *http.Request, debug, ip, ua string, adCodeId int) (u *url.URL, err error) {
	//gen hash
	creation := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	ttl := time.Duration(time.Second * 20).String()
	h := crc32.NewIEEE()
	_, err = h.Write([]byte(ip + ua + creation + ttl))
	if err != nil {
		return
	}
	//M-X
	mv := make(url.Values)
	if s.AdId == 0 {
		mv.Add("ca", strconv.Itoa(s.CampaignId))
		mv.Add("u", s.Url)
	} else {
		mv.Add("a", strconv.Itoa(s.AdId))
	}

	mv.Add("c", strconv.Itoa(adCodeId))
	mv.Add("h", strconv.Itoa(int(h.Sum32())))
	mv.Add("p", strconv.FormatFloat(float64(s.Price), 'f', 2, 32))
	mv.Add("ct", creation)
	mv.Add("tt", ttl)

	var b bytes.Buffer
	w := base64.NewEncoder(base64.URLEncoding, &b)
	fw, err := flate.NewWriter(w, flate.BestCompression)
	if err != nil {
		return
	}
	_, err = fw.Write([]byte(mv.Encode()))
	if err != nil {
		return
	}
	fw.Flush()
	fw.Close()
	w.Close()

	host := req.Host

	if debug != "" {
		hosts := strings.Split(req.Host, ":")
		host = hosts[0] + ":10082"
	} else if host == "xf.tubecontext.com" {
		host = "xg.tubecontext.com"
	} else {
		//todo
	}

	u = &url.URL{
		Host:     host,
		Path:     "/z/",
		RawQuery: "m=" + b.String(),
	}

	if debug != "" {
		u.RawQuery += "&debug=" + debug
	}

	return
}
