package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"tc/openrtbex"
)

type FeedResult struct {
	hc *HttpContext
}

func NewFeedResult(hc *HttpContext) *FeedResult {
	return &FeedResult{
		hc: hc,
	}
}

func (s *FeedResult) write() {
	hc := s.hc
	var out string

	if hc.IsDebug {
		hc.W.Header().Add("Content-Type", "text/html")
		out = NewDebugResult(hc).compile()
		io.WriteString(hc.W, out)
		return
	}

	q := hc.R.URL.Query()
	isCrypt := hc.AdCode.IsCrypt()
	player := q.Get("p")

	if hc.AdZone.TypeId == openrtbex.AdZoneTypePopunder {
		NewPopResult(hc.Win).redirect(hc.W)
		return
	}

	if (q.Get("j") == "" && q.Get("action") == "") &&
		(hc.AdCode.TypeId == openrtbex.AdCodeTypeInVideoPauseRoll || hc.AdCode.TypeId == openrtbex.AdCodeTypeInVideoOverlay) &&
		!strings.Contains(player, "12traff") {
		isCrypt = true
		out = NewXmlResult(hc).compile()
	} else {
		out = NewJsonResult(hc).compile()
		isCrypt = !strings.Contains(player, "12traff") && q.Get("ra") == "" && isCrypt
	}

	if out == "" {
		hc.ErrCode = http.StatusNoContent
		hc.WriteError(errors.New("Feed return empty"), "")
		return
	}

	if isCrypt {
		hc.W.Header().Add("Content-Type", "binary/octet-stream;charset=utf-8")
		out = hc.Rc4.crypt(out)
	} else {
		hc.W.Header().Add("Content-Type", "text/javascript;charset=utf-8")
	}

	io.WriteString(hc.W, out)
}
