// generated from tables.go
package main

import (
	"tc/stat"
	"time"
)

type (
	statsByCountryForCampaignsTable struct {
		data map[statsByCountryForCampaignsKey]*statsByCountryForCampaigns
	}

	statsByCountryForCampaigns struct {
		advCounters
		statsByCountryForCampaignsKey `table:"key"`
		Price                         float64
		Tax                           float64
		IsWebmaster                   bool `table:"noupdate"`
		PaymentType                   int  `table:"noupdate"`
		CampaignTypeId                int  `table:"noupdate"`
		UserId                        int  `table:"noupdate"`
		FriendId                      int  `table:"noupdate"`
	}
	statsByCountryForCampaignsKey struct {
		CampaignId int
		CountryId  int
		TypeId     int
		ForDate    time.Time
	}
)

func (s *statsByCountryForCampaignsTable) Name() string {
	return "statsByCountryForCampaigns"
}

func (s *statsByCountryForCampaignsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByCountryForCampaignsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByCountryForCampaignsKey]*statsByCountryForCampaigns, config.tableCap)
	}
	for _, e := range stats {
		key := statsByCountryForCampaignsKey{
			e.Campaign.Id,
			e.Geo.CountryId,
			e.Campaign.Type.ToInt(),
			truncateDay(e.Date),
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsByCountryForCampaigns{statsByCountryForCampaignsKey: key}
			record.UserId = e.Campaign.UserId
			record.IsWebmaster = e.Campaign.IsWebmaster
			record.CampaignTypeId = e.Campaign.Type.ToInt()
			record.PaymentType = e.PaymentType.ToInt()
			record.FriendId = e.Campaign.BrokerId
			s.data[key] = record
		}
		record.advCounters.add(&e.Counters)
		record.Price += stat.Float(e.Price + e.Tax + e.Referrals)
		record.Tax += stat.Float(e.Tax + e.Referrals)
	}
}

func (s *statsByCountryForCampaignsTable) UpdateSql() string {
	return makeUpdateSql(statsByCountryForCampaigns{}, s.Name())
}

func (s *statsByCountryForCampaignsTable) InsertSql() string {
	return makeInsertSql(statsByCountryForCampaigns{}, s.Name())
}
