package main

import (
	"fmt"
	"github.com/bsm/openrtb"
	"log"
	"net/rpc"
	"runtime/debug"
	"sync"
	"tc/openrtbex"
	"tc/subscribe"
	"time"
)

type (
	Auction struct {
		sync.RWMutex
		wg          sync.WaitGroup
		subscribers *subscribe.Subscribers
		success     []*openrtb.Response
		fail        []*openrtb.Response
	}
)

func NewAuction(subscribers *subscribe.Subscribers) *Auction {
	return &Auction{
		subscribers: subscribers,
		success:     make([]*openrtb.Response, 0),
		fail:        make([]*openrtb.Response, 0),
	}
}

// at - 1, 2
func (s *Auction) Play(req *openrtb.Request, limit, at int) (win BidList, loss BidList) {
	s.requestAll(req)

	win = make(BidList, 0)
	loss = make(BidList, 0)
	bidList := make(BidList, 0)

	s.RLock()

	for _, resp := range s.success {
		for i := 0; i < len(resp.Seatbid); i++ {
			for k := 0; k < len(resp.Seatbid[i].Bid); k++ {
				bidList = append(bidList, &resp.Seatbid[i].Bid[k])
			}
		}
	}

	s.RUnlock()

	bidList = bidList.Sort()
	var winPrice float32

	for i, bid := range bidList {
		status := &openrtbex.BidStatusExt{
			Status:        openrtbex.ConfirmStatusWin,
			ClearingPrice: bid.Price,
			WinningPrice:  &winPrice,
		}

		if len(win) < limit {
			if i+1 == at {
				winPrice = *bid.Price
			}

			status.Status = openrtbex.ConfirmStatusWin
			bid.Ext["bidStatusExt"] = status
			win = append(win, bid)
			continue
		}

		status.Status = openrtbex.ConfirmStatusLoss
		status.LossReason = openrtbex.LossReasonPrice
		bid.Ext["bidStatusExt"] = status
		loss = append(loss, bid)
	}

	return
}

func (s *Auction) GetResults() (success []*openrtb.Response, fail []*openrtb.Response) {
	s.RLock()
	defer s.RUnlock()

	for _, sr := range s.success {
		success = append(success, sr)
	}

	for _, fr := range s.fail {
		fail = append(fail, fr)
	}

	return
}

func (s *Auction) confirm(r *openrtb.Response, background bool) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}

		do(measure("bid-request"))
	}()

	if !background {
		defer s.wg.Done()
	}

	sub := r.Ext["sub"].(*subscribe.Subscriber)

	resp := &openrtb.Response{}
	*resp = *r
	resp.Ext = openrtb.Extensions{
		"responseExt": r.Ext["responseExt"],
	}

	start := time.Now()

	//todo persistent connection for all request or pool
	conn, err := rpc.Dial(sub.Network, sub.Address)

	if err != nil {
		log.Println(sub.Name, "Dial error:", err)
		//todo write error to bid or r
		return
	}

	success := false
	err = conn.Call("Controller.SendConfirm", resp, &success)

	if err != nil {
		log.Println(sub.Name, err)
		//todo write error to bid or r
		return
	}

	log.Println("Confirm ", sub.Name, " elapsed:", time.Now().Sub(start))
}

func (s *Auction) Confirm(background bool, req *openrtb.Request) {
	for _, resp := range s.success {
		//		winCount := 0
		//		for i := 0; i < len(resp.Seatbid); i++ {
		//			for k := 0; k < len(resp.Seatbid[i].Bid); k++ {
		//				st := resp.Seatbid[i].Bid[k].Ext["bidStatusExt"].(*openrtbex.BidStatusExt)
		//
		//				if st.Status == openrtbex.ConfirmStatusWin {
		//					winCount++
		//				}
		//			}
		//		}
		//		sub := resp.Ext["sub"].(*Subscriber)
		//
		//		if winCount == 0 && !sub.LossConfirm {
		//			return
		//		}

		if !background {
			s.wg.Add(1)
		}

		resp.Ext["responseExt"] = &openrtbex.ResponseExt{
			UserId:   req.User.Id,
			DeviceIp: req.Device.Ip,
		}

		go s.confirm(resp, background)
	}

	if !background {
		//todo timeout
		s.wg.Wait()
	}
}

func (s *Auction) addSuccess(sub *subscribe.Subscriber, resp *openrtb.Response) {
	s.Lock()
	defer s.Unlock()

	resp.Ext["sub"] = sub
	s.success = append(s.success, resp)
}

func (s *Auction) addFail(sub *subscribe.Subscriber, resp *openrtb.Response) {
	s.Lock()
	defer s.Unlock()

	resp.Ext["sub"] = sub
	s.fail = append(s.fail, resp)
}

func (s *Auction) request(sub *subscribe.Subscriber, req *openrtb.Request) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("PANIC:", r)
			return
		}

		do(measure("bid-request"))
	}()
	defer s.wg.Done()

	start := time.Now()

	//todo persistent connection for all request or pool
	conn, err := rpc.Dial(sub.Network, sub.Address)

	if err != nil {
		log.Println("Dial error:", err)
		resp := &openrtb.Response{
			Ext: openrtb.Extensions{
				"debugExt": openrtbex.DebugExt{
					Error:      err.Error(),
					CustomData: fmt.Sprintf("%v", time.Now().Sub(start)),
				},
			},
		}

		s.addFail(sub, resp)
		return
	}

	var resp openrtb.Response
	err = conn.Call("Controller.SendRequest", req, &resp)

	if err != nil {
		resp.Ext = openrtb.Extensions{
			"debugExt": openrtbex.DebugExt{
				Error:      err.Error(),
				CustomData: fmt.Sprintf("%v", time.Now().Sub(start)),
			},
		}

		s.addFail(sub, &resp)
		return
	}

	if resp.Ext == nil {
		resp.Ext = openrtb.Extensions{}
	}

	if d, ok := resp.Ext["debugExt"]; ok {
		dExt := d.(openrtbex.DebugExt)
		dExt.CustomData = fmt.Sprintf("%v", time.Now().Sub(start))
		resp.Ext["debugExt"] = dExt

		if dExt.Error != "" {
			s.addFail(sub, &resp)
			return
		}
	}

	s.addSuccess(sub, &resp)
}

func (s *Auction) requestAll(req *openrtb.Request) {
	dsp := s.subscribers.GetAll()
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	tn := reqExt.CodeTypeId.AdCodeTypeNew()

	for _, sub := range dsp {
		//allow all
		if len(sub.AllowTypes) == 0 {
			s.wg.Add(1)
			go s.request(sub, req)
			continue
		}
		for _, t := range sub.AllowTypes {
			if t != tn {
				continue
			}

			s.wg.Add(1)
			go s.request(sub, req)
			break
		}
	}

	//todo timeout
	s.wg.Wait()
}
