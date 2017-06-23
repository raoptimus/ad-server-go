package main

import (
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"strconv"
	"tc/openrtbex"
)

const AD_SHOW_LIMIT_PER_USER = 4

type Bid struct {
	Ad        *Ad
	Bid       *openrtb.Bid
	Priority  int
	Price     float32
	ShowCount int
}

type BidCtr struct {
	Id          string  `bson:"_id"`
	Ctr         float64 `bson:"Ctr"`
	IsNew       bool    `bson:"IsNew"`
	StatBlockId int     `bson:"StatBlockId"`
}

func NewBid(ad *Ad, req *openrtb.Request, ss *Session) (bid *Bid, err error) {
	if ad.Campaign.Deleted {
		return nil, errors.New(fmt.Sprintf("Campaign %d has been removed", ad.CampaignId))
	}
	if ad.Deleted {
		return nil, errors.New(fmt.Sprintf("Ad %d has been removed", ad.Id))
	}
	if ad.Campaign.IsLimited {
		return nil, errors.New(fmt.Sprintf("Campaign %d isLimited", ad.CampaignId))
	}
	if !ad.Campaign.IsDebit() {
		return nil, errors.New(fmt.Sprintf("User %d has no coin", ad.UserId))
	}

	wasClicked := ss.campWasClicked(ad.CampaignId)

	if !ad.Campaign.IsAllowRaw && wasClicked {
		return nil, errors.New(fmt.Sprintf("Campaign %d was clicked", ad.CampaignId))
	}

	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)

	var showCount int

	if ad.Campaign.TypeId == openrtbex.CampaignTypePopunder {
		if wasClicked {
			showCount = 1
		}
	} else {
		showCount = ss.getAdShowCount(strconv.Itoa(ad.Id))

		if showCount >= AD_SHOW_LIMIT_PER_USER {
			return nil, errors.New(fmt.Sprintf("Ad %d can't be showed because of the limit impressions per user", ad.Id))
		}
	}

	isNew := false
	statBlockId := 0
	geoPrice := ad.Campaign.GeoPrice(reqExt.GeoId)
	var ctr float32
	var cpm float32
	var cpc float32

	if ad.Campaign.PaymentType == openrtbex.PaymentTypeCPC {
		//получаем ctr
		cpc = geoPrice
		bidCtr := ad.CtrStorage.GetOneOrDefault(ad, &reqExt)
		ctr = float32(bidCtr.Ctr) + 0.001
		cpm = ctr * float32(10) * cpc
		isNew = bidCtr.IsNew && ss.isNewRequest()
		statBlockId = bidCtr.StatBlockId
	} else {
		ctr = 0
		cpm = geoPrice
		cpc = 0
	}

	i := Btou(ad.Campaign.IsBase())<<2 | Btou(isNew)
	priority := priorities[i]

	bidType := openrtbex.BidTypePay

	if ad.Campaign.IsBase() {
		bidType = openrtbex.BidTypeFree
	} else if ad.Campaign.IsWebmaster {
		bidType = openrtbex.BidTypeWm
	}

	w := 0
	h := 0

	switch reqExt.CodeTypeId {
	case
		openrtbex.AdCodeTypeInVideoPauseRoll,
		openrtbex.AdCodeTypeTeasers,
		openrtbex.AdCodeTypeInVideoOverlay,
		openrtbex.AdCodeTypeInEmbedOverlay,
		openrtbex.AdCodeTypeInHtml5VideoPauseRoll,
		openrtbex.AdCodeTypeInHtml5VideoOverlay:
		{
			w = 250
			h = 250
		}
	case
		openrtbex.AdCodeTypeInVideoPostRoll,
		openrtbex.AdCodeTypeInVideoPreRoll,
		openrtbex.AdCodeTypeBanners300x250,
		openrtbex.AdCodeTypeInEmbedPreRoll,
		openrtbex.AdCodeTypeMobileBanners300x250:
		{
			//todo if new player 12traffic and new version then 250x250
			w = 300
			h = 250
		}
	case openrtbex.AdCodeTypeMobileBanners300x100:
		{
			w = 300
			h = 100
		}
	case openrtbex.AdCodeTypeMobileBanners300x50:
		{
			w = 300
			h = 50
		}
	}

	bidExt := openrtbex.BidExt{
		PaymentType:  ad.Campaign.PaymentType,
		Title:        ad.Title,
		Width:        w,
		Height:       h,
		Image:        ad.ImageUrl(w, h),
		Priority:     priority,
		Ctr:          ctr,
		Type:         bidType,
		Price:        geoPrice,
		CampaignId:   int(ad.Campaign.Id),
		AdId:         int(ad.Id),
		StatBlockId:  statBlockId,
		UserId:       int(ad.UserId),
		BrokerId:     ad.Campaign.FriendId,
		CampaignType: ad.Campaign.TypeId,
	}

	adId := strconv.Itoa(ad.Id)

	oBid := &openrtb.Bid{
		Id:    &adId,
		Price: &cpm,
		Ext: openrtb.Extensions{
			"bidExt": bidExt,
		},
	}

	bid = &Bid{
		Ad:        ad,
		Bid:       oBid,
		Priority:  priority,
		Price:     cpm,
		ShowCount: showCount,
	}

	return
}

func Btou(b bool) uint {
	if b {
		return 1
	}
	return 0
}

var priorities = []int{
	// // |Free|Shown|IsNew|
	3, // |	  0|	0|    0|
	1, // |	  0|	0|    1|
	4, // |	  0|	1|    0|
	2, // |	  0|	1|    1|
	7, // |	  1|	0|    0|
	5, // |	  1|	0|    1|
	8, // |	  1|	1|    0|
	6, // |	  1|	1|    1|
}
