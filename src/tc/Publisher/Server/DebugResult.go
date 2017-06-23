package main

import (
	"encoding/json"
	"fmt"
	"github.com/bsm/openrtb"
	"log"
	"strconv"
	"tc/openrtbex"
	"tc/subscribe"
)

type DebugResult struct {
	hc *HttpContext
}

func NewDebugResult(hc *HttpContext) *DebugResult {
	return &DebugResult{
		hc: hc,
	}
}

func (s *DebugResult) compile() string {
	out := "<html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\" /><style>b {color: white;} </style></head>"
	out += "<body style=\"background-color: black; color: #ccc; font-size: 14px; word-wrap: break-word; width: 98%; font-family: arial;\"><pre>"
	out += s.compileInfo()

	success, fail := s.hc.Auction.GetResults()
	out += s.compileResp(success, "lightgreen")
	out += s.compileResp(fail, "red")

	out += "</pre></body></html>"
	return out
}

func (s *DebugResult) compileInfo() string {
	hc := s.hc
	reqExt := hc.Req.Ext["requestExt"].(openrtbex.RequestExt)
	vals := map[string]string{
		"AdType":       hc.AdCode.TypeId.ToString(),
		"AdZoneType":   hc.AdCode.AdZone.TypeId.ToString(),
		"AdZoneId":     strconv.Itoa(hc.AdCode.AdZone.Id),
		"AdCodeId":     strconv.Itoa(hc.AdCode.Id),
		"SiteId":       strconv.Itoa(hc.Site.Id),
		"SiteHost":     hc.Site.Host,
		"SiteCategory": hc.Site.CategoryId.ToFormatRTB(),
		//todo
		"Ip":        *hc.Req.Device.Ip,
		"GeoId":     strconv.Itoa(reqExt.GeoId),
		"Region":    "",
		"CountryId": strconv.Itoa(reqExt.CountryId),
		"CityId":    strconv.Itoa(reqExt.CityId),
		"Device":    reqExt.Device.ToString(),
		"Os":        reqExt.Os.ToString(),
		"Broser":    reqExt.Browser.ToString(),
		"Operator":  reqExt.Operator.ToString(),
	}
	out := "[(<b style=\"color: white;\">Info</b>)]$ "

	for l, v := range vals {
		out += fmt.Sprintf("<b>%s</b>: %s, ", l, v)
	}

	out += "<br><hr><br>"
	return out
}

func (s *DebugResult) compileResp(respList []*openrtb.Response, color string) string {
	hc := s.hc
	out := ""

	for _, resp := range respList {
		debugExt := resp.Ext["debugExt"].(openrtbex.DebugExt)
		sub := resp.Ext["sub"].(*subscribe.Subscriber)
		out += fmt.Sprintf("[(<b style=\"color: %s;\">%s</b>)]$ <b>Request</b>: %v, <b>Response</b>: %v, <b>Status</b>: %v <b>CustomData</b>: %v<br><br>",
			color,
			sub.Name,
			debugExt.HttpRequestContent,
			debugExt.HttpResponseContent,
			debugExt.HttpStatusCode,
			debugExt.CustomData,
		)

		for _, seatbid := range resp.Seatbid {
			for i, bid := range seatbid.Bid {
				bidStatusExt := bid.Ext["bidStatusExt"].(*openrtbex.BidStatusExt)
				var status string
				var bidColor string

				switch bidStatusExt.Status {
				case openrtbex.ConfirmStatusWin:
					{
						status = "Win"
						bidColor = "green"
						bidExt := bid.Ext["bidExt"].(openrtbex.BidExt)
						log.Println(hc.R.Host, hc.R.URL.String(), hc.R.RequestURI)
						u, err := bidExt.ClickUrl(hc.R, hc.Debug, hc.Ip, hc.Ua, hc.AdCode.Id)
						if err != nil {
							bidExt.Url = err.Error()
						} else {
							bidExt.Url = u.String()
						}
						bid.Ext["bidExt"] = bidExt
						log.Println(u.String())
					}
				case openrtbex.ConfirmStatusLoss:
					{
						status = "Loss"
						bidColor = "lightblue"
					}
				case openrtbex.ConfirmStatusError:
					{
						status = "Error"
						bidColor = "red"
					}
				}

				out += fmt.Sprintf("[%s (<b style=\"color: %s;\">%s<sup>%d</sup></b>)]$ <b>Bid</b>: %v<br><br>",
					status,
					bidColor,
					sub.Name,
					i,
					s.enc(bid))
			}
		}
	}

	return out
}

func (s *DebugResult) enc(item interface{}) string {
	b, err := json.MarshalIndent(&item, "", "    ")

	if err != nil {
		log.Println(err)
		return ""
	}

	return string(b)
}
