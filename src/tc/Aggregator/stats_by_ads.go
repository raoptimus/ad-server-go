package main

import (
	"tc/stat"
	"time"
)

type (
	statsByAdsTable struct {
		data map[statsByAdsKey]*statsByAds
	}

	statsByAds struct {
		advCounters
		statsByAdsKey `table:"key"`
		SumCost       float64
		IsWebmaster   bool `table:"noupdate"`
		UserId        int  `table:"noupdate"`
		CampaignId    int  `table:"noupdate"`
	}
	statsByAdsKey struct {
		AdId    int
		ForDate time.Time
	}
)

func (s *statsByAdsTable) Name() string {
	return "statsByAds"
}

func (s *statsByAdsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByAdsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByAdsKey]*statsByAds, config.tableCap)
	}
	for _, e := range stats {
		if e.Campaign.Ad.Id == 0 {
			continue
		}
		key := statsByAdsKey{e.Campaign.Ad.Id, truncateDay(e.Date)}
		record, ok := s.data[key]
		if !ok {
			record = &statsByAds{statsByAdsKey: key}
			record.CampaignId = e.Campaign.Id
			record.UserId = e.Campaign.UserId
			record.IsWebmaster = e.Campaign.IsWebmaster
			s.data[key] = record
		}
		record.advCounters.add(&e.Counters)
		record.SumCost += stat.Float(e.Price + e.Tax + e.Referrals)
	}
}

func (s *statsByAdsTable) UpdateSql() string {
	return makeUpdateSql(statsByAds{}, s.Name())
}

func (s *statsByAdsTable) InsertSql() string {
	return makeInsertSql(statsByAds{}, s.Name())
}
