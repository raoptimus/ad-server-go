package main

import (
	"errors"
	"github.com/bsm/openrtb"
	"log"
	"net"
	"net/rpc"
	"runtime/debug"
	"tc/openrtbex"
	"tc/subscribe"
)

type ApiController struct {
	*Controller
}

func NewApiController(c *Controller) *ApiController {
	cc := &ApiController{
		Controller: c,
	}
	go cc.listen()
	return cc
}

type Event struct {
	Site struct {
		Id        int
		UserId    int
		RefUserId int
		Category  openrtbex.Category

		AdZone struct {
			Id   int
			Type openrtbex.AdZoneType

			AdCode2 struct {
				Id   int
				Type openrtbex.AdCodeType
			}
		}
	}
	Campaign struct {
		Id       int
		UserId   int
		BrokerId int
		Type     openrtbex.CampaignType

		Ad struct {
			Id int
		}
	}
	Geo struct {
		Id        int
		CountryId int
	}

	Operator openrtbex.Operator
	Os       openrtbex.Os
	Device   openrtbex.Device
	Browser  openrtbex.Browser

	RawEvent struct {
		//            EventKey    `bson:",inline"`
		PaymentType openrtbex.PaymentType
		//            Type        EventType click
	}
}

func (s *ApiController) ClickConfirm(resp *openrtb.Response, unused *bool) error {
	bid := resp.Seatbid[0].Bid[0]
	bidExt := bid.Ext["bidExt"].(openrtbex.BidExt)
	bidStatusExt := bid.Ext["bidStatusExt"].(openrtbex.BidStatusExt)

	if bidExt.BrokerId == 0 {
		sub := s.Subscriber.FindByBroker(bidExt.BrokerId)
		if sub == nil {
			return errors.New("Subscriber not found")
		}
		//todo send click confirm to subscribe
	}
	//todo get isUnique, isDouble from session
	if bidStatusExt.Status == openrtbex.ConfirmStatusError {
		log.Println("Click is bad:" + bidStatusExt.Error) //event type as bad click
	}

	//todo send event click to stats
	e := Event{
		//			Site:   {},
		//			AdZone: {},

		//			Geo: {},

		Campaign: struct {
			Id       int
			UserId   int
			BrokerId int
			Type     openrtbex.CampaignType

			Ad struct {
				Id int
			}
		}{
			Id:       bidExt.CampaignId,
			UserId:   bidExt.UserId,
			BrokerId: bidExt.BrokerId,
			Type:     bidExt.CampaignType,

			Ad: struct {
				Id int
			}{
				Id: bidExt.AdId,
			},
		},
	}
	log.Println(e)
	return nil
}

func (s *ApiController) Subscribe(sub *subscribe.Subscriber, result *bool) error {
	if err := s.Subscriber.Valid(sub); err != nil {
		return err
	}

	s.Subscriber.Subscribe(sub)
	*result = true
	return nil
}

func (s *ApiController) UnSubscribe(name *string, result *bool) error {
	s.Subscriber.UnSubscribe(*name)
	*result = true
	return nil
}

func (s *ApiController) listen() {
	server, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		panic(err)
	}

	rpc.Register(s)

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
			//			defer conn.Close()
			//            conn.SetDeadline(time.Now().Add(600 * time.Millisecond))
			rpc.ServeConn(conn)
		}(conn)
	}
}
