/*
(let ((class "statsBySitesHourly")
      (file "stats_by_sites_hourly.go"))
  (copy-file "stats_by_sites.go" file 1)
  (pop-to-buffer (find-file file))
  (delete-region (point-min) (search-forward "*\/"))
  (insert "// GENERATED! DO NOT EDIT!\n")
  (while (search-forward "statsBySites" nil t)
    (replace-match class nil t)))
*/

package main

import (
	tc "tc/openrtbex"
	"tc/stat"
	"time"
)

type (
	statsBySitesTable struct {
		data map[statsBySitesKey]*statsBySites
	}

	statsBySites struct {
		stat.Counters   //must be first due to alignment
		statsBySitesKey `table:"key"`
		Earnings        float64
		Tax             float64
		IsAllowAdult    bool `table:"noupdate"`
		CategoryId      int  `table:"noupdate"`
		UserId          int  `table:"noupdate"`
		CampaignTypeId  int  `table:"noupdate"`
	}
	statsBySitesKey struct {
		SiteId  int
		TypeId  int
		ForDate time.Time
	}
)

func (s *statsBySitesTable) Name() string {
	return "statsBySites"
}

func (s *statsBySitesTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsBySitesTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsBySitesKey]*statsBySites, config.tableCap)
	}
	for _, e := range stats {
		key := statsBySitesKey{
			e.Site.Id,
			e.Site.AdZone.AdCode.Type.ToInt(),
			truncateDay(e.Date),
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsBySites{statsBySitesKey: key}
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

func (s *statsBySitesTable) UpdateSql() string {
	return makeUpdateSql(statsBySites{}, s.Name())
}

func (s *statsBySitesTable) InsertSql() string {
	return makeInsertSql(statsBySites{}, s.Name())
}
