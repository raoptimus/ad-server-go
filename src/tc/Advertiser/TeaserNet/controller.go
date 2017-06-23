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
	"net/url"
	"runtime/debug"
	"strings"
	"tc/openrtbex"
	"time"
)

const API_URL = "http://..."
const HTTP_GET_ADS_TIMEOUT = 3000 * time.Millisecond
const HTTP_CONFIRM_TIMEOUT = 10000 * time.Millisecond

type Controller struct {
	sessions *SessionStorage
}

func NewController() *Controller {
	c := &Controller{
		sessions: NewSessionStorage("TeaserNet", "unix", "/tmp/redis.sock", LIFE_SESSION),
	}
	go c.listenServe()

	return c
}

func (s *Controller) SendRequest(req *openrtb.Request, resp *openrtb.Response) error {
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

	ss := s.sessions.get(*req.User.Id, *req.Device.Ip)
	defer s.sessions.releaseUnlock(ss) // free session

	bidRequest := NewBidRequest(req, ss)
	b, err := json.Marshal(&bidRequest)

	if err != nil {
		err = errors.New("Request json encode error: " + err.Error())

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if isDebug {
		debugExt.HttpRequestContent = string(b)
	}

	transport := &http.Transport{}
	dial := &net.Dialer{Timeout: HTTP_GET_ADS_TIMEOUT}
	transport.Dial = dial.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport

	data := url.Values{"json": {string(b)}}
	httpReq, err := http.NewRequest("POST", API_URL, strings.NewReader(data.Encode()))
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Connection", "close")

	timer := time.AfterFunc(HTTP_GET_ADS_TIMEOUT, func() {
		transport.CancelRequest(httpReq)
	})
	httpResp, err := client.Do(httpReq)
	timer.Stop()

	if err != nil {
		err = errors.New("Request error: " + err.Error())
		log.Println(err) //todo
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
		if httpResp.StatusCode == 204 {
			log.Println(string(b))
		}
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
		err = errors.New("Read response error: " + err.Error())

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if isDebug {
		debugExt.HttpResponseContent = string(body)
	}

	var bidResponse *BidResponse

	if err := json.Unmarshal(body, &bidResponse); err != nil {
		err = errors.New("Response encode json error: " + err.Error())
		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if bidResponse.Error != "" {
		err = errors.New("TN return error: " + bidResponse.Error)

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	if len(bidResponse.Teasers) <= 0 {
		err := errors.New(fmt.Sprintf("Teasers not found; adCodeId: %v, tnBlockId: %v, tnSiteId: %v\n %v\n%v",
			*bidRequest.MiscId, *bidRequest.BlockId, *bidRequest.SiteId, string(b), string(body)))

		if isDebug {
			debugExt.Error = err.Error()
			return nil
		}
		return err
	}

	bidResponse.ToRTBResponse(bidRequest, resp)

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

	respExt := resp.Ext["responseExt"].(openrtbex.ResponseExt)
	ss := s.sessions.get(*respExt.UserId, *respExt.DeviceIp)
	defer s.sessions.releaseUnlock(ss)

	bidConfirm := NewBidConfirm(resp, ss)
	b, err := json.Marshal(&bidConfirm)

	if err != nil {
		err = errors.New("BidConfirm serialization is failed: " + err.Error())
		return err
	}

	transport := &http.Transport{}
	dialer := &net.Dialer{Timeout: HTTP_CONFIRM_TIMEOUT}
	transport.Dial = dialer.Dial
	transport.MaxIdleConnsPerHost = 1000
	client := &http.Client{}
	client.Transport = transport

	data := url.Values{"json": {string(b)}}
	httpReq, err := http.NewRequest("POST", API_URL, strings.NewReader(data.Encode()))
	httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Add("Content-Length", fmt.Sprintf("%d", len(string(b))))
	httpReq.Header.Add("Accept", "application/json")

	timer := time.AfterFunc(HTTP_CONFIRM_TIMEOUT, func() {
		transport.CancelRequest(httpReq)
	})
	httpResp, err := client.Do(httpReq)
	timer.Stop()

	if err != nil {
		return errors.New("BidConfirm post is failed: " + err.Error())
	}

	if httpResp.StatusCode != 200 {
		transport.CancelRequest(httpReq)
		err = errors.New("BidConfirm request is canceled")
		return err
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		err = errors.New("Read the body after confirm is failed: " + err.Error())
		return err
	}

	var result *openrtbex.ConfirmResult

	if err := json.Unmarshal(body, &result); err != nil {
		err = errors.New("Unserialization confirm result is failed: " + err.Error())
		return err
	}

	*success = result.Result

	if !result.Result {
		err = errors.New("Result confirm not ok, return: " + string(body))
		return err
	}

	return nil
}

func (s *Controller) listenServe() {
	server, err := net.Listen("tcp", "127.0.0.1:8082")

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
