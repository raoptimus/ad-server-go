package main

type AdCode struct {
	Id       int    `bson:"_id"`
	AdZoneId int    `bson:"AdZoneId"`
	TypeId   int    `bson:"TypeId"`
	SiteId   string `bson:"SiteId"`
	Styles   string `bson:"Styles"`
	Title    string `bson:"Title"`
}
