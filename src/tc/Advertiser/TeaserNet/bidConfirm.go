package main

import (
	"github.com/bsm/openrtb"
	"tc/openrtbex"
)

type BidConfirm struct {
	Action  string    `json:"action"`
	Charset string    `json:"charset"`
	Hash    *string   `json:"hash"`
	Teasers []*string `json:"teasers"`
}

const (
	BID_CONFIRM_ACTION  = "confirmView"
	BID_CONFIRM_CHARSET = "utf-8"
)

func NewBidConfirm(resp *openrtb.Response, ss *Session) *BidConfirm {
	bidConfirm := &BidConfirm{}
	bidConfirm.Action = BID_CONFIRM_ACTION
	bidConfirm.Charset = BID_CONFIRM_CHARSET
	bidConfirm.Hash = resp.Customdata
	bidConfirm.Teasers = make([]*string, 0)

	for _, seatbid := range resp.Seatbid {
		for _, bid := range seatbid.Bid {
			status := bid.Ext["bidStatusExt"].(openrtbex.BidStatusExt)

			if status.Status != openrtbex.ConfirmStatusWin {
				continue
			}

			bidConfirm.Teasers = append(bidConfirm.Teasers, bid.Id)
			ss.setAdAsShown(bid.Id)
		}
	}

	return bidConfirm
}
