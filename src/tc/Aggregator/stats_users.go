package main

import (
	"gopkg.in/mgo.v2/bson"
	"tc/stat"
	"time"
)

type (
	statsByUsersTable struct {
		data map[statsByUsersKey]*statsByUsers
	}

	statsByUsers struct {
		statsByUsersKey
		Earn      float64
		Outcome   float64
		EarnRef   float64
		RefUserId int
	}
	statsByUsersKey struct {
		UserId  int
		ForDate time.Time
	}
)

func (s *statsByUsersTable) Name() string {
	return "statsByUsers"
}

func (s *statsByUsersTable) Rows() Rows {
	l := make(Rows, 0, len(s.data))
	for _, row := range s.data {
		l = append(l, row)
	}
	return l
}

func (s *statsByUsersTable) Add(stats stat.Slice) {
	if s.data == nil {
		s.data = make(map[statsByUsersKey]*statsByUsers, config.tableCap)
	}
	for _, e := range stats {
		/* Advertiser */
		key := statsByUsersKey{e.Campaign.UserId, truncateDay(e.Date)}
		record, ok := s.data[key]
		if !ok {
			record = &statsByUsers{statsByUsersKey: key}
			record.RefUserId = e.Site.RefUserId
			s.data[key] = record
		}
		if e.IsPaid() {
			record.Outcome += stat.Float(e.Price + e.Tax + e.Referrals)
		}
		/* Webmaster */
		key = statsByUsersKey{e.Site.UserId, truncateDay(e.Date)}
		record, ok = s.data[key]
		if !ok {
			record = &statsByUsers{statsByUsersKey: key}
			record.RefUserId = e.Site.RefUserId
			s.data[key] = record
		}
		if e.IsPaid() {
			record.Earn += stat.Float(e.Price)
			if e.Site.RefUserId > RefAnonymous {
				record.EarnRef += stat.Float(e.Referrals)
			}
		}
	}
}

func (s *statsByUsersTable) UpdateSql() string {
	return `SELECT * FROM tc."tc_Payments_InsertOrUpdateSystemPayment2"(:userid, :earn, :outcome, :earnref, :fordate, :refuserid)`
}

func (s *statsByUsersTable) InsertSql() string {
	return "SELECT 1" //must not be executed
}

func (s *statsByUsers) Update() {
	inc := s.Earn - s.Outcome
	if inc == 0 {
		return
	}
	if config.debug.disableBalance {
		inc = 0
	}
	cmd := bson.D{
		{"findAndModify", "Users"},
		{"query", bson.M{"_id": s.UserId}},
		{"update", bson.M{
			"$inc": bson.M{
				"Balance": inc,
			},
		}},
		{"fields", bson.M{
			"Balance":       true,
			"Notifications": true,
			"Profile":       true,
		}},
		{"new", true},
	}
	result := struct {
		Value struct {
			Balance       float64 `bson:"Balance"`
			Notifications struct {
				LowBalance bool `bson:"LowBalance"`
			} `bson:"Notifications"`
			Profile struct {
				Notifyminbalance float64 `bson:"notifyminbalance"`
			} `bson:"Profile"`
		}
	}{}
	err := mongo.main.Run(cmd, &result)
	if checkMongoErr(err, mongo.main) != nil {
		return
	}
	user := result.Value

	if user.Notifications.LowBalance || user.Balance > user.Profile.Notifyminbalance {
		return
	}
	notification := bson.M{"$set": bson.M{"Notifications.LowBalance": true}}
	checkMongoErr(mongo.main.C("Users").UpdateId(s.UserId, notification), mongo.main)
}
