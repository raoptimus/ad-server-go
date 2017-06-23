package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"tc/stat"
	"time"
)

func saveFull(slice stat.Slice) {
	full := mongo.stats.C(config.mongoStatsCollection)
	for _, e := range slice {
		year, month, day := e.Date.Date()
		e.Date = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

		id := fmt.Sprint(e.Key())
		b := bson.M{}
		b["$setOnInsert"] = e.RawEvent
		b["$inc"] = e.Counters
		_, err := full.UpsertId(id, b)

		checkMongoErr(err, mongo.stats)
	}
}
