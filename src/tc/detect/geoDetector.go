package detect

import (
	"github.com/abh/geoip"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type GeoDetector struct {
	db      *geoip.GeoIP
	country *CountryStorage
	city    *CityStorage
	region  *RegionStorage
}

const (
	GeoCityDbFileName string = "/geo/GeoIPCity.dat"
	GeoOrgDbFileName  string = "/geo/GeoIPOrg.dat"
)

func NewGeoDetector() *GeoDetector {
	g := &GeoDetector{}
	g.country = NewCountryStorage()
	g.city = NewCityStorage(g)
	g.region = NewRegionStorage()

	curDir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		panic(err)
	}

	r, err := geoip.Open(filepath.Dir(curDir) + GeoCityDbFileName)

	if err != nil {
		log.Println(err)
	}

	g.db = r
	return g
}

func (s *GeoDetector) Detect(ip string) (geoId int, cc *Country, c *City, r *Region) {
	l := s.db.GetRecord(ip)

	if l == nil {
		return 0, s.country.GetUnknown(), s.city.GetUnknownCity(), s.region.GetUnknownRegion()
	}

	//	log.Printf("%v", l)

	cc = s.country.GetByCode(strings.ToLower(l.CountryCode))
	r = s.region.GetRegionByCode(cc.Code)
	c = s.city.GetCityOrUnknown(cc.Id, l.City)
	geoId = 0

	for gId, geo := range GeoGroups {
		if geo.CountryId == cc.Id {
			if geo.CityId == c.Id || geo.CityId == 0 {
				geoId = gId
				break
			}
		}
	}

	return
}
