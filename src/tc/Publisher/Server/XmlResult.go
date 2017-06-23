package main

import (
	"fmt"
	"strings"
	"tc/openrtbex"
)

type XmlResult struct {
	hc *HttpContext
}

func NewXmlResult(hc *HttpContext) *XmlResult {
	return &XmlResult{
		hc: hc,
	}
}

func (s *XmlResult) compile() string {
	hc := s.hc
	city := hc.City.Title(hc.LangIsRu)

	feed := `<?xml version="1.1" encoding="UTF-8" ?>`
	feed += "<result>"
	adZone := hc.AdCode.AdZone
	style := adZone.GetStyle()
	isRus := strings.Contains(*hc.Req.Device.Language, "ru")
	siteExt := hc.Req.Site.Ext["siteExt"].(openrtbex.SiteExt)

	switch hc.AdCode.TypeId {
	case openrtbex.AdCodeTypeInVideoPauseRoll:
		{
			feed += `<thumbs_ad stop="1" pause="1" show="1">`
			feed += "<backgroundcolor>" + style.BackgroundColor + "</backgroundcolor>"
			feed += "<bordercolor>" + style.BorderColor + "</bordercolor>"
			feed += "<fontcolor>" + style.FontColor + "</fontcolor>"
			if isRus {
				feed += "<header>Реклама</header>"
			} else {
				feed += "<header>Advertising</header>"
			}
		}

	case openrtbex.AdCodeTypeInVideoOverlay:
		{
			feed += `<imagetext_ad show="1" hidden_time="20" show_time="10" rotate_time="15">`
			feed += "<backgroundcolor>" + style.BackgroundColor + "</backgroundcolor>"
			feed += "<bordercolor>" + style.BorderColor + "</bordercolor>"
			feed += "<fontcolor>" + style.FontColor + "</fontcolor>"
			feed += "<target>_blank</target>"
			feed += "<allow_flash>1</allow_flash>"

			if hc.AdCode.IsShowWatermark {
				feed += "<adurl>" + fmt.Sprintf(REF_URL, siteExt.Id) + "</adurl>"
				feed += `<adtitle><![CDATA[<font color="` + style.FontColor + `">ads by 12traffic.com</font>]]></adtitle>`
			}
		}
	}

	for _, bid := range hc.Win {
		bidExt := bid.Ext["bidExt"].(openrtbex.BidExt)
		url := bidExt.Url //todo обернуть и экранировать
		title := bidExt.Url
		title = strings.Replace(title, "#city#", "%city%", -1)
		title = strings.Replace(title, "%city%", city, -1)

		switch hc.AdCode.TypeId {
		case openrtbex.AdCodeTypeInVideoPauseRoll:
			{
				feed += `<thumb width="160" height="160">`
				feed += `<title><![CDATA[<font color="` + style.FontColor + `" size="` + style.FontSize + `"><b>` + title + `</b></font>]]></title>`
				feed += "<src>" + bidExt.Image + "</src>"
				feed += "<url>" + url + "</url>"
				feed += "</thumb>"
			}
		case openrtbex.AdCodeTypeInVideoOverlay:
			{
				feed += `<imagetext>`
				feed += `<title><![CDATA[<u><font color="` + style.FontColor + `"><b>` + title + `</b></font></u>]]></title>`
				feed += `<description><![CDATA[<u><font color="` + style.FontColor + `"><b>` + title + `</b></font></u>]]></description>`
				feed += "<src>" + bidExt.Image + "</src>"
				feed += "<url>" + url + "</url>"
				feed += "</thumb>"
			}

		}
	}

	feed += "</result>"
	return feed
}
