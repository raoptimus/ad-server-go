package openrtbex

type PaymentType int

const (
	PaymentTypeCPC PaymentType = 0
	PaymentTypeCPM PaymentType = 1
)

func (pt PaymentType) ToInt() int {
	return int(pt)
}

type Cur string

const (
	CurRub  Cur = "RUB"
	CurUsd  Cur = "USD"
	CurEuro Cur = "EURO"
)

type ConfirmStatus int

const (
	ConfirmStatusWin   ConfirmStatus = 1
	ConfirmStatusLoss  ConfirmStatus = 2
	ConfirmStatusError ConfirmStatus = 3
)

type ConfirmType int

const (
	ConfirmTypeClick ConfirmType = 0
	ConfirmTypeShow  ConfirmType = 1
)

type LossReason string

const (
	LossReasonPrice            LossReason = "price"
	LossReasonDisqualification LossReason = "disqualification"
)

type ErrorReason string

const (
	//An error occurred while establishing the connection.
	ErrorReasonConnection ErrorReason = "ERROR_CONNECT"
	//An empty response was returned.
	ErrorReasonEmptyResponse ErrorReason = "ERROR_EMPTY"
	//There was an error resolving the endpoint domain.
	ErrorReasonUnResolvedDomain ErrorReason = "ERROR_NXDOMAIN"
	//An error occurred while parsing the response
	ErrorReasonUnParsedResponse ErrorReason = "ERROR_PARSE"
	//An error occurred while receiving the response.
	ErrorReasonUnReceivedResponse ErrorReason = "ERROR_RECEIVE"
	//Some other unknown error occurred.
	ErrorReasonUnknown ErrorReason = "ERROR_UNKNOWN"
	//The original bid was not received (it timed out).
	ErrorReasonTimeout ErrorReason = "TIMEOUT"
)

type Device int

const (
	DeviceDesktop Device = 1
	DeviceMobile  Device = 2
	DeviceTablet  Device = 3
)

func (s Device) ToInt() int {
	return int(s)
}

func (s Device) ToString() string {
	switch s {
	case DeviceMobile:
		return "Mobile"
	case DeviceTablet:
		return "Tablet"
	default:
		return "Desktop"
	}
}

func (s Device) ToOpenRTB() *int {
	var d int
	switch s {
	case DeviceMobile:
		d = 4
	case DeviceTablet:
		d = 5
	default:
		d = 2
	}

	return &d
}

type Operator int

const (
	OperatorUnknown Operator = 1
	OperatorMTS     Operator = 2
	OperatorBeeline Operator = 3
	OperatorMegafon Operator = 4
	OperatorTele2   Operator = 5
)

func (s Operator) ToInt() int {
	return int(s)
}

func (s Operator) ToString() string {
	switch s {
	case OperatorMTS:
		return "MTS"
	case OperatorBeeline:
		return "Beeline"
	case OperatorMegafon:
		return "Megafon"
	case OperatorTele2:
		return "Tele2"
	default:
		return "Unknown"
	}
}

type Browser int

const (
	BrowserUnknown     Browser = 0
	BrowserIE          Browser = 1
	BrowserSafari      Browser = 2
	BrowserChrome      Browser = 3
	BrowserFirefox     Browser = 4
	BrowserOpera       Browser = 5
	BrowserOperaMobile Browser = 6
	BrowserOperaMini   Browser = 7
)

func (s Browser) ToInt() int {
	return int(s)
}

func (s Browser) ToString() string {
	switch s {
	case BrowserIE:
		return "IE"
	case BrowserSafari:
		return "Safari"
	case BrowserChrome:
		return "Chrome"
	case BrowserFirefox:
		return "Firefox"
	case BrowserOpera:
		return "Opera"
	case BrowserOperaMini:
		return "OperaMini"
	case BrowserOperaMobile:
		return "OperaMobile"
	default:
		return "Unknown"
	}
}

type Os int

const (
	OsUnknown    Os = 0
	OsIOs        Os = 1
	OsAndroid    Os = 2
	OsWindows    Os = 3
	OsSymbian    Os = 4
	OsBlackBerry Os = 5
	OsMacintosh  Os = 6
	OsLinux      Os = 7
)

func (s Os) ToInt() int {
	return int(s)
}

func (s Os) ToString() string {
	switch s {
	case OsIOs:
		return "IOs"
	case OsAndroid:
		return "Android"
	case OsWindows:
		return "Windows"
	case OsSymbian:
		return "Symbian"
	case OsBlackBerry:
		return "BlackBerry"
	case OsMacintosh:
		return "Macintosh"
	case OsLinux:
		return "Linux"
	default:
		return "Unknown"
	}
}

type AdCodeTypeNew int

const (
	AdCodeTypeNewInVideoPauseRoll AdCodeTypeNew = 0
	AdCodeTypeNewInVideoPreRoll   AdCodeTypeNew = 1
	AdCodeTypeNewInVideoPostRoll  AdCodeTypeNew = 2
	AdCodeTypeNewInVideoOverlay   AdCodeTypeNew = 3

	AdCodeTypeNewTeasers  AdCodeTypeNew = 4
	AdCodeTypeNewBanners  AdCodeTypeNew = 5
	AdCodeTypeNewPopunder AdCodeTypeNew = 6
)

type AdCodeType int

const (
	AdCodeTypeInVideoPauseRoll      AdCodeType = 0
	AdCodeTypeTeasers               AdCodeType = 1
	AdCodeTypeInVideoOverlay        AdCodeType = 2
	AdCodeTypeInVideoPostRoll       AdCodeType = 4
	AdCodeTypeInVideoPreRoll        AdCodeType = 5
	AdCodeTypeInEmbedOverlay        AdCodeType = 6
	AdCodeTypeInHtml5VideoPauseRoll AdCodeType = 7
	AdCodeTypeInHtml5VideoOverlay   AdCodeType = 8
	AdCodeTypeBanners300x250        AdCodeType = 9
	AdCodeTypeInEmbedPreRoll        AdCodeType = 10
	AdCodeTypePopunder              AdCodeType = 11
	AdCodeTypeMobileBanners300x250  AdCodeType = 20
	AdCodeTypeMobileBanners300x100  AdCodeType = 21
	//deprecated
	AdCodeTypeMobileBanners300x50 AdCodeType = 22
	AdCodeTypeMobilePopunder      AdCodeType = 23
)

func (s AdCodeType) ToInt() int {
	return int(s)
}

func (s AdCodeType) AdCodeTypeNew() AdCodeTypeNew {
	switch s {
	case AdCodeTypeInVideoPauseRoll, AdCodeTypeInHtml5VideoPauseRoll:
		return AdCodeTypeNewInVideoPauseRoll
	case AdCodeTypeInVideoOverlay, AdCodeTypeInEmbedOverlay, AdCodeTypeInHtml5VideoOverlay:
		return AdCodeTypeNewInVideoOverlay
	case AdCodeTypeInVideoPostRoll:
		return AdCodeTypeNewInVideoPostRoll
	case AdCodeTypeInVideoPreRoll, AdCodeTypeInEmbedPreRoll:
		return AdCodeTypeNewInVideoPreRoll
	case AdCodeTypeBanners300x250, AdCodeTypeMobileBanners300x250, AdCodeTypeMobileBanners300x100, AdCodeTypeMobileBanners300x50:
		return AdCodeTypeNewBanners
	case AdCodeTypePopunder, AdCodeTypeMobilePopunder:
		return AdCodeTypeNewPopunder
	default:
		return AdCodeTypeNewTeasers
	}
}

func (s AdCodeType) ToString() string {
	switch s {
	case AdCodeTypeInVideoPauseRoll:
		return "InVideoPauseRoll"
	case AdCodeTypeTeasers:
		return "Teasers"
	case AdCodeTypeInVideoOverlay:
		return "InVideoOverlay"
	case AdCodeTypeInVideoPostRoll:
		return "InVideoPostRoll"
	case AdCodeTypeInVideoPreRoll:
		return "InVideoPreRoll"
	case AdCodeTypeInEmbedOverlay:
		return "InEmbedOverlay"
	case AdCodeTypeInHtml5VideoPauseRoll:
		return "InHtml5VideoPauseRoll"
	case AdCodeTypeInHtml5VideoOverlay:
		return "InHtml5VideoOverlay"
	case AdCodeTypeBanners300x250:
		return "Banners300x250"
	case AdCodeTypeInEmbedPreRoll:
		return "InEmbedPreRoll"
	case AdCodeTypePopunder:
		return "Popunder"
	case AdCodeTypeMobileBanners300x250:
		return "MobileBanners300x250"
	case AdCodeTypeMobileBanners300x100:
		return "MobileBanners300x100"
	case AdCodeTypeMobileBanners300x50:
		return "MobileBanners300x50"
	case AdCodeTypeMobilePopunder:
		return "MobilePopunder"
	default:
		return "Unknown"
	}
}

type PlayerType int

const (
	PlayerTypeTubeContext PlayerType = 1
	PlayerTypeJW          PlayerType = 2
	PlayerTypeKernel      PlayerType = 3
	PlayerTypeUppod       PlayerType = 5
	PlayerTypeEmbed       PlayerType = 7
	PlayerTypeWp          PlayerType = 8
	PlayerTypeDle         PlayerType = 9
	PlayerType12Traffic   PlayerType = 10
)

func (s PlayerType) ToInt() int {
	return int(s)
}

func (s PlayerType) ToString() string {
	switch s {
	case PlayerTypeTubeContext:
		return "TubeContext"
	case PlayerTypeJW:
		return "JW"
	case PlayerTypeKernel:
		return "Kernel"
	case PlayerTypeUppod:
		return "Uppod"
	case PlayerTypeEmbed:
		return "Embed"
	case PlayerTypeWp:
		return "Wp"
	case PlayerTypeDle:
		return "Dle"
	case PlayerType12Traffic:
		return "12Traffic"
	default:
		return "Unknown"
	}
}

type AdZoneType int

const (
	AdZoneTypeInVideo        AdZoneType = 1
	AdZoneTypeBanners        AdZoneType = 2
	AdZoneTypeTeasers        AdZoneType = 3
	AdZoneTypeMobileBanners  AdZoneType = 4
	AdZoneTypePopunder       AdZoneType = 10
	AdZoneTypeMobilePopunder AdZoneType = 11
)

func (s AdZoneType) ToInt() int {
	return int(s)
}

func (s AdZoneType) ToString() string {
	switch s {
	case AdZoneTypeInVideo:
		return "InVideo"
	case AdZoneTypeBanners, AdZoneTypeMobileBanners:
		return "Banners"
	case AdZoneTypeTeasers:
		return "Teasers"
	case AdZoneTypePopunder, AdZoneTypeMobilePopunder:
		return "Popunder"
	default:
		return ""
	}
}

type BidType int

const (
	BidTypePay  BidType = 0
	BidTypeFree BidType = 1
	BidTypeWm   BidType = 2
)

type CampaignType int

const (
	CampaignTypeDefault  CampaignType = 0
	CampaignTypePopunder CampaignType = 1
	CampaignTypeInVideo  CampaignType = 2
	CampaignTypeBanners  CampaignType = 3
)

func (ct CampaignType) ToInt() int {
	return int(ct)
}
