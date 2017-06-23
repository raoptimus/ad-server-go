package main

type (
	storeContext struct {
		Users *UserStorage
		Sites *SiteStorage
	}
)

var StoreContext *storeContext

func init() {
	StoreContext = &storeContext{}
	StoreContext.Users = NewUserStorage()
	StoreContext.Sites = NewSiteStorage()

	go StoreContext.Users.Updater()
}
