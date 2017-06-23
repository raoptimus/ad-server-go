package main

import (
	"encoding/json"
	"github.com/bsm/openrtb"
	"log"
	"sort"
	"tc/openrtbex"
)

type BidList []*openrtb.Bid

func (s BidList) Sort() BidList {
	defer do(measure("list-sort"))
	sort.Sort(s)
	return s
}

func (s BidList) Len() int {
	return len(s)
}

func (s BidList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s BidList) Less(i, j int) bool {
	e1 := s[i].Ext["bidExt"].(openrtbex.BidExt)
	e2 := s[j].Ext["bidExt"].(openrtbex.BidExt)

	if e1.Priority != e2.Priority {
		return e1.Priority < e2.Priority
	}

	if s[i].Price == nil {
		b, _ := json.Marshal(s[i])
		log.Println(string(b))
	}

	if s[j].Price == nil {
		b, _ := json.Marshal(s[j])
		log.Println(string(b))
	}

	return *s[i].Price > *s[j].Price
}
