package detect

type (
	Region struct {
		Id     int
		NameEn string
		NameRu string
	}
	RegionStorage struct {
		items map[string]*Region
	}
)

func (s *RegionStorage) GetUnknownRegion() *Region {
	return s.items["un"]
}

func (s *RegionStorage) GetRegionByCode(code string) *Region {
	r, ok := s.items[code]

	if ok {
		return r
	}

	return s.GetUnknownRegion()
}

func NewRegionStorage() *RegionStorage {
	c := &RegionStorage{
		items: map[string]*Region{
			"un": &Region{Id: 0, NameEn: "Other", NameRu: "Другие"},
			"ru": &Region{Id: 1, NameEn: "Russian", NameRu: "Россия"},
		},
	}

	cis := &Region{Id: 2, NameEn: "CIS", NameRu: "СНГ"}

	for _, cc := range []string{"ua", "kz", "by", "md", "am", "az", "ge", "kg", "tj", "uz", "tm"} {
		c.items[cc] = cis
	}

	eu := &Region{Id: 3, NameEn: "Europa", NameRu: "Европа"}

	for _, cc := range []string{"de", "fr", "gb", "es", "it", "pt", "be", "nl", "lu", "at", "ch", "dk", "is", "ie", "pl", "cz", "sk", "bg", "hu", "ro", "no", "fi", "se", "ee", "lv", "lt", "gr", "cy", "mt", "ad", "rs", "mk", "al", "ba", "si", "hr", "me"} {
		c.items[cc] = eu
	}

	as := &Region{Id: 4, NameEn: "Asia", NameRu: "Азия"}

	for _, cc := range []string{"af", "bd", "bh", "bn", "bt", "vn", "hk", "il", "in", "id", "jo", "iq", "ir", "ye", "kh", "qa", "cn", "kw", "la", "lb", "my", "mv", "mn", "mm", "np", "ae", "om", "pk", "sa", "kp", "sg", "sy", "th", "tw", "tr", "ph", "lk", "kr", "jp"} {
		c.items[cc] = as
	}

	na := &Region{Id: 5, NameEn: "North America", NameRu: "Северная Америка"}

	for _, cc := range []string{"us", "gt", "hn", "gl", "ca", "cr", "cu", "mx", "ni", "pa"} {
		c.items[cc] = na
	}

	sa := &Region{Id: 6, NameEn: "South America", NameRu: "Южная Америка"}

	for _, cc := range []string{"ai", "ag", "ar", "aw", "bs", "bb", "bz", "bo", "br", "vg", "ve", "vi", "ht", "gy", "gp", "gd", "do", "dm", "ky", "co", "mq", "ms", "an", "py", "pe", "pr", "vc", "kn", "lc", "sr", "tt", "uy", "fk", "gf", "cl", "ec", "sv", "jm"} {
		c.items[cc] = sa
	}

	au := &Region{Id: 7, NameEn: "Australia and Oceania", NameRu: "Австралия и Океания"}

	for _, cc := range []string{"au", "as", "vu", "gu", "ki", "mh", "fm", "nz", "nc", "ck", "pw", "pg", "ws", "sb", "to", "fj", "pf"} {
		c.items[cc] = au
	}

	af := &Region{Id: 8, NameEn: "Africa", NameRu: "Африка"}

	for _, cc := range []string{"dz", "ao", "bj", "bw", "bf", "bi", "ga", "gm", "gh", "gw", "gn", "cd", "dj", "eg", "zm", "zw", "cv", "cm", "ke", "km", "cg", "ci", "ls", "lr", "ly", "mu", "mr", "mg", "yt", "mw", "ml", "ma", "mz", "na", "ng", "ne", "re", "rw", "st", "sz", "sc", "sn", "so", "sd", "sl", "tz", "tg", "tn", "ug", "cf", "td", "gq", "er", "et", "za"} {
		c.items[cc] = af
	}

	return c
}
