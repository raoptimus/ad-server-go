package openrtbex

type RequestExt struct {
	Device     Device     `json:"device"`
	Os         Os         `json:"os"`
	Browser    Browser    `json:"browser"`
	CodeId     int        `json:"codeId"`
	ZoneId     int        `json:"zoneId"`
	CodeTypeId AdCodeType `json:"codeTypeId"`
	PlayerType PlayerType `json:"playerType"`
	ZoneType   AdZoneType `json:"zoneType"`
	GeoId      int        `json:"geoId"`
	CountryId  int        `json:"countryId"`
	CityId     int        `json:"cityId"`
	Operator   Operator   `json:"operator"`
	//TODO: move to Debug struct
	IsDebug bool `json:"debug"`
	Debug   struct {
		DisableSession  bool
		DisableRotation bool
	}
}
