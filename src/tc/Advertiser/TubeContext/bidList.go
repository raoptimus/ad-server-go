package main

import (
	"github.com/bsm/openrtb"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"tc/openrtbex"
)

type BidList []*Bid

func NewBidList(adList []*Ad, req *openrtb.Request, ss *Session, limit int, disableRotation bool) BidList {
	bidList := make(BidList, 0)

	for _, ad := range adList {
		bid, err := NewBid(ad, req, ss)

		if err != nil {
			log.Println(err)
			continue
		}

		bidList = append(bidList, bid)
	}

	bidList.Sort()
	if disableRotation {
		if len(bidList) > limit {
			return bidList[:limit]
		}
		return bidList
	}
	return bidList.Uniq(limit, ss)
}

func (s BidList) ToRtbBidList() []openrtb.Bid {
	list := make([]openrtb.Bid, len(s))

	for i, bid := range s {
		list[i] = *bid.Bid
	}

	return list
}

func (s BidList) Len() int {
	return len(s)
}

func (s BidList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s BidList) Less(i, j int) bool {
	if s[i].Priority != s[j].Priority {
		return s[i].Priority < s[j].Priority
	}

	if s[i].ShowCount != s[j].ShowCount {
		return s[i].ShowCount < s[j].ShowCount
	}

	return s[i].Price > s[j].Price
}

func (s BidList) Sort() BidList {
	defer do(measure("list-sort"))
	sort.Sort(s)
	return s
}

func (s BidList) Uniq(limit int, ss *Session) BidList {
	defer do(measure("list-uniq"))

	if len(s) == 0 {
		return s
	}

	firstBid := s[0]

	if firstBid.Ad.Campaign.TypeId == openrtbex.CampaignTypePopunder {
		random := rand.Intn(100) + 1
		probability := 0
		mostExpensive := firstBid

		for _, i := range rand.Perm(len(s)) {
			bid := s[i]
			probability += bid.Ad.Campaign.Percent

			if random <= probability && bid.Ad.Campaign.Percent > 0 {
				return BidList{bid}
			}

			if mostExpensive.Price < bid.Price {
				mostExpensive = bid
			}
		}

		return BidList{mostExpensive}

	} else {
		list := make(BidList, 0)
		campIdList := make(map[int]bool)

		//select all ads with the highest cpm
		for _, bid := range s {
			b, ok := campIdList[bid.Ad.CampaignId]

			if ok && b {
				continue
			}

			if !ss.lockAd(strconv.Itoa(bid.Ad.Id)) {
				continue
			}

			list = append(list, bid)
			campIdList[bid.Ad.CampaignId] = true

			if len(list) == limit {
				break
			}
		}

		return list
	}
}
