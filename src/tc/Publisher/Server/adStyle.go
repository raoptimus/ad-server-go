package main

import (
	"net/url"
	"strconv"
)

type AdStyle struct {
	FontSize        string
	FontColor       string
	BorderColor     string
	BackgroundColor string
	TextAlign       string
	BlockCount      int
	Margin          int
	RotateTime      int
	DelayClose      int
	Delay           int
	ShowTime        int
	HiddenTime      int
}

var defaultStyle = &AdStyle{
	FontSize:        "",
	FontColor:       "",
	BorderColor:     "",
	BackgroundColor: "",
	TextAlign:       "",
	BlockCount:      1,
	Margin:          1,
	RotateTime:      15,
	DelayClose:      10,
	Delay:           0,
	ShowTime:        10,
	HiddenTime:      20,
}

func NewAdStyle(style string) *AdStyle {
	q, err := url.ParseQuery(style)

	if err != nil {
		return defaultStyle
	}

	s := &AdStyle{}
	s.FontSize = s.inOr(q.Get("font-size"), defaultStyle.FontSize)
	s.FontColor = s.inOr(q.Get("font-color"), defaultStyle.FontColor)
	s.BorderColor = s.inOr(q.Get("border-color"), defaultStyle.BorderColor)
	s.BackgroundColor = s.inOr(q.Get("background-color"), defaultStyle.BackgroundColor)
	s.TextAlign = s.inOr(q.Get("text-align"), defaultStyle.TextAlign)
	s.BlockCount = s.atoiOr(q.Get("blockcount"), defaultStyle.BlockCount)
	s.Margin = s.atoiOr(q.Get("margin"), defaultStyle.Margin)
	s.RotateTime = s.atoiOr(q.Get("rotatetime"), defaultStyle.RotateTime)
	s.DelayClose = s.atoiOr(q.Get("delayclose"), defaultStyle.DelayClose)
	s.Delay = s.atoiOr(q.Get("delay"), defaultStyle.Delay)
	s.ShowTime = s.atoiOr(q.Get("showtime"), defaultStyle.ShowTime)
	s.HiddenTime = s.atoiOr(q.Get("hiddentime"), defaultStyle.HiddenTime)

	return s
}

func (s *AdStyle) inOr(in string, def string) string {
	if in == "" {
		return def
	}

	return in
}

func (s *AdStyle) atoiOr(in string, def int) int {
	if in == "" {
		return def
	}

	i, err := strconv.Atoi(in)

	if err != nil {
		return def
	}

	return i
}
