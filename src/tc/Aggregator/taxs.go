package main

import (
	"tc/stat"
	"time"
)

type (
	taxsTable struct {
		data map[taxsKey]*taxs
	}

	taxs struct {
		taxsKey  `table:"key"`
		Earnings float64
	}
	taxsKey struct {
		ForDate time.Time
	}
)

func (s *taxsTable) Name() string {
	return "taxs"
}

func (s *taxsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *taxsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[taxsKey]*taxs, config.tableCap)
	}
	for _, e := range stats {
		key := taxsKey{truncateDay(e.Date)}
		record, ok := s.data[key]
		if !ok {
			record = &taxs{taxsKey: key}
			s.data[key] = record
		}
		record.Earnings += stat.Float(e.Tax)
	}
}

func (s *taxsTable) UpdateSql() string {
	return makeUpdateSql(taxs{}, s.Name())
}

func (s *taxsTable) InsertSql() string {
	return makeInsertSql(taxs{}, s.Name())
}
