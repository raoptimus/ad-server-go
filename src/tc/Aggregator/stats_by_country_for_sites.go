// generated from tables.go
package main

import (
	tc "tc/openrtbex"
	"tc/stat"
	"time"
)

type (
	statsByCountryForSitesTable struct {
		data map[statsByCountryForSitesKey]*statsByCountryForSites
	}

	statsByCountryForSites struct {
		stat.Counters             //must be first due to alignment
		statsByCountryForSitesKey `table:"key"`
		Price                     float64
		Tax                       float64
		IsAllowAdult              bool `table:"noupdate"`
		CampaignTypeId            int  `table:"noupdate"`
		CategoryId                int  `table:"noupdate"`
		UserId                    int  `table:"noupdate"`
	}
	statsByCountryForSitesKey struct {
		SiteId    int
		CountryId int
		TypeId    int
		ForDate   time.Time
	}
)

func (s *statsByCountryForSitesTable) Name() string {
	return "statsByCountryForSites"
}

func (s *statsByCountryForSitesTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByCountryForSitesTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByCountryForSitesKey]*statsByCountryForSites, config.tableCap)
	}
	for _, e := range stats {
		key := statsByCountryForSitesKey{
			e.Campaign.Id,
			e.Geo.CountryId,
			e.Campaign.Type.ToInt(),
			truncateDay(e.Date),
		}
		record, ok := s.data[key]
		if !ok {
			record = &statsByCountryForSites{statsByCountryForSitesKey: key}
			record.UserId = e.Site.UserId
			record.CampaignTypeId = e.Campaign.Type.ToInt()
			record.IsAllowAdult = e.Site.Category == tc.CategoryAdult
			record.CategoryId = e.Site.Category.ToInt()

			s.data[key] = record
		}
		if e.IsPaid() {
			record.Price += stat.Float(e.Price)
			record.Tax += stat.Float(e.Tax)
		}
	}
}

func (s *statsByCountryForSitesTable) UpdateSql() string {
	return makeUpdateSql(statsByCountryForSites{}, s.Name())
}

func (s *statsByCountryForSitesTable) InsertSql() string {
	return makeInsertSql(statsByCountryForSites{}, s.Name())
}
