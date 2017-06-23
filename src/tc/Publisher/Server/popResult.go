package main

import (
	"fmt"
	"io"
	"net/http"
	"tc/openrtbex"
)

type PopResult struct {
	win BidList
}

func NewPopResult(win BidList) *PopResult {
	return &PopResult{
		win: win,
	}
}

func (s *PopResult) redirect(w http.ResponseWriter) {
	bidExt := s.win[0].Ext["bidExt"].(openrtbex.BidExt)

	out := "<html><header></header><body>" +
		"<script>document.location.href='%s';</script>" +
		"<noscript><meta http-equiv=\"refresh\" content=\"1;url=%s\" /></noscript>" +
		"</body></html>"

	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	io.WriteString(w, fmt.Sprintf(out, bidExt.Url))
}

func (s *PopResult) close(w http.ResponseWriter) {
	out := "<html><header></header><body><script>window.close();</script></body></html>"
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	io.WriteString(w, out)
}
