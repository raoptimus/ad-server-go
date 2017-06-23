// GENERATED! DO NOT EDIT!

package main

import (
	tc "tc/openrtbex"
	"tc/stat"
	"time"
)

type (
	statsBySitesHourlyTable struct {
		data map[statsBySitesHourlyKey]*statsBySitesHourly
	}

	statsBySitesHourly struct {
		stat.Counters         //must be first due to alignment
		statsBySitesHourlyKey `table:"key"`
		Earnings              float64
		Tax                   float64
		IsAllowAdult          bool `table:"noupdate"`
		CategoryId            int  `table:"noupdate"`
		UserId                int  `table:"noupdate"`
		CampaignTypeId        int  `table:"noupdate"`
	}
	statsBySitesHourlyKey struct {
		SiteId  int
		TypeId  int
		ForDate time.Time
	}
)

func (s *statsBySitesHourlyTable) Name() string {
	return "statsBySitesHourly"
}

func (s *statsBySitesHourlyTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsBySitesHourlyTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsBySitesHourlyKey]*statsBySitesHourly, config.tableCap)
	}
	for _, e := range stats {
		key := statsBySitesHourlyKey{
			e.Site.Id,
			e.Site.AdZone.AdCode.Type.ToInt(),
			e.Date,
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsBySitesHourly{statsBySitesHourlyKey: key}
			record.UserId = e.Site.UserId
			record.CampaignTypeId = e.Campaign.Type.ToInt()
			record.IsAllowAdult = e.Site.Category == tc.CategoryAdult
			record.CategoryId = e.Site.Category.ToInt()

			s.data[key] = record
		}

		record.Counters.Add(&e.Counters)

		if e.IsPaid() {
			record.Earnings += stat.Float(e.Price)
			record.Tax += stat.Float(e.Tax)
		}
	}
}

func (s *statsBySitesHourlyTable) UpdateSql() string {
	return makeUpdateSql(statsBySitesHourly{}, s.Name())
}

func (s *statsBySitesHourlyTable) InsertSql() string {
	return makeInsertSql(statsBySitesHourly{}, s.Name())
}
