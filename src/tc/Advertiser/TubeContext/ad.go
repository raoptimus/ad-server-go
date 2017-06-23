package main

import (
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"tc/openrtbex"
)

type (
	Ad struct {
		Id          int        `bson:"_id"`
		CampaignId  int        `bson:"CampaignId"`
		Title       string     `bson:"Title"`
		Description string     `bson:"Description"`
		URL         string     `bson:"URL"`
		UserId      int        `bson:"UserId"`
		IsWebmaster bool       `bson:"IsWebmaster"`
		CityId      int        `bson:"CityId"`
		Images      []*AdImage `bson:"Images"`

		Deleted    bool        `json:"-"`
		Campaign   *Campaign   `json:"-"`
		CtrStorage *CtrStorage `json:"-"`
	}
	AdImage struct {
		Id       bson.ObjectId `bson:"_id"`
		Width    int           `bson:"Width"`
		Height   int           `bson:"Height"`
		Ext      string        `bson:"Ext"`
		Original bool          `bson:"Original"`
	}
)

var imageUrls = []string{
	"http://it....com",
	"http://th.....com",
	"http://im.....com",
}

func (ad *Ad) FormatAllowed(typeId openrtbex.AdCodeType, newPlayer bool) bool {
	if newPlayer {
		return ad.newPlayerFormatAllowed(typeId)
	}
	return ad.oldPlayerFormatAllowed(typeId)
}

func (ad *Ad) Filter(req *openrtb.Request) error {
	reqExt := req.Ext["requestExt"].(openrtbex.RequestExt)
	isNewPlayer := false
	reason := ""

	switch {
	case ad.Deleted:
		{
			reason = "Deleted"
		}
	case ad.CityId > 0 && ad.CityId != reqExt.CityId:
		{
			reason = fmt.Sprintf("City %d", ad.CityId)
		}
	case !ad.FormatAllowed(reqExt.CodeTypeId, isNewPlayer):
		{
			reason = fmt.Sprintf("No images for typeId %s (%d) (newPlayer: %v): (%v)",
				reqExt.CodeTypeId.ToString(), reqExt.CodeTypeId.ToInt(), isNewPlayer, ad.Images)
		}
	default:
		return nil
	}

	return errors.New(fmt.Sprintf("Ad %d is denied because: %+v", ad.Id, reason))
}

func (ad *Ad) ImageUrl(w, h int) string {
	if w+h > 0 {
		for _, img := range ad.Images {
			if img.Original {
				continue
			}

			if w%h == img.Width%img.Height {
				return ad.compileImgUrl(img)
			}
		}
	}

	return ad.randImgPath() + "/404.jpg"
}

//-----------------------------------------------------------------------------------

func (ad *Ad) randImgPath() string {
	return imageUrls[rand.Intn(len(imageUrls))]
}

func (ad *Ad) compileImgUrl(img *AdImage) string {
	return fmt.Sprintf("%s/%x.%s?AdId=%d", ad.randImgPath(), string(img.Id), img.Ext, ad.Id)
}

func (ad *Ad) newPlayerFormatAllowed(typeId openrtbex.AdCodeType) bool {
	switch typeId {
	case openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInHtml5VideoPauseRoll:
		return ad.hasBigSquareImage() || ad.hasImageWithSize(300, 250)
	case openrtbex.AdCodeTypeTeasers:
		return ad.hasBigSquareImage()
	case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInEmbedOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
		return ad.hasSquareImage()
	case openrtbex.AdCodeTypeInVideoPostRoll, openrtbex.AdCodeTypeInVideoPreRoll, openrtbex.AdCodeTypeBanners300x250, openrtbex.AdCodeTypeInEmbedPreRoll, openrtbex.AdCodeTypeMobileBanners300x250:
		return ad.hasImageWithSize(300, 250) || ad.hasImageWithSize(250, 250)
	case openrtbex.AdCodeTypeMobileBanners300x100:
		return ad.hasImageWithSize(300, 100)
	case openrtbex.AdCodeTypeMobileBanners300x50:
		return ad.hasImageWithSize(300, 50)
	}
	return false
}

func (ad *Ad) oldPlayerFormatAllowed(typeId openrtbex.AdCodeType) bool {
	switch typeId {
	case openrtbex.AdCodeTypeInVideoPauseRoll, openrtbex.AdCodeTypeInHtml5VideoPauseRoll:
		return ad.hasBigSquareImage()
	case openrtbex.AdCodeTypeTeasers:
		return ad.hasBigSquareImage()
	case openrtbex.AdCodeTypeInVideoOverlay, openrtbex.AdCodeTypeInEmbedOverlay, openrtbex.AdCodeTypeInHtml5VideoOverlay:
		return ad.hasSquareImage() || ad.hasImageWithHeight(80)
	case openrtbex.AdCodeTypeInVideoPostRoll, openrtbex.AdCodeTypeInVideoPreRoll, openrtbex.AdCodeTypeBanners300x250, openrtbex.AdCodeTypeInEmbedPreRoll, openrtbex.AdCodeTypeMobileBanners300x250:
		return ad.hasImageWithSize(300, 250)
	case openrtbex.AdCodeTypeMobileBanners300x100:
		return ad.hasImageWithSize(300, 100)
	case openrtbex.AdCodeTypeMobileBanners300x50:
		return ad.hasImageWithSize(300, 50)
	}

	return false
}

func (ad *Ad) hasSquareImage() bool {
	for _, img := range ad.Images {
		if img.Width == img.Height {
			return true
		}
	}
	return false
}

func (ad *Ad) hasBigSquareImage() bool {
	for _, img := range ad.Images {
		if img.Width == img.Height && img.Width > 80 {
			return true
		}
	}
	return false
}

func (ad *Ad) hasImageWithSize(w, h int) bool {
	for _, img := range ad.Images {
		if img.Width == w && img.Height == h {
			return true
		}
	}
	return false
}

func (ad *Ad) hasImageWithHeight(h int) bool {
	for _, img := range ad.Images {
		if img.Height == h {
			return true
		}
	}
	return false
}
