package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"tc/tc"
	"time"
)

const API_URL = "http://dsp.kadam.ru/bid/9"
const CAMPAIGN_ID = 1
const USER_ID = 1
const PRIORITY = 3

var httpGetBidTimeout = time.Duration(300 * time.Millisecond)
var httpConfirmBidTimeout = time.Duration(1000 * time.Millisecond)

type Controller int

func (s *Controller) SendRequest(req *openrtb.Request, resp *openrtb.Response) error {
	reqExt := req.Ext["requestExt"].(tc.RequestExt)
	debugExt := &tc.DebugExt{}
	isDebug := reqExt.IsDebug
	respId := string(reqExt.CodeId)
	resp.Id = &respId

	if isDebug {
		resp.Ext = openrtb.Extensions{
			"debugExt": debugExt,
		}
	}

	imp := req.Imp[0]
	impExt := imp.Ext["impExt"].(tc.ImpExt)
	width := *imp.Banner.W
	height := *imp.Banner.H

	if width > 600 && height == 80 { //overlay as teaser, no banner
		width = 100
		height = 100
	}

	typeId := BlockTypeTeaser

	switch reqExt.ZoneType {
	case tc.AdZoneTypeInVideo: //1
		{ //InVideo
			if width%height == 0 { //sq
				width = 200
				height = 200
				typeId = BlockTypeTeaser
			} else {
				typeId = BlockTypeBanner
			}
		}
	case tc.AdZoneTypeBanners, tc.AdZoneTypeMobileBanners:
		{ //Banners
			typeId = BlockTypeBanner
		}
	case tc.AdZoneTypeTeasers:
		{ //Teasers
			typeId = BlockTypeTeaser
		}
	case tc.AdZoneTypePopunder, tc.AdZoneTypeMobilePopunder:
		{ //Popunders
			typeId = BlockTypeClickUnder
		}
	}

	req.Ext = openrtb.Extensions{ //rewrite our ext
		"blocks": []*Block{
			&Block{
				Id:     &reqExt.CodeId,
				Type:   &typeId,
				Limit:  &impExt.Count,
				Width:  &width,
				Height: &height,
			},
		},
	}

	req.Imp = []openrtb.Impression{ //rewrite our imp
		openrtb.Impression{
			Id: imp.Id,
			Banner: &openrtb.Banner{
				W: &width,
				H: &height,
			},
		},
	}

	req.Site.Ext = openrtb.Extensions{} //clear our ext

	b, err := json.Marshal(&req)

	if err != nil {
		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if reqExt.IsDebug {
		debugExt.HttpRequestContent = string(b)
	}

	transport := &http.Transport{}
	dial := &net.Dialer{Timeout: httpGetBidTimeout}
	transport.Dial = dial.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport

	httpReq, err := http.NewRequest("POST", API_URL, strings.NewReader(string(b)))
	httpReq.Header.Add("x-openrtb-version", "2.2")
	httpReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Add("Content-Length", fmt.Sprintf("%d", len(string(b))))
	httpReq.Header.Add("Accept", "application/json")
	//	req.Header.Add("Connection", "close");

	timer := time.AfterFunc(httpGetBidTimeout, func() {
		transport.CancelRequest(httpReq)
	})
	httpResp, err := client.Do(httpReq)
	timer.Stop()

	if err != nil {
		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if reqExt.IsDebug {
		debugExt.HttpStatusCode = httpResp.StatusCode
	}

	if httpResp.StatusCode != 200 {
		transport.CancelRequest(httpReq)
		err = errors.New(fmt.Sprintf("Http status: %d, request canceled", httpResp.StatusCode))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	defer httpReq.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if isDebug {
		debugExt.HttpResponseContent = string(body)
	}

	if err := json.Unmarshal(body, &resp); err != nil { //todo записать string(body) в ответ
		err := errors.New(fmt.Sprintf("Error: %s\nResponse: '%s'\nStatus: %d", err, string(body), httpResp.StatusCode))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if len(resp.Seatbid) <= 0 {
		return nil
	}

	bidList := resp.Seatbid[0].Bid

	if len(bidList) <= 0 {
		return nil
	}

	bidFirst := bidList[0]
	ads := bidFirst.Ext["ads"].([]interface{})

	copyBidList := make([]openrtb.Bid, 0)

	for _, adI := range ads {
		ad := adI.(map[string]interface{})
		bid := openrtb.Bid{}
		bid.Impid = bidFirst.Impid
		//bid.Id = bidFirst.Id
		bid.Nurl = bidFirst.Nurl
		*bid.Nurl = strings.Replace(*bid.Nurl, "${AUCTION_ID}", *resp.Id, 1)
		cpm := float32(ad["cpm"].(float64))
		bid.Price = &cpm

		bid.Id = resp.Id
		adId := fmt.Sprintf("%d", int(ad["id"].(float64)))
		bid.Adid = &adId
		bid.Id = &adId

		bid.Ext = openrtb.Extensions{
			"bidExt": tc.BidExt{
				PaymentType: tc.PaymentTypeCPM,
				Title:       ad["title"].(string),
				Image:       ad["image"].(string),
				Url:         ad["url"].(string),
				Width:       ad["width"].(int),
				Height:      ad["height"].(int),
				Price:       cpm,
				Priority:    PRIORITY,
				Type:        tc.BidTypePay,
				CampaignId:  CAMPAIGN_ID,
				UserId:      USER_ID,
			},
		}

		copyBidList = append(copyBidList, bid)
	}

	resp.Seatbid[0].Bid = copyBidList

	return nil
}

func (s *Controller) SendConfirm(resp *openrtb.Response, ret *int) error {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}
	}()

	bidIdList := make([]string, 0)

	for _, seatbid := range resp.Seatbid {
		for _, bid := range seatbid.Bid {
			status := bid.Ext["bidStatusExt"].(tc.BidStatusExt)

			if status.Status != tc.ConfirmStatusWin {
				continue
			}

			bidIdList = append(bidIdList, *bid.Id)
		}
	}

	firstBid := resp.Seatbid[0].Bid[0]
	firstStatus := firstBid.Ext["bidStatusExt"].(tc.BidStatusExt)
	winPrice := firstStatus.WinningPrice
	bidIds := strings.Join(bidIdList, ",")
	url := strings.Replace(*firstBid.Nurl, "${AUCTION_AD_ID}", bidIds, 1)
	url = strings.Replace(url, "${AUCTION_PRICE}", fmt.Sprintf("%v", winPrice), 1)
	url = strings.Replace(url, "${AUCTION_ID}", *resp.Id, 1)
	err := s.confirm(url)

	return err
}

func (s *Controller) confirm(url string) error {
	transport := &http.Transport{}
	dialer := &net.Dialer{Timeout: httpConfirmBidTimeout}
	transport.Dial = dialer.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport
	req, err := http.NewRequest("GET", url, nil)

	timer := time.AfterFunc(httpConfirmBidTimeout, func() {
		transport.CancelRequest(req)
	})
	resp, err := client.Do(req)
	timer.Stop()

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		transport.CancelRequest(req)
		return errors.New(fmt.Sprintf("Confirm request is canceled, status: %d", resp.StatusCode))
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var result tc.ConfirmResult

	if err := json.Unmarshal(body, &result); err != nil {
		return errors.New(fmt.Sprintf("Confirm error: %s, response: '%s', status: %d", err, string(body), resp.StatusCode))
	}

	if !result.Result {
		return errors.New("Result confirm not ok")
	}

	return nil
}
