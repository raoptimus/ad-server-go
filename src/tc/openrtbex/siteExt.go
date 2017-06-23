package openrtbex

type SiteExt struct {
	Id           int      `json:"id"`
	UserId       int      `json:"userId"`
	IsPremium    bool     `json:"premium"`
	CategoryId   Category `json:"catId"`
	IsAllowAdult bool     `json:"isAdult"`
}
