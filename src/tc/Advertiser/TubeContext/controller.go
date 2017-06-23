package main

import (
	"encoding/json"
	"github.com/bsm/openrtb"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"runtime"
	"runtime/debug"
	"tc/openrtbex"
	"tc/subscribe"
	"time"
)

const SERVER_SUBSCRIBE_NETWORK = "tcp"
const SERVER_SUBSCRIBE_ADDR = ":8888"
const SERVER_SUBSCRIBE_METHOD = "ApiController.Subscribe"
const SERVER_UNSUBSCRIBE_METHOD = "ApiController.UnSubscribe"

type Controller struct {
	journal  *JournalUpdater
	adFilter *AdFilter
}

func NewController() *Controller {
	c := &Controller{
		journal:  NewJournalUpdater(),
		adFilter: NewAdFilter(),
	}

	go c.listenServe()
	go c.listenHttp()
	go c.subscribe()
	return c
}

func (s *Controller) SendRequest(req *openrtb.Request, resp *openrtb.Response) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC:", r)
			log.Println(runtime.Caller(4))
		}
	}()
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	debugExt := &openrtbex.DebugExt{}
	isDebug := reqExt.IsDebug

	if isDebug {
		resp.Ext = openrtb.Extensions{
			"debugExt": debugExt,
		}
	}

	respId := string(reqExt.CodeId)
	resp.Id = &respId

	var ss *Session
	if reqExt.Debug.DisableSession {
		ss = NewSession("", "")
	} else {
		ss = StoreContext.Sessions.get(*req.User.Id, *req.Device.Ip)
		defer StoreContext.Sessions.releaseUnlock(ss) // free session
	}

	adList := s.adFilter.Filter(req)
	limit := req.Imp[0].Ext["impExt"].(openrtbex.ImpExt).Count
	bidList := NewBidList(adList, req, ss, limit, reqExt.Debug.DisableRotation).ToRtbBidList()

	resp.Seatbid = []openrtb.Seatbid{
		openrtb.Seatbid{
			Bid: bidList,
		},
	}

	if isDebug {
		b, _ := json.Marshal(&req)
		debugExt.HttpRequestContent = string(b)
		b, _ = json.Marshal(&resp)
		debugExt.HttpResponseContent = string(b)
	}

	return nil
}

func (s *Controller) SendConfirm(resp *openrtb.Response, success *bool) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC:", r)
			return
		}
	}()
	respExt := resp.Ext["responseExt"].(openrtbex.ResponseExt)
	ss := StoreContext.Sessions.get(*respExt.UserId, *respExt.DeviceIp)
	defer StoreContext.Sessions.releaseUnlock(ss)

	for _, seatbid := range resp.Seatbid {
		for _, bid := range seatbid.Bid {
			bidStatusExt := bid.Ext["bidStatusExt"].(openrtbex.BidStatusExt)

			if bidStatusExt.Status != openrtbex.ConfirmStatusWin {
				continue
			}

			ss.incAdShowCount(*bid.Id)
		}
	}

	*success = true
	return nil
}

func (s *Controller) subscribe() {
	client, err := rpc.Dial(SERVER_SUBSCRIBE_NETWORK, SERVER_SUBSCRIBE_ADDR)
	if err != nil {
		log.Println("Dial to server subscribe error:", err)
		return
	}
	sub := subscribe.Subscriber{
		Name:        "TubeContext",
		Network:     "tcp",
		Address:     "127.0.0.1:8081",
		LossConfirm: true,
		BrokerId:    0,
	}
	var result bool
	err = client.Call(SERVER_SUBSCRIBE_METHOD, sub, &result)
	if err != nil {
		log.Println("Call subscribe error:", err)
		return
	}

	log.Println("Server is subscribed")
}

func (s *Controller) unSubscribe() {
	client, err := rpc.Dial(SERVER_SUBSCRIBE_NETWORK, SERVER_SUBSCRIBE_ADDR)
	if err != nil {
		log.Println("Dial to server subscribe error:", err)
		return
	}
	defer client.Close()
	name := "TubeContext"
	var result bool
	err = client.Call(SERVER_UNSUBSCRIBE_METHOD, name, &result)
	if err != nil {
		log.Println("Call subscribe error:", err)
		return
	}

	log.Println("Server is unSubscribed")
}

func (s *Controller) listenServe() {
	server, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	rpc.Register(s)
	openrtbex.GobRegisterExt()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Panic(err)
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

func (s *Controller) listenHttp() {
	server := &http.Server{
		Addr:           ":8841",
		Handler:        http.HandlerFunc(s.handleHttp),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *Controller) handleHttp(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query().Get("request")
	if request == "" {
		http.Error(w, "no request", http.StatusBadRequest)
		return
	}

	adapter := struct {
		Site struct {
			Ext struct {
				SiteExt openrtbex.SiteExt
			}
		}
		Imp []struct {
			Ext struct {
				ImpExt openrtbex.ImpExt
			}
		}
		Ext struct {
			RequestExt openrtbex.RequestExt
		}
	}{}
	err := json.Unmarshal([]byte(request), &adapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	reqId := "42" //ignored
	req := &openrtb.Request{
		Id: &reqId,
		Site: &openrtb.Site{Ext: openrtb.Extensions{
			"siteExt": adapter.Site.Ext.SiteExt,
		}},
		Imp: []openrtb.Impression{
			openrtb.Impression{
				Ext: openrtb.Extensions{
					"impExt": adapter.Imp[0].Ext.ImpExt,
				},
			},
		},
		Ext: openrtb.Extensions{
			"requestExt": adapter.Ext.RequestExt,
		},
	}
	var resp openrtb.Response
	err = s.SendRequest(req, &resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dec := json.NewEncoder(w)
	err = dec.Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
