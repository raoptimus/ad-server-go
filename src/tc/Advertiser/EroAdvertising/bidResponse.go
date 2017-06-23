package main

import (
	"github.com/bsm/openrtb"
	"strconv"
	"strings"
	"tc/openrtbex"
)

type BidResponse struct {
	Results BidResponseResults `json:"results"`
}

type BidResponseResults struct {
	Ads []Ad `json:"ad"`
}

var CUR string = "RUB"

type pFloat float32

func (f *pFloat) UnmarshalJSON(b []byte) error {
	s := string(b)

	if s == "0" {
		s = "0.0"
	}

	s = strings.Trim(s, "\"")
	ret, err := strconv.ParseFloat(s, 32)

	if err != nil {
		return err
	}

	*f = pFloat(ret)
	return nil
}

func (f *pFloat) toFloat32() float32 {
	return float32(*f)
}

type Ad struct {
	ClickUrl  *string `json:"clickurl"`
	BannerUrl *string `json:"bannerurl"`
	PingUrl   *string `json:"pingurl"`
	Cpm       *pFloat `json:"cpm"`
	Ctr       *pFloat `json:"ctr"`
	Cpc       *pFloat `json:"cpc"`
	IsCpm     *int    `json:"iscpm"`
}

func (s *BidResponse) ToRTBResponse(resp *openrtb.Response, cur float32) {
	seatbid := openrtb.Seatbid{}
	seatbid.Bid = make([]openrtb.Bid, len(s.Results.Ads))

	for i, ad := range s.Results.Ads {
		bid := openrtb.Bid{}
		//		id := fmt.Sprintf("%d", i+1)
		//		bid.Adid = &id
		//		bid.Id = &id
		bid.Nurl = ad.PingUrl
		cpc := ad.Cpc.toFloat32() * cur
		ctr := ad.Ctr.toFloat32()
		cpm := ad.Cpm.toFloat32() * cur

		var paymentType openrtbex.PaymentType

		if *ad.IsCpm == 1 {
			cpc = 0.0
			bid.Price = &cpm
			paymentType = openrtbex.PaymentTypeCPM
		} else {
			cpm = cpc * 10.0 * ctr
			bid.Price = &cpm
			paymentType = openrtbex.PaymentTypeCPC
		}

		bid.Ext = openrtb.Extensions{
			"bidExt": openrtbex.BidExt{
				PaymentType: paymentType,
				Image:       *ad.BannerUrl,
				Url:         *ad.ClickUrl,
				Price:       cpm,
				Ctr:         ctr,
			},
		}

		seatbid.Bid[i] = bid
	}

	resp.Seatbid = []openrtb.Seatbid{seatbid}
	resp.Cur = &CUR
}
