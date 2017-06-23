package main

type (
	Site struct {
		Id       int    `bson:"_id"`
		UserId   int    `bson:"UserId"`
		IsActive bool   `bson:"IsActive"`
		Approved bool   `bson:"Approved"`
		Host     string `bson:"Host"`
	}
)
