package main

import "os"

type (
	storeContext struct {
		Sessions  *SessionStorage
		Users     *UserStorage
		Campaigns *CampaignStorage
	}
)

var StoreContext *storeContext

func init() {
	StoreContext = &storeContext{}
	os.Chmod("/tmp/redis.sock", os.ModePerm)
	StoreContext.Sessions = NewSessionStorage("TubeContext", "unix", "/tmp/redis.sock", LIFE_SESSION)
	StoreContext.Users = NewUserStorage()
	StoreContext.Campaigns = NewCampaignStorage()

	go StoreContext.Users.Updater()
}
