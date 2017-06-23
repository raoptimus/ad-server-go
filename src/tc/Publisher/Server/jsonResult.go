package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"tc/openrtbex"
)

const REF_URL = "http://aff....com/?u=%d"
const REF_TITLE = "ads by ...com"
const OVERLAY_SEC_AFTER_SHOW = 20

type (
	JsonResult struct {
		hc *HttpContext
	}
)

func NewJsonResult(hc *HttpContext) *JsonResult {
	return &JsonResult{
		hc: hc,
	}
}

func (s *JsonResult) compile() string {
	hc := s.hc
	city := hc.City.Title(hc.LangIsRu)

	feed := make(map[string]interface{})
	adZone := hc.AdCode.AdZone
	style := adZone.GetStyle()
	siteExt := hc.Req.Site.Ext["siteExt"].(openrtbex.SiteExt)

	switch hc.AdCode.TypeId {
	case
		openrtbex.AdCodeTypeInVideoPauseRoll,
		openrtbex.AdCodeTypeInHtml5VideoPauseRoll:
		{
			if hc.LangIsRu {
				feed["header"] = "Реклама"
			} else {
				feed["header"] = "Advertising"
			}

			feed["show"] = 1
			feed["pause"] = 1
			feed["stop"] = 1
			feed["backgroundColor"] = style.BackgroundColor
			feed["borderColor"] = style.BorderColor
			feed["fontColor"] = style.FontColor
			feed["type"] = 1
		}

	case
		openrtbex.AdCodeTypeTeasers,
		openrtbex.AdCodeTypeBanners300x250,
		openrtbex.AdCodeTypeMobileBanners300x250,
		openrtbex.AdCodeTypeMobileBanners300x100,
		openrtbex.AdCodeTypeMobileBanners300x50:
		{
			feed["background-color"] = style.BackgroundColor
			feed["border-color"] = style.BorderColor
			feed["font-color"] = style.FontColor
			feed["color"] = style.FontColor
			feed["font-size"] = style.FontSize
			feed["width"] = 250
			feed["height"] = 250
			feed["margin"] = style.Margin
		}
	case
		openrtbex.AdCodeTypeInVideoOverlay,
		openrtbex.AdCodeTypeInEmbedOverlay,
		openrtbex.AdCodeTypeInHtml5VideoOverlay:
		{
			feed["show"] = 1
			feed["hiddenTime"] = 20
			feed["showTime"] = OVERLAY_SEC_AFTER_SHOW
			feed["rotateTime"] = 15
			feed["backgroundColor"] = style.BackgroundColor
			feed["borderColor"] = style.BorderColor
			feed["fontColor"] = style.FontColor
			feed["target"] = "_blank"
			feed["type"] = 3
			feed["clickClose"] = 1
			feed["delayClose"] = 0
			feed["blockId"] = hc.BlockId
		}
	case
		openrtbex.AdCodeTypeInVideoPreRoll,
		openrtbex.AdCodeTypeInEmbedPreRoll:
		{
			feed["show"] = 1
			feed["delay"] = 0
			feed["autoClose"] = 1
			feed["type"] = 2
			feed["blockId"] = hc.BlockId
		}
	}

	if hc.AdCode.IsShowWatermark {
		feed["adUrl"] = fmt.Sprintf(REF_URL, siteExt.Id)
		feed["adTarget"] = "_blank"
		feed["adTitle"] = "ads by 12traffic.com"
		feed["adFontColor"] = style.FontColor
	}

	switch hc.AdCode.TypeId {
	case
		openrtbex.AdCodeTypeInVideoPreRoll,
		openrtbex.AdCodeTypeInVideoPostRoll,
		openrtbex.AdCodeTypeInEmbedPreRoll:
		{
			if hc.LangIsRu {
				feed["Advertising"] = "Реклама"
				feed["header"] = "Реклама"
				feed["CloseAndPlay"] = "Закрыть рекламу"
				feed["SkipAd"] = "Пропустить рекламу"
			} else {
				feed["Advertising"] = "Advertising"
				feed["header"] = "Advertising"
				feed["CloseAndPlay"] = "Close & Play"
				feed["SkipAd"] = "Skip Ad"
			}
		}
	}

	adList := make([]map[string]interface{}, 0)

	for _, bid := range hc.Win {
		bidExt := bid.Ext["bidExt"].(openrtbex.BidExt)
		ad := make(map[string]interface{})
		u, err := bidExt.ClickUrl(hc.R, hc.Debug, hc.Ip, hc.Ua, hc.AdCode.Id)
		if err != nil {
			return ""
		}
		title := bidExt.Title
		title = strings.Replace(title, "#city#", "%city%", -1)
		title = strings.Replace(title, "%city%", city, -1)
		tags := map[string]string{
			"title":     "title",
			"desc":      "description",
			"url":       "url",
			"src":       "src",
			"width":     "width",
			"height":    "height",
			"textAlign": "text-align",
			"margin":    "margin",
		}

		switch hc.AdCode.TypeId {
		case
			openrtbex.AdCodeTypeTeasers,
			openrtbex.AdCodeTypeBanners300x250,
			openrtbex.AdCodeTypeMobileBanners300x250,
			openrtbex.AdCodeTypeMobileBanners300x100,
			openrtbex.AdCodeTypeMobileBanners300x50:
			{
				tags["src"] = "imgurl"
				tags["desc"] = "desc"
			}
		case
			openrtbex.AdCodeTypeInVideoPreRoll,
			openrtbex.AdCodeTypeInVideoPostRoll,
			openrtbex.AdCodeTypeInEmbedPreRoll:
			{
				tags["desc"] = "desc"
			}
		}

		ad[tags["title"]] = title
		ad[tags["desc"]] = title
		ad[tags["url"]] = u.String()
		ad[tags["src"]] = bidExt.Image
		ad[tags["width"]] = bidExt.Width
		ad[tags["height"]] = bidExt.Height
		ad[tags["textAlign"]] = style.TextAlign
		ad[tags["margin"]] = style.Margin

		adList = append(adList, ad)
	}

	feed["ads"] = adList

	b, err := json.Marshal(feed)
	if err != nil {
		return ""
	}
	return string(b)
}
