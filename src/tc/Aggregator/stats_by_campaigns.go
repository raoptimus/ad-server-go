/*
(let ((class "statsByCampaignsHourly")
      (file "stats_by_campaigns_hourly.go"))
  (copy-file "stats_by_campaigns.go" file 1)
  (pop-to-buffer (find-file file))
  (delete-region (point-min) (search-forward "*\/"))
  (insert "// GENERATED! DO NOT EDIT!\n")
  (while (search-forward "statsByCampaigns" nil t)
    (replace-match class nil t)))
*/
package main

import (
	"tc/stat"
	"time"
)

type (
	advCounters struct {
		ShowCount        int64
		UniqueShowCount  int64
		ClickCount       int64
		UniqueClickCount int64
		DoubleClickCount int64
		BadClickCount    int64
	}
	statsByCampaignsTable struct {
		data map[statsByCampaignsKey]*statsByCampaigns
	}

	statsByCampaigns struct {
		advCounters
		statsByCampaignsKey `table:"key"`
		Costs               float64
		Tax                 float64
		IsWebmaster         bool `table:"noupdate"`
		PaymentType         int  `table:"noupdate"`
		CampaignTypeId      int  `table:"noupdate"`
		UserId              int  `table:"noupdate"`
		FriendId            int  `table:"noupdate"`
	}
	statsByCampaignsKey struct {
		CampaignId int
		TypeId     int
		ForDate    time.Time
	}
)

func (ac *advCounters) add(t *stat.Counters) {
	ac.ShowCount += t.ShowCount +
		t.FreeShowCount +
		t.WmShowCount
	ac.UniqueShowCount += t.UniqueShowCount +
		t.UniqueFreeShowCount +
		t.UniqueWmShowCount
	ac.ClickCount += t.ClickCount +
		t.FreeClickCount +
		t.WmClickCount
	ac.UniqueClickCount += t.UniqueClickCount +
		t.UniqueFreeClickCount +
		t.UniqueWmClickCount
	ac.BadClickCount += t.BadClickCount
	ac.DoubleClickCount += t.DoubleClickCount
}

func (s *statsByCampaignsTable) Name() string {
	return "statsByCampaigns"
}

func (s *statsByCampaignsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByCampaignsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByCampaignsKey]*statsByCampaigns, config.tableCap)
	}
	for _, e := range stats {
		key := statsByCampaignsKey{
			e.Campaign.Id,
			e.Site.AdZone.AdCode.Type.ToInt(),
			truncateDay(e.Date),
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsByCampaigns{statsByCampaignsKey: key}
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

func (s *statsByCampaignsTable) UpdateSql() string {
	return makeUpdateSql(statsByCampaigns{}, s.Name())
}

func (s *statsByCampaignsTable) InsertSql() string {
	return makeInsertSql(statsByCampaigns{}, s.Name())
}
