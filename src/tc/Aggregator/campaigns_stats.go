package main

import (
	"gopkg.in/mgo.v2/bson"
	"tc/mongodb"
	"tc/stat"
	"time"
)

type (
	campaignsStatsTable struct {
		data map[campaignsStatsKey]*campaignsStats
	}

	campaignsStats struct {
		campaignsStatsKey `table:"key"`
		ClickCountToday   int64
		ShowCountToday    int64
		SumCostToday      float64
		ResetedDate       time.Time `table:"noupdate"`
	}
	campaignsStatsKey struct {
		CampaignId int
	}
)

//force interface check
var _campaignsStatsTable ComplexTable = &campaignsStatsTable{}

func (s *campaignsStatsTable) Name() string {
	return "campaigns_stats"
}

func (s *campaignsStatsTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *campaignsStatsTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[campaignsStatsKey]*campaignsStats, config.tableCap)
	}
	for _, e := range stats {
		key := campaignsStatsKey{e.Campaign.Id}
		record, ok := s.data[key]
		if !ok {
			record = &campaignsStats{campaignsStatsKey: key}
			s.data[key] = record
			record.ResetedDate = truncateDay(e.Date)
		}
		record.ClickCountToday += e.ClickCount
		record.ShowCountToday += e.ShowCount
		record.SumCostToday += stat.Float(e.Price + e.Tax + e.Referrals)
	}
}

func (s *campaignsStatsTable) UpdateSql() string {
	return makeUpdateSql(campaignsStats{}, s.Name())
}

func (s *campaignsStatsTable) InsertSql() string {
	return makeInsertSql(campaignsStats{}, s.Name())
}

func (s *campaignsStatsTable) BeforeSql() string {
	return `UPDATE "tc"."campaigns_stats"
                SET "ClickCountToday" = 0, "SumCostToday" = 0, "ShowCountToday" = 0, "ResetedDate" = :reseteddate
               WHERE "CampaignId" = :campaignid AND "ResetedDate" < :reseteddate`
}

// db.runCommand({findAndModify: "Campaigns", query: {_id: 17}, update: {$inc: {ClickCount: 13}}, fields: {IsLimited: 1, ClickCount: 1, MaxClickCountPerDay: 1}, new: true});
func (c *campaignsStats) Update() {
	if config.debug.disableLimits {
		return
	}
	cmd := bson.D{
		{"findAndModify", "Campaigns"},
		{"query", bson.M{"_id": c.CampaignId}},
		{"update", bson.M{
			"$inc": bson.M{
				"ClickCountToday": c.ClickCountToday,
				"ShowCountToday":  c.ShowCountToday,
				"SumCostToday":    c.SumCostToday,
			},
		}},
		{"fields", bson.M{
			"IsLimited":           true,
			"ClickCountToday":     true,
			"ShowCountToday":      true,
			"SumCostToday":        true,
			"MaxClickCountPerDay": true,
			"MaxShowCountPerDay":  true,
			"MaxCostPerDay":       true,
		}},
		{"new", true},
	}
	result := struct {
		Value struct {
			IsLimited           bool `bson:"IsLimited"`
			ClickCountToday     int  `bson:"ClickCountToday"`
			ShowCountToday      int  `bson:"ShowCountToday"`
			SumCostToday        int  `bson:"SumCostToday"`
			MaxClickCountPerDay int  `bson:"MaxClickCountPerDay"`
			MaxShowCountPerDay  int  `bson:"MaxShowCountPerDay"`
			MaxCostPerDay       int  `bson:"MaxCostPerDay"`
		}
	}{}
	err := mongo.main.Run(cmd, &result)
	if checkMongoErr(err, mongo.main) != nil {
		return
	}
	limited := func(cur, max int) bool {
		return max != 0 && cur >= max
	}
	limits := result.Value
	isLimited := limited(limits.ClickCountToday, limits.MaxClickCountPerDay) ||
		limited(limits.ShowCountToday, limits.MaxShowCountPerDay) ||
		limited(limits.SumCostToday, limits.MaxCostPerDay)

	if isLimited != limits.IsLimited {
		set := bson.M{"$set": bson.M{"IsLimited": isLimited}}
		if checkMongoErr(mongo.main.C("Campaigns").UpdateId(c.CampaignId, set), mongo.main) != nil {
			return
		}
		journal := mongo.main.C("Journal")
		id, err := mongodb.GetNewIncId(journal)
		if checkMongoErr(err, mongo.main) != nil {
			return
		}
		record := bson.M{
			"OperationId": 5, // limit; see JournalOperationId.php
			"ObjectId":    c.CampaignId,
			"AddedDate":   time.Now().UTC(),
			"_id":         id,
		}
		if checkMongoErr(journal.Insert(record), mongo.main) != nil {
			return
		}
	}
}
