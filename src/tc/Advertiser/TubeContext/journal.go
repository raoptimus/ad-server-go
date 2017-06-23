package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"tc/data"
	"time"
)

const (
	CampaignStart = iota
	CampaignStop
	CampaignCreate
	CampaignDelete
	CampaignUpdate
	CampaignLimit
	AdStart
	AdStop
	AdCreate
	AdDelete
	AdUpdate
	AdDecline
	AdApprove
	UnlimitAllCampaigns
	SiteStart
	SiteStop
	SiteCreate
	SiteDelete
	SiteUpdate
	SiteDecline
	SiteApproved
	AdZoneStart2
	AdZoneStop
	AdZoneCreate
	AdZoneDelete
	AdZoneUpdate
	AdZoneStart
	UserCreate
	UserDelete
	UserUpdate
	UserBanned
	UserUnBanned
)

type (
	Journal struct {
		Id          int `bson:"_id"`
		ObjectId    int `bson:"ObjectId"`
		OperationId int `bson:"OperationId"`
	}
	JournalUpdater struct {
		LastId int `bson:"_id"`
	}
)

func NewJournalUpdater() *JournalUpdater {
	j := &JournalUpdater{}
	err := data.DataContext.Journal.Find(nil).Select(bson.M{"_id": 1}).Sort("-_id").One(j)

	if err != nil && err != mgo.ErrNotFound {
		log.Panicln("", err)
	}

	log.Println("Journal last id on load: ", j.LastId)

	j.start()
	return j
}

func (j *JournalUpdater) start() {
	go func() {
		for {
			time.Sleep(time.Minute * 1)
			j.update()
		}
	}()
}

func (j *JournalUpdater) update() {
	log.Println("Journal: starting update")

	query := data.DataContext.Journal.
		Find(bson.M{"_id": bson.M{"$gt": j.LastId}}).
		Select(bson.M{"_id": 1, "ObjectId": 1, "OperationId": 1})

	rc, err := query.Count()

	if err != nil {
		log.Println("Journal: Error read: ", err)
		return
	}

	if rc == 0 {
		return
	}

	log.Println("Journal: New records found ", rc)
	iter := query.Sort("_id").Iter()
	var record Journal

	for iter.Next(&record) {
		j.apply(&record)
		j.LastId = record.Id
	}
}

func (j *JournalUpdater) apply(jr *Journal) {
	log.Printf("Journal: applying: %+v\n", jr)
	id := jr.ObjectId

	switch jr.OperationId {
	case CampaignStart, CampaignCreate, CampaignUpdate:
		{
			StoreContext.Campaigns.LoadById(id)
		}
	case CampaignDelete, CampaignStop:
		{
			StoreContext.Campaigns.Delete(id)
		}
	case CampaignLimit:
		{
			StoreContext.Campaigns.Limit(id)
		}
	case AdStart, AdCreate, AdUpdate, AdApprove:
		{
			StoreContext.Campaigns.LoadReloadAd(id)
		}
	case AdStop, AdDelete:
		{
			StoreContext.Campaigns.DeleteAd(id)
		}
	case UnlimitAllCampaigns:
		{
			StoreContext.Campaigns.UnlimitAll()
		}
	default:
		//ignore
	}
	log.Printf("%+v applied\n", jr)
}
