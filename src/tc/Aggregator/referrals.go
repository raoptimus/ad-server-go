package main

import (
	"tc/stat"
	"time"
)

const RefAnonymous = 1

type (
	referralsTable struct {
		data map[referralsKey]*referrals
	}

	referrals struct {
		referralsKey `table:"key"`
		ShowCount    int64
		ClickCount   int64
		Price        float64
	}
	referralsKey struct {
		UserId    int
		RefUserId int
		ForDate   time.Time
	}
)

func (s *referralsTable) Name() string {
	return "referrals"
}

func (s *referralsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *referralsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[referralsKey]*referrals, config.tableCap)
	}
	for _, e := range stats {
		if e.Site.RefUserId <= RefAnonymous {
			continue
		}
		key := referralsKey{
			e.Site.UserId,
			e.Site.RefUserId,
			truncateDay(e.Date),
		}
		record, ok := s.data[key]
		if !ok {
			record = &referrals{referralsKey: key}
			s.data[key] = record
		}
		record.Price += stat.Float(e.Referrals)
		record.ShowCount += e.ShowCount
		record.ClickCount += e.ClickCount
	}
}

func (s *referralsTable) UpdateSql() string {
	return makeUpdateSql(referrals{}, s.Name())
}

func (s *referralsTable) InsertSql() string {
	return makeInsertSql(referrals{}, s.Name())
}
