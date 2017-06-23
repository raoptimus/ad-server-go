// GENERATED! DO NOT EDIT!

package main

import (
	"tc/stat"
	"time"
)

type (
	statsByCampaignsHourlyTable struct {
		data map[statsByCampaignsHourlyKey]*statsByCampaignsHourly
	}

	statsByCampaignsHourly struct {
		advCounters
		statsByCampaignsHourlyKey `table:"key"`
		Costs                     float64
		Tax                       float64
		IsWebmaster               bool `table:"noupdate"`
		PaymentType               int  `table:"noupdate"`
		CampaignTypeId            int  `table:"noupdate"`
		UserId                    int  `table:"noupdate"`
		FriendId                  int  `table:"noupdate"`
	}
	statsByCampaignsHourlyKey struct {
		CampaignId int
		TypeId     int
		ForDate    time.Time
	}
)

func (s *statsByCampaignsHourlyTable) Name() string {
	return "statsByCampaignsHourly"
}

func (s *statsByCampaignsHourlyTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByCampaignsHourlyTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByCampaignsHourlyKey]*statsByCampaignsHourly, config.tableCap)
	}
	for _, e := range stats {
		key := statsByCampaignsHourlyKey{
			e.Campaign.Id,
			e.Site.AdZone.AdCode.Type.ToInt(),
			e.Date,
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsByCampaignsHourly{statsByCampaignsHourlyKey: key}
			record.UserId = e.Campaign.UserId
			record.IsWebmaster = e.Campaign.IsWebmaster
			record.CampaignTypeId = e.Campaign.Type.ToInt()
			record.PaymentType = e.PaymentType.ToInt()
			record.FriendId = e.Campaign.BrokerId
			s.data[key] = record
		}
		record.advCounters.add(&e.Counters)
		record.Costs += stat.Float(e.Price + e.Tax + e.Referrals)
		record.Tax += stat.Float(e.Tax + e.Referrals)
	}
}

func (s *statsByCampaignsHourlyTable) UpdateSql() string {
	return makeUpdateSql(statsByCampaignsHourly{}, s.Name())
}

func (s *statsByCampaignsHourlyTable) InsertSql() string {
	return makeInsertSql(statsByCampaignsHourly{}, s.Name())
}
