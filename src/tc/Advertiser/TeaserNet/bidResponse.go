package main

import (
	"fmt"
	"github.com/bsm/openrtb"
	"log"
	"tc/openrtbex"
)

var CUR string = "RUB"

const USER_ID = 5858
const CAMPAIGN_ID = 36887
const BID_PRIORITY = 3

type BidResponse struct {
	Error   string   `json:"error",omitempty`
	Hash    string   `json:"hash"`
	Teasers []Teaser `json:"teasers"`
}

type Teaser struct {
	Id    int     `json:"id"`
	Price float32 `json:"price,string"`
	Ctr   float32 `json:"ctr,string"`
	Title string  `json:"title"`
	Img   string  `json:"img"`
	Url   string  `json:"url"`
}

func (s *BidResponse) ToRTBResponse(req *BidRequest, resp *openrtb.Response) {
	resp.Customdata = &s.Hash
	seatbid := openrtb.Seatbid{}

	seatbid.Bid = make([]openrtb.Bid, len(s.Teasers))

	for i, teaser := range s.Teasers {
		bid := openrtb.Bid{}
		id := fmt.Sprintf("%d", teaser.Id)
		bid.Adid = &id
		bid.Id = &id
		ctr := teaser.Ctr * 100 //в долях

		if ctr == 0 {
			ctr = 0.01 //TODO fixed problem
		}

		if ctr > 7 {
			log.Println("HIGH CTR", req)
		}

		teaser.Price = teaser.Price * 0.6 //40% маржа тизернета
		teaser.Ctr = ctr
		cpm := (teaser.Price * 10.0 * teaser.Ctr)
		bid.Price = &cpm

		bid.Ext = openrtb.Extensions{
			"bidExt": openrtbex.BidExt{
				PaymentType: openrtbex.PaymentTypeCPC,
				Title:       teaser.Title,
				Image:       teaser.Img,
				Url:         teaser.Url,
				Price:       teaser.Price,
				Ctr:         teaser.Ctr,
				Priority:    BID_PRIORITY,
				UserId:      USER_ID,
				CampaignId:  CAMPAIGN_ID,
			},
		}

		seatbid.Bid[i] = bid
	}

	resp.Seatbid = []openrtb.Seatbid{seatbid}
	resp.Cur = &CUR
}
