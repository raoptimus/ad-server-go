package main

import (
	"github.com/bsm/openrtb"
	"tc/openrtbex"
)

const (
	CATEGORY_ADULT string  = "IAB25-3"
	MIN_BID        float32 = 0.01
)

type BidRequest struct {
	Action  string `json:"action"`
	Charset string `json:"charset"`
	Width   *int   `json:"width"`
	MiscId  *int   `json:"misc_id"`
	BlockId *int   `json:"block_id"`
	SiteId  *int   `json:"site_id"`
	Site    *Site  `json:"site"`
	User    *User  `json:"user"`
	Ad      *Ad    `json:"ad"`
}

type Site struct {
	Referer *string `json:"referer"`
	Page    *string `json:"page"`
}

type User struct {
	Ua   *string `json:"ua"`
	Ip   *string `json:"ip"`
	Lang *string `json:"lang"`
}

type Ad struct {
	Show      []int    `json:"show,omitempty"`
	Amount    int      `json:"amount"`
	MinBid    *float32 `json:"min_bid"`
	NoContent bool     `json:"no_content"`
}

const (
	BID_REQUEST_ACTION  = "getTeasers"
	BID_REQUEST_CHARSET = "utf-8"
)

func NewBidRequest(request *openrtb.Request, ss *Session) *BidRequest {
	req := &BidRequest{}
	req.Action = BID_REQUEST_ACTION
	req.Charset = BID_REQUEST_CHARSET
	req.Width = request.Imp[0].Banner.W
	reqExt := request.Ext["requestExt"].(openrtbex.RequestExt)
	req.MiscId = &reqExt.CodeId
	req.Site = new(Site)
	req.Site.Page = request.Site.Page
	//	req.Site.Referer = request.Site.Ref
	req.User = new(User)
	req.User.Ua = request.Device.Ua
	req.User.Ip = request.Device.Ip
	req.User.Lang = request.Device.Language
	req.Ad = new(Ad)
	req.Ad.NoContent = true
	impExt := request.Imp[0].Ext["impExt"].(openrtbex.ImpExt)
	req.Ad.Amount = impExt.Count
	req.Ad.Show = ss.getAdShownList()

	//	minBid := *request.Imp[0].Bidfloor / 10.0
	var minBid float32 = MIN_BID
	req.Ad.MinBid = &minBid

	isAdult := StringArray(request.Site.Cat).Contains(CATEGORY_ADULT)
	isPremium := request.Site.Ext["siteExt"].(openrtbex.SiteExt).IsPremium
	typeId := reqExt.CodeTypeId
	isMobile := request.Device.DeviceType() == openrtb.DEVICE_TYPE_MOBILE

	var blockId int
	var siteId int

	if isMobile {
		if isAdult {
			blockId = 477319
			siteId = 227135
		} else {
			blockId = 477321
			siteId = 227138
		}
	} else {
		if isAdult {
			if isPremium {
				switch typeId {
				case openrtbex.AdCodeTypeInEmbedOverlay:
					{
						blockId = 550396
						siteId = 247040
					}

				case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
					{
						blockId = 529151
						siteId = 241852
					}

				case openrtbex.AdCodeTypeInHtml5VideoPauseRoll, openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInEmbedPreRoll:
					{
						blockId = 529153
						siteId = 241853
					}

				case openrtbex.AdCodeTypeTeasers:
					{
						blockId = 529154
						siteId = 241854
					}

				default:
					{
						blockId = 498175
						siteId = 218977
					}
				}
			} else {
				switch typeId {
				case openrtbex.AdCodeTypeInEmbedOverlay:
					{
						blockId = 550398
						siteId = 247042
					}

				case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
					{
						blockId = 529155
						siteId = 241855
					}

				case openrtbex.AdCodeTypeInHtml5VideoPauseRoll, openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInEmbedPreRoll:
					{
						blockId = 529156
						siteId = 241856
					}

				case openrtbex.AdCodeTypeTeasers:
					{
						blockId = 529157
						siteId = 241857
					}

				default:
					{
						blockId = 280439
						siteId = 115392
					}
				}
			}
		} else {
			if isPremium {
				switch typeId {
				case openrtbex.AdCodeTypeInEmbedOverlay:
					{
						blockId = 550394
						siteId = 247039
					}

				case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
					{
						blockId = 529164
						siteId = 241867
					}

				case openrtbex.AdCodeTypeInHtml5VideoPauseRoll, openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInEmbedPreRoll:
					{
						blockId = 529166
						siteId = 241868
					}

				case openrtbex.AdCodeTypeTeasers:
					{
						blockId = 529167
						siteId = 241869
					}

				default:
					{
						blockId = 286015
						siteId = 117855
					}
				}
			} else {
				switch typeId {
				case openrtbex.AdCodeTypeInEmbedOverlay:
					{
						blockId = 550397
						siteId = 247041
					}

				case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
					{
						blockId = 529168
						siteId = 241870
					}

				case openrtbex.AdCodeTypeInHtml5VideoPauseRoll, openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInEmbedPreRoll:
					{
						blockId = 529169
						siteId = 241871
					}

				case openrtbex.AdCodeTypeTeasers:
					{
						blockId = 529170
						siteId = 241872
					}

				default:
					{
						blockId = 286015
						siteId = 117855
					}
				}
			}
		}
	}

	req.SiteId = &siteId
	req.BlockId = &blockId
	return req
}
