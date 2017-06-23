package openrtbex

import (
	"math"
	"strconv"
)

type Category int

const (
	CategoryAdult Category = 25003
	CategoryOther Category = 24000
)

func (s Category) ToInt() int {
	return int(s)
}

func (s Category) ToFormatRTB() string {
	var d float64 = float64(s) / 1000.0
	g := math.Floor(d)
	var i float64 = (d - g) * 1000.0
	id := int(math.Floor(i))
	rtb := "IAB" + strconv.Itoa(int(g))

	if id > 0 {
		rtb += "-" + strconv.Itoa(id)
	}

	return rtb
}
