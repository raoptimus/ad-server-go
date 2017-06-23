package openrtbex

import "encoding/gob"

func GobRegisterExt() {
	gob.Register(ImpExt{})
	gob.Register(SiteExt{})
	gob.Register(BidExt{})
	gob.Register(DebugExt{})
	gob.Register(BidStatusExt{})
	gob.Register(RequestExt{})
	gob.Register(ResponseExt{})
}
