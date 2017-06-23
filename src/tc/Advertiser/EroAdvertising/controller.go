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
	"net/rpc"
	"runtime/debug"
	"strconv"
	"strings"
	"tc/openrtbex"
	"time"
)

const API_URL = "..."
const HTTP_GET_ADS_TIMEOUT = 300 * time.Millisecond
const HTTP_CONFIRM_ADS_TIMEOUT = 2000 * time.Millisecond

type Controller struct {
	yaCur *openrtbex.YaCur
}

func NewController() *Controller {
	c := &Controller{
		yaCur: openrtbex.NewYaCur(),
	}

	go c.listenServe()
	return c
}

func (s *Controller) SendRequest(req *openrtb.Request, resp *openrtb.Response) error {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}
	}()
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	debugExt := &openrtbex.DebugExt{}
	isDebug := reqExt.IsDebug
	respId := string(reqExt.CodeId)
	resp.Id = &respId

	if isDebug {
		resp.Ext = openrtb.Extensions{
			"debugExt": debugExt,
		}
	}

	bidRequest := NewBidRequest(req)

	if isDebug {
		b, err := json.Marshal(&bidRequest)

		if err == nil {
			debugExt.HttpRequestContent = string(b)
		}
	}

	transport := &http.Transport{}
	dial := &net.Dialer{Timeout: HTTP_GET_ADS_TIMEOUT}
	transport.Dial = dial.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport

	params := bidRequest.ToUrlValues()
	url := API_URL + "?" + params.Encode()

	if isDebug {
		debugExt.HttpRequestContent = url
	}

	httpReq, err := http.NewRequest("GET", url, nil)
	httpReq.Header.Add("Content-Type", "application/json; charset=utf-8")
	httpReq.Header.Add("Accept", "application/json")

	timer := time.AfterFunc(HTTP_GET_ADS_TIMEOUT, func() {
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

	if isDebug {
		debugExt.HttpStatusCode = httpResp.StatusCode
	}

	if httpResp.StatusCode != 200 {
		transport.CancelRequest(httpReq)
		err := errors.New(fmt.Sprintf("Status: %d, request canceled", httpResp.StatusCode))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	bodyStr := string(body)

	if isDebug {
		debugExt.HttpResponseContent = bodyStr
	}

	if bodyStr == "{\"results\":null}" || bodyStr == "Access denied" { //server can't return 403
		transport.CancelRequest(httpReq)
		err := errors.New(fmt.Sprintf("Status: %d, request canceled", httpResp.StatusCode))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	var bidResponse *BidResponse

	if err := json.Unmarshal(body, &bidResponse); err != nil {

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if len(bidResponse.Results.Ads) <= 0 {
		err := errors.New("Banners not found, server returned: " + string(body))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	cur := s.yaCur.GetCur(openrtbex.CurEuro)
	bidResponse.ToRTBResponse(resp, cur)

	return nil
}

func (s *Controller) SendConfirm(resp *openrtb.Response, success *bool) error {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}
	}()

	*success = false
	errs := make([]string, 0)

	for _, seatbid := range resp.Seatbid {
		for _, bid := range seatbid.Bid {
			status := bid.Ext["bidStatusExt"].(openrtbex.BidStatusExt)

			if status.Status != openrtbex.ConfirmStatusWin {
				continue
			}

			err := s.confirmWin(bid.Nurl)

			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ". "))
	} else {
		*success = true
	}

	return nil
}

func (s *Controller) confirmWin(url *string) error {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
		}
	}()

	transport := &http.Transport{}
	dialer := &net.Dialer{Timeout: HTTP_CONFIRM_ADS_TIMEOUT}
	transport.Dial = dialer.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport

	httpReq, err := http.NewRequest("GET", *url, nil)

	timer := time.AfterFunc(HTTP_CONFIRM_ADS_TIMEOUT, func() {
		transport.CancelRequest(httpReq)
	})
	httpResp, err := client.Do(httpReq)
	timer.Stop()

	if err != nil {
		return err
	}

	if httpResp.StatusCode != 200 {
		transport.CancelRequest(httpReq)
		return errors.New("Confirm return status: " + strconv.Itoa(httpResp.StatusCode))
	}

	defer httpResp.Body.Close()

	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return err
	}

	result := string(body)

	if result != "ok" {
		return errors.New("Confirm return invalid status: " + result)
	}

	return nil
}

func (s *Controller) listenServe() {
	server, err := net.Listen("tcp", "127.0.0.1:8084")

	if err != nil {
		panic(err)
	}

	rpc.Register(s)
	openrtbex.GobRegisterExt()

	for {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				log.Println("PANIC:", r)
				return
			}
		}()

		conn, err := server.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					debug.PrintStack()
					log.Println("PANIC:", r)
					return
				}
			}()
			defer conn.Close()
			conn.SetDeadline(time.Now().Add(600 * time.Millisecond))
			rpc.ServeConn(conn)
		}(conn)
	}
}
