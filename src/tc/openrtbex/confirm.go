package openrtbex

import (
	"github.com/bsm/openrtb"
)

type Confirm struct {
	Results []ConfirmBid       `json:"results"`
	Ext     openrtb.Extensions `json:"ext"`
}

type ConfirmBid struct {
	BidId               *string       `json:"bid_id"`
	Nurl                *string       `json:"nurl"`
	ClearingPriceMicros *int          `json:"clearing_price_micros"`
	WinningBidMicros    *int          `json:"winning_bid_micros"`
	LossReason          LossReason    `json:"loss_reason"`
	ErrorReason         ErrorReason   `json:"error_reason"`
	Status              ConfirmStatus `json:"status"`
}

func (bid *ConfirmBid) SetBidId(id string) {
	if bid.BidId == nil {
		bid.BidId = new(string)
	}

	*bid.BidId = id
}
