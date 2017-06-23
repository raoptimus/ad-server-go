package detect

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"tc/data"
)

type City struct {
	Id        int      `bson:"_id"`
	CountryId int      `bson:"CountryId"`
	EngName   []string `bson:"EngName"`
	RusName   string   `bson:"RusName"`
}

type CityStorage struct {
	items map[int]map[string]*City
	geo   *GeoDetector
}

func NewCityStorage(g *GeoDetector) *CityStorage {
	cs := &CityStorage{
		items: make(map[int]map[string]*City),
		geo:   g,
	}
	go cs.load()
	return cs
}

func (s *City) Title(isRus bool) string {
	if isRus {
		return s.RusName
	}

	if len(s.EngName) > 0 {
		return s.EngName[0]
	}

	return ""
}

func (s *CityStorage) load() {
	items := make(map[int]map[string]*City)
	ccList := s.geo.country.GetAll()
	var iter *mgo.Iter

	for _, cc := range ccList {
		cList := make(map[string]*City)

		//get cities for country
		iter = data.DataContext.Cities.Find(bson.M{"CountryId": cc.Id}).Iter()
		c := &City{}

		for iter.Next(c) {
			for _, name := range c.EngName {
				cList[name] = c
			}

			c = &City{}
		}

		if err := iter.Err(); err != nil {
			log.Panicln(err)
		}

		items[cc.Id] = cList
	}

	s.items = items
}

var defCity *City = &City{
	Id:        0,
	CountryId: 0,
	EngName:   []string{""},
	RusName:   "",
}

func (s *CityStorage) GetUnknownCity() *City {
	return defCity
}

func (s *CityStorage) GetCityOrUnknown(countryId int, name string) *City {
	cList, ok := s.items[countryId]

	if ok && cList != nil && len(cList) > 0 {
		c, ok := cList[name]

		if ok {
			return c
		}
	}
	//todo set CountryId
	return s.GetUnknownCity()
}
