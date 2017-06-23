package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	mgo "gopkg.in/mgo.v2"
	"log"
	"tc/bootstrap/config"
)

type (
	dataContext struct {
		//main
		Campaigns *mgo.Collection
		Sites     *mgo.Collection
		Ads       *mgo.Collection
		Journal   *mgo.Collection
		Users     *mgo.Collection
		Cities    *mgo.Collection
		AdZones   *mgo.Collection
		AdCodes   *mgo.Collection
		//index
		AdIndex *mgo.Collection
		//pgsql
		PgSqlDb *sqlx.DB
	}
)

var DataContext *dataContext

func connectToMongoDb(url string) *mgo.Database {
	session, err := mgo.Dial(url)
	if err != nil {
		log.Panicln(url, err)
	}

	session.SetMode(mgo.Monotonic, true)
	return session.DB("")
}

func connectToPgSqlDb(url string) *sqlx.DB {
	return sqlx.MustConnect("postgres", config.PgSqlServer())
}

func init() {
	main := connectToMongoDb(config.MongoServerMain())
	index := connectToMongoDb(config.MongoServerAdIndex())
	pgSqlDb := connectToPgSqlDb(config.PgSqlServer())

	DataContext = &dataContext{
		//main
		Campaigns: main.C("Campaigns"),
		Sites:     main.C("Sites"),
		Ads:       main.C("Ads"),
		Journal:   main.C("Journal"),
		Users:     main.C("Users"),
		Cities:    main.C("Cities"),
		AdZones:   main.C("AdZones"),
		AdCodes:   main.C("AdCodes"),
		//index
		AdIndex: index.C("ads_index"),
		//pgsql
		PgSqlDb: pgSqlDb,
	}
}
