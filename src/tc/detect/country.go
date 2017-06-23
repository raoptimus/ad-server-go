package detect

type (
	Country struct {
		Id     int
		Code   string
		NameEn string
		NameRu string
	}
	CountryStorage struct {
		items map[string]*Country
	}
)

func (s *Country) Title(isRus bool) string {
	if isRus {
		return s.NameRu
	}

	return s.NameEn
}
func (s *CountryStorage) GetUnknown() *Country {
	return s.GetByCode("un")
}

func (s *CountryStorage) GetByCode(code string) *Country {
	cc, ok := s.items[code]

	if ok {
		return cc
	}

	return s.items["un"]
}

func (s *CountryStorage) GetAll() []*Country {
	ccList := make([]*Country, 0)

	for _, cc := range s.items {
		ccList = append(ccList, cc)
	}

	return ccList
}

func NewCountryStorage() *CountryStorage {
	c := &CountryStorage{
		items: map[string]*Country{
			"un": &Country{Id: 0, Code: "un", NameEn: "Other", NameRu: "Любая другая"},
			"ru": &Country{Id: 1, Code: "ru", NameEn: "Russia", NameRu: "Россия"},
			"by": &Country{Id: 11, Code: "by", NameEn: "Belarus", NameRu: "Беларусь"},
			"ua": &Country{Id: 12, Code: "ua", NameEn: "Ukraine", NameRu: "Украина"},
			"kz": &Country{Id: 13, Code: "kz", NameEn: "Kazakhstan", NameRu: "Казахстан"},
			"md": &Country{Id: 14, Code: "md", NameEn: "Moldova, Republic of", NameRu: "Молдова"},
			"uz": &Country{Id: 15, Code: "uz", NameEn: "Uzbekistan", NameRu: "Узбекистан"},
			"ee": &Country{Id: 16, Code: "ee", NameEn: "Estonia", NameRu: "Эстония"},
			"lv": &Country{Id: 17, Code: "lv", NameEn: "Latvia", NameRu: "Латвия"},
			"il": &Country{Id: 18, Code: "il", NameEn: "Israel", NameRu: "Израиль"},
			"cn": &Country{Id: 19, Code: "cn", NameEn: "China", NameRu: "Китай"},
			"no": &Country{Id: 20, Code: "no", NameEn: "Norway", NameRu: "Норвегия"},
			"pl": &Country{Id: 21, Code: "pl", NameEn: "Poland", NameRu: "Польша"},
			"am": &Country{Id: 22, Code: "am", NameEn: "Armenia", NameRu: "Армения"},
			"az": &Country{Id: 23, Code: "az", NameEn: "Azerbaijan", NameRu: "Азербайджан"},
			"ge": &Country{Id: 24, Code: "ge", NameEn: "Georgia", NameRu: "Грузия"},
			"gb": &Country{Id: 25, Code: "gb", NameEn: "United Kingdom", NameRu: "Англия"},
			"bg": &Country{Id: 26, Code: "bg", NameEn: "Bulgaria", NameRu: "Болгария"},
			"us": &Country{Id: 27, Code: "us", NameEn: "United States", NameRu: "США"},
			"de": &Country{Id: 28, Code: "de", NameEn: "Germany", NameRu: "Германия"},
			"nl": &Country{Id: 29, Code: "nl", NameEn: "Netherlands", NameRu: "Нидерланды"},
			"it": &Country{Id: 30, Code: "it", NameEn: "Italy", NameRu: "Италия"},
			"tm": &Country{Id: 31, Code: "tm", NameEn: "Turkmenistan", NameRu: "Туркменистан"},
			"a2": &Country{Id: 32, Code: "a2", NameEn: "Satellite Provider", NameRu: "Спутник"},
			"be": &Country{Id: 33, Code: "be", NameEn: "Belgium", NameRu: "Бельгия"},
			"ma": &Country{Id: 34, Code: "ma", NameEn: "Morocco", NameRu: "Морокко"},
			"tj": &Country{Id: 35, Code: "tj", NameEn: "Tajikistan", NameRu: "Таджикистан"},
			"jo": &Country{Id: 36, Code: "jo", NameEn: "Jordan", NameRu: "Иордания"},
			"dk": &Country{Id: 37, Code: "dk", NameEn: "Denmark", NameRu: "Дания"},
			"kg": &Country{Id: 38, Code: "kg", NameEn: "Kyrgyzstan", NameRu: "Кыргызстан"},
			"lt": &Country{Id: 39, Code: "lt", NameEn: "Lithuania", NameRu: "Литва"},
			"th": &Country{Id: 40, Code: "th", NameEn: "Thailand", NameRu: "Таиланд"},
			"sk": &Country{Id: 41, Code: "sk", NameEn: "Slovakia", NameRu: "Словакия"},
			"ch": &Country{Id: 42, Code: "ch", NameEn: "Switzerland", NameRu: "Швейцария"},
			"es": &Country{Id: 43, Code: "es", NameEn: "Spain", NameRu: "Испания"},
			"gr": &Country{Id: 44, Code: "gr", NameEn: "Greece", NameRu: "Греция"},
			"tr": &Country{Id: 45, Code: "tr", NameEn: "Turkey", NameRu: "Турция"},
			"ie": &Country{Id: 46, Code: "ie", NameEn: "Ireland", NameRu: "Ирландия"},
			"kr": &Country{Id: 47, Code: "kr", NameEn: "Korea, Republic of", NameRu: "Республика Корея"},
			"fr": &Country{Id: 48, Code: "fr", NameEn: "France", NameRu: "Франция"},
			"ly": &Country{Id: 49, Code: "ly", NameEn: "Libyan Arab Jamahiriya", NameRu: "Ливийская Арабская Джамахирия"},
			"jp": &Country{Id: 50, Code: "jp", NameEn: "Japan", NameRu: "Япония"},
			"ps": &Country{Id: 51, Code: "ps", NameEn: "Palestinian Territory", NameRu: "Палестинская территория"},
			"id": &Country{Id: 52, Code: "id", NameEn: "Indonesia", NameRu: "Индонезия"},
			"ca": &Country{Id: 53, Code: "ca", NameEn: "Canada", NameRu: "Канада"},
			"tw": &Country{Id: 54, Code: "tw", NameEn: "Taiwan", NameRu: "Тайван"},
			"cy": &Country{Id: 55, Code: "cy", NameEn: "Cyprus", NameRu: "Кипр"},
			"me": &Country{Id: 56, Code: "me", NameEn: "Montenegro", NameRu: "Черногория"},
			"pt": &Country{Id: 57, Code: "pt", NameEn: "Portugal", NameRu: "Португалия"},
			"rs": &Country{Id: 58, Code: "rs", NameEn: "Serbia", NameRu: "Сербия"},
			"br": &Country{Id: 59, Code: "br", NameEn: "Brazil", NameRu: "Бразилия"},
			"se": &Country{Id: 60, Code: "se", NameEn: "Sweden", NameRu: "Швеция"},
			"mn": &Country{Id: 61, Code: "mn", NameEn: "Mongolia", NameRu: "Монголия"},
			"at": &Country{Id: 62, Code: "at", NameEn: "Austria", NameRu: "Австрия"},
			"nz": &Country{Id: 63, Code: "nz", NameEn: "New Zealand", NameRu: "Новая Зеландия"},
			"cz": &Country{Id: 64, Code: "cz", NameEn: "Czech Republic", NameRu: "Чешская Республика"},
			"hr": &Country{Id: 65, Code: "hr", NameEn: "Croatia", NameRu: "Хорватия"},
			"in": &Country{Id: 66, Code: "in", NameEn: "India", NameRu: "Индия"},
			"bd": &Country{Id: 67, Code: "bd", NameEn: "Bangladesh", NameRu: "Бангладеш"},
			"my": &Country{Id: 68, Code: "my", NameEn: "Malaysia", NameRu: "Малайзия"},
			"pe": &Country{Id: 69, Code: "pe", NameEn: "Peru", NameRu: "Перу"},
			"mx": &Country{Id: 70, Code: "mx", NameEn: "Mexico", NameRu: "Мексика"},
			"fi": &Country{Id: 71, Code: "fi", NameEn: "Finland", NameRu: "Финляндия"},
			"sd": &Country{Id: 72, Code: "sd", NameEn: "Sudan", NameRu: "Судан"},
			"ar": &Country{Id: 73, Code: "ar", NameEn: "Argentina", NameRu: "Аргентина"},
			"ro": &Country{Id: 74, Code: "ro", NameEn: "Romania", NameRu: "Румыния"},
			"cr": &Country{Id: 75, Code: "cr", NameEn: "Costa Rica", NameRu: "Коста-Рика"},
			"ye": &Country{Id: 76, Code: "ye", NameEn: "Yemen", NameRu: "Йемен"},
			"au": &Country{Id: 77, Code: "au", NameEn: "Australia", NameRu: "Австралия"},
			"ir": &Country{Id: 78, Code: "ir", NameEn: "Iran, Islamic Republic of", NameRu: "Иран"},
			"lu": &Country{Id: 79, Code: "lu", NameEn: "Luxembourg", NameRu: "Люксембург"},
			"eu": &Country{Id: 80, Code: "eu", NameEn: "Europe", NameRu: "Европа"},
			"vn": &Country{Id: 81, Code: "vn", NameEn: "Vietnam", NameRu: "Вьетнам"},
			"si": &Country{Id: 82, Code: "si", NameEn: "Slovenia", NameRu: "Словения"},
			"pa": &Country{Id: 83, Code: "pa", NameEn: "Panama", NameRu: "Панама"},
			"cl": &Country{Id: 84, Code: "cl", NameEn: "Chile", NameRu: "Чили"},
			"lk": &Country{Id: 85, Code: "lk", NameEn: "Sri Lanka", NameRu: "Шри-Ланка"},
			"ba": &Country{Id: 86, Code: "ba", NameEn: "Bosnia and Herzegovina", NameRu: "Босния и Герцеговина"},
			"et": &Country{Id: 87, Code: "et", NameEn: "Ethiopia", NameRu: "Эфиопия"},
			"mk": &Country{Id: 88, Code: "mk", NameEn: "Macedonia", NameRu: "Македония"},
			"ph": &Country{Id: 89, Code: "ph", NameEn: "Philippines", NameRu: "Филиппины"},
			"hu": &Country{Id: 90, Code: "hu", NameEn: "Hungary", NameRu: "Венгрия"},
			"eg": &Country{Id: 91, Code: "eg", NameEn: "Egypt", NameRu: "Египет"},
			"za": &Country{Id: 92, Code: "za", NameEn: "South Africa", NameRu: "Южная Африка"},
			"dz": &Country{Id: 93, Code: "dz", NameEn: "Algeria", NameRu: "Алжир"},
			"ni": &Country{Id: 94, Code: "ni", NameEn: "Nicaragua", NameRu: "Никарагуа"},
			"sg": &Country{Id: 95, Code: "sg", NameEn: "Singapore", NameRu: "Сингапур"},
			"ve": &Country{Id: 96, Code: "ve", NameEn: "Venezuela", NameRu: "Венесуэла"},
			"iq": &Country{Id: 97, Code: "iq", NameEn: "Iraq", NameRu: "Ирак"},
			"do": &Country{Id: 98, Code: "do", NameEn: "Dominican Republic", NameRu: "Доминиканская Республика"},
			"lb": &Country{Id: 99, Code: "lb", NameEn: "Lebanon", NameRu: "Ливан"},
			"pk": &Country{Id: 100, Code: "pk", NameEn: "Pakistan", NameRu: "Пакистан"},
			"ng": &Country{Id: 101, Code: "ng", NameEn: "Nigeria", NameRu: "Нигерия"},
			"ao": &Country{Id: 102, Code: "ao", NameEn: "Angola", NameRu: "Ангола"},
			"la": &Country{Id: 103, Code: "la", NameEn: "Lao People's Democratic Republic", NameRu: "Лаос"},
			"mm": &Country{Id: 104, Code: "mm", NameEn: "Myanmar", NameRu: "Мьянмы"},
			"mt": &Country{Id: 105, Code: "mt", NameEn: "Malta", NameRu: "Мальта"},
			"ec": &Country{Id: 106, Code: "ec", NameEn: "Ecuador", NameRu: "Эквадор"},
			"sa": &Country{Id: 107, Code: "sa", NameEn: "Saudi Arabia", NameRu: "Саудовская Аравия"},
			"hk": &Country{Id: 108, Code: "hk", NameEn: "Hong Kong", NameRu: "Гонконг"},
			"a1": &Country{Id: 109, Code: "a1", NameEn: "Anonymous Proxy", NameRu: "Прокси"},
			"ke": &Country{Id: 110, Code: "ke", NameEn: "Kenya", NameRu: "Кения"},
			"af": &Country{Id: 111, Code: "af", NameEn: "Afghanistan", NameRu: "Афганистан"},
			"sy": &Country{Id: 112, Code: "sy", NameEn: "Syrian Arab Republic", NameRu: "Сирийская Арабская Республика"},
			"co": &Country{Id: 113, Code: "co", NameEn: "Colombia", NameRu: "Колумбия"},
			"bs": &Country{Id: 114, Code: "bs", NameEn: "Bahamas", NameRu: "Багамские острова"},
			"na": &Country{Id: 115, Code: "na", NameEn: "Namibia", NameRu: "Намибия"},
			"tn": &Country{Id: 116, Code: "tn", NameEn: "Tunisia", NameRu: "Tunisia"},
			"bo": &Country{Id: 117, Code: "bo", NameEn: "Bolivia", NameRu: "Боливия"},
			"lr": &Country{Id: 118, Code: "lr", NameEn: "Liberia", NameRu: "Либерия"},
			"tz": &Country{Id: 119, Code: "tz", NameEn: "Tanzania, United Republic of", NameRu: "Танзания"},
			"al": &Country{Id: 120, Code: "al", NameEn: "Albania", NameRu: "Албания"},
			"mc": &Country{Id: 121, Code: "mc", NameEn: "Monaco", NameRu: "Монако"},
			"bf": &Country{Id: 122, Code: "bf", NameEn: "Burkina Faso", NameRu: "Буркина-Фасо"},
			"ae": &Country{Id: 123, Code: "ae", NameEn: "United Arab Emirates", NameRu: "Объединенные Арабские Эмираты"},
			"mr": &Country{Id: 124, Code: "mr", NameEn: "Mauritania", NameRu: "Мавритания"},
			"gn": &Country{Id: 125, Code: "gn", NameEn: "Guinea", NameRu: "Гвинея"},
			"zw": &Country{Id: 126, Code: "zw", NameEn: "Zimbabwe", NameRu: "Зимбабве"},
			"gi": &Country{Id: 127, Code: "gi", NameEn: "Gibraltar", NameRu: "Гибралтар"},
			"kh": &Country{Id: 128, Code: "kh", NameEn: "Cambodia", NameRu: "Камбоджа"},
			"bh": &Country{Id: 129, Code: "bh", NameEn: "Bahrain", NameRu: "Бахрейн"},
			"tl": &Country{Id: 130, Code: "tl", NameEn: "Timor-Leste", NameRu: "Тимор-Лешти"},
			"mv": &Country{Id: 131, Code: "mv", NameEn: "Maldives", NameRu: "Мальдивы"},
			"rw": &Country{Id: 132, Code: "rw", NameEn: "Rwanda", NameRu: "Rwanda"},
			"kw": &Country{Id: 133, Code: "kw", NameEn: "Kuwait", NameRu: "Кувейт"},
			"uy": &Country{Id: 134, Code: "uy", NameEn: "Uruguay", NameRu: "Уругвай"},
			"mo": &Country{Id: 135, Code: "mo", NameEn: "Macau", NameRu: "Macau"},
			"is": &Country{Id: 136, Code: "is", NameEn: "Iceland", NameRu: "Iceland"},
			"qa": &Country{Id: 137, Code: "qa", NameEn: "Qatar", NameRu: "Qatar"},
			"mu": &Country{Id: 138, Code: "mu", NameEn: "Mauritius", NameRu: "Mauritius"},
			"aw": &Country{Id: 139, Code: "aw", NameEn: "Aruba", NameRu: "Aruba"},
			"ci": &Country{Id: 140, Code: "ci", NameEn: "Cote D'Ivoire", NameRu: "Cote D'Ivoire"},
			"cv": &Country{Id: 141, Code: "cv", NameEn: "Cape Verde", NameRu: "Cape Verde"},
			"tt": &Country{Id: 142, Code: "tt", NameEn: "Trinidad and Tobago", NameRu: "Trinidad and Tobago"},
			"sn": &Country{Id: 143, Code: "sn", NameEn: "Senegal", NameRu: "Senegal"},
			"np": &Country{Id: 144, Code: "np", NameEn: "Nepal", NameRu: "Nepal"},
			"gt": &Country{Id: 145, Code: "gt", NameEn: "Guatemala", NameRu: "Guatemala"},
			"ax": &Country{Id: 146, Code: "ax", NameEn: "Aland Islands", NameRu: "Aland Islands"},
			"gu": &Country{Id: 147, Code: "gu", NameEn: "Guam", NameRu: "Guam"},
			"gh": &Country{Id: 148, Code: "gh", NameEn: "Ghana", NameRu: "Ghana"},
			"pf": &Country{Id: 149, Code: "pf", NameEn: "French Polynesia", NameRu: "French Polynesia"},
			"sm": &Country{Id: 5070, Code: "sm", NameEn: "San Marino", NameRu: "San Marino"},
			"sl": &Country{Id: 5071, Code: "sl", NameEn: "Sierra Leone", NameRu: "Sierra Leone"},
			"gq": &Country{Id: 5072, Code: "gq", NameEn: "Equatorial Guinea", NameRu: "Equatorial Guinea"},
			"gg": &Country{Id: 5073, Code: "gg", NameEn: "Guernsey", NameRu: "Guernsey"},
			"hn": &Country{Id: 5074, Code: "hn", NameEn: "Honduras", NameRu: "Honduras"},
			"ug": &Country{Id: 5075, Code: "ug", NameEn: "Uganda", NameRu: "Uganda"},
			"dj": &Country{Id: 5076, Code: "dj", NameEn: "Djibouti", NameRu: "Djibouti"},
			"cd": &Country{Id: 5077, Code: "cd", NameEn: "Congo, The Democratic Republic of the", NameRu: "Congo, The Democratic Republic of the"},
			"pr": &Country{Id: 5078, Code: "pr", NameEn: "Puerto Rico", NameRu: "Puerto Rico"},
			"bj": &Country{Id: 5079, Code: "bj", NameEn: "Benin", NameRu: "Benin"},
			"mz": &Country{Id: 5080, Code: "mz", NameEn: "Mozambique", NameRu: "Mozambique"},
			"sv": &Country{Id: 5081, Code: "sv", NameEn: "El Salvador", NameRu: "El Salvador"},
			"ga": &Country{Id: 5082, Code: "ga", NameEn: "Gabon", NameRu: "Gabon"},
			"mg": &Country{Id: 5083, Code: "mg", NameEn: "Madagascar", NameRu: "Madagascar"},
			"je": &Country{Id: 5084, Code: "je", NameEn: "Jersey", NameRu: "Jersey"},
			"fo": &Country{Id: 5085, Code: "fo", NameEn: "Faroe Islands", NameRu: "Faroe Islands"},
			"ne": &Country{Id: 5086, Code: "ne", NameEn: "Niger", NameRu: "Niger"},
			"py": &Country{Id: 5087, Code: "py", NameEn: "Paraguay", NameRu: "Paraguay"},
			"jm": &Country{Id: 5088, Code: "jm", NameEn: "Jamaica", NameRu: "Jamaica"},
			"bb": &Country{Id: 5089, Code: "bb", NameEn: "Barbados", NameRu: "Barbados"},
			"kp": &Country{Id: 5090, Code: "kp", NameEn: "Korea, Democratic People's Republic of", NameRu: "Korea, Democratic People's Republic of"},
			"gy": &Country{Id: 5091, Code: "gy", NameEn: "Guyana", NameRu: "Guyana"},
			"ml": &Country{Id: 5092, Code: "ml", NameEn: "Mali", NameRu: "Mali"},
			"om": &Country{Id: 5093, Code: "om", NameEn: "Oman", NameRu: "Oman"},
			"ad": &Country{Id: 5095, Code: "ad", NameEn: "Andorra", NameRu: "Andorra"},
			"fj": &Country{Id: 5096, Code: "fj", NameEn: "Fiji", NameRu: "Fiji"},
			"nc": &Country{Id: 5097, Code: "nc", NameEn: "New Caledonia", NameRu: "New Caledonia"},
			"bw": &Country{Id: 5098, Code: "bw", NameEn: "Botswana", NameRu: "Botswana"},
			"gp": &Country{Id: 5099, Code: "gp", NameEn: "Guadeloupe", NameRu: "Guadeloupe"},
			"bz": &Country{Id: 5100, Code: "bz", NameEn: "Belize", NameRu: "Belize"},
			"bn": &Country{Id: 5101, Code: "bn", NameEn: "Brunei Darussalam", NameRu: "Brunei Darussalam"},
			"sc": &Country{Id: 5102, Code: "sc", NameEn: "Seychelles", NameRu: "Seychelles"},
			"im": &Country{Id: 5103, Code: "im", NameEn: "Isle of Man", NameRu: "Isle of Man"},
			"ky": &Country{Id: 5104, Code: "ky", NameEn: "Cayman Islands", NameRu: "Cayman Islands"},
			"mw": &Country{Id: 5105, Code: "mw", NameEn: "Malawi", NameRu: "Malawi"},
			"re": &Country{Id: 5106, Code: "re", NameEn: "Reunion", NameRu: "Reunion"},
			"an": &Country{Id: 5107, Code: "an", NameEn: "Netherlands Antilles", NameRu: "Netherlands Antilles"},
			"vi": &Country{Id: 5108, Code: "vi", NameEn: "Virgin Islands, U.S.", NameRu: "Virgin Islands, U.S."},
			"gf": &Country{Id: 5109, Code: "gf", NameEn: "French Guiana", NameRu: "French Guiana"},
			"zm": &Country{Id: 5110, Code: "zm", NameEn: "Zambia", NameRu: "Zambia"},
			"li": &Country{Id: 5111, Code: "li", NameEn: "Liechtenstein", NameRu: "Liechtenstein"},
			"mq": &Country{Id: 5112, Code: "mq", NameEn: "Martinique", NameRu: "Martinique"},
			"bm": &Country{Id: 5113, Code: "bm", NameEn: "Bermuda", NameRu: "Bermuda"},
			"sz": &Country{Id: 5114, Code: "sz", NameEn: "Swaziland", NameRu: "Swaziland"},
			"ap": &Country{Id: 5115, Code: "ap", NameEn: "Asia/Pacific Region", NameRu: "Asia/Pacific Region"},
			"cu": &Country{Id: 5116, Code: "cu", NameEn: "Cuba", NameRu: "Cuba"},
			"sr": &Country{Id: 5117, Code: "sr", NameEn: "Suriname", NameRu: "Suriname"},
			"tg": &Country{Id: 5118, Code: "tg", NameEn: "Togo", NameRu: "Togo"},
			"cg": &Country{Id: 5119, Code: "cg", NameEn: "Congo", NameRu: "Congo"},
			"as": &Country{Id: 5120, Code: "as", NameEn: "American Samoa", NameRu: "American Samoa"},
			"cm": &Country{Id: 5121, Code: "cm", NameEn: "Cameroon", NameRu: "Cameroon"},
			"bt": &Country{Id: 5122, Code: "bt", NameEn: "Bhutan", NameRu: "Bhutan"},
			"mp": &Country{Id: 5123, Code: "mp", NameEn: "Northern Mariana Islands", NameRu: "Northern Mariana Islands"},
			"pg": &Country{Id: 5124, Code: "pg", NameEn: "Papua New Guinea", NameRu: "Papua New Guinea"},
			"tc": &Country{Id: 5125, Code: "tc", NameEn: "Turks and Caicos Islands", NameRu: "Turks and Caicos Islands"},
			"td": &Country{Id: 5126, Code: "td", NameEn: "Chad", NameRu: "Chad"},
			"ht": &Country{Id: 5127, Code: "ht", NameEn: "Haiti", NameRu: "Haiti"},
			"kn": &Country{Id: 5128, Code: "kn", NameEn: "Saint Kitts and Nevis", NameRu: "Saint Kitts and Nevis"},
			"dm": &Country{Id: 5129, Code: "dm", NameEn: "Dominica", NameRu: "Dominica"},
			"so": &Country{Id: 5130, Code: "so", NameEn: "Somalia", NameRu: "Somalia"},
			"gd": &Country{Id: 5131, Code: "gd", NameEn: "Grenada", NameRu: "Grenada"},
			"lc": &Country{Id: 5132, Code: "lc", NameEn: "Saint Lucia", NameRu: "Saint Lucia"},
			"km": &Country{Id: 5133, Code: "km", NameEn: "Comoros", NameRu: "Comoros"},
			"pm": &Country{Id: 5134, Code: "pm", NameEn: "Saint Pierre and Miquelon", NameRu: "Saint Pierre and Miquelon"},
			"gm": &Country{Id: 5135, Code: "gm", NameEn: "Gambia", NameRu: "Gambia"},
			"gl": &Country{Id: 5136, Code: "gl", NameEn: "Greenland", NameRu: "Greenland"},
			"cf": &Country{Id: 5137, Code: "cf", NameEn: "Central African Republic", NameRu: "Central African Republic"},
			"vg": &Country{Id: 5138, Code: "vg", NameEn: "Virgin Islands, British", NameRu: "Virgin Islands, British"},
			"pw": &Country{Id: 5139, Code: "pw", NameEn: "Palau", NameRu: "Palau"},
			"fm": &Country{Id: 5140, Code: "fm", NameEn: "Micronesia, Federated States of", NameRu: "Micronesia, Federated States of"},
			"bi": &Country{Id: 5141, Code: "bi", NameEn: "Burundi", NameRu: "Burundi"},
			"mh": &Country{Id: 5142, Code: "mh", NameEn: "Marshall Islands", NameRu: "Marshall Islands"},
			"vc": &Country{Id: 5143, Code: "vc", NameEn: "Saint Vincent and the Grenadines", NameRu: "Saint Vincent and the Grenadines"},
			"ag": &Country{Id: 5144, Code: "ag", NameEn: "Antigua and Barbuda", NameRu: "Antigua and Barbuda"},
			"vu": &Country{Id: 5145, Code: "vu", NameEn: "Vanuatu", NameRu: "Vanuatu"},
			"ai": &Country{Id: 5146, Code: "ai", NameEn: "Anguilla", NameRu: "Anguilla"},
			"ck": &Country{Id: 5147, Code: "ck", NameEn: "Cook Islands", NameRu: "Cook Islands"},
			"ls": &Country{Id: 5148, Code: "ls", NameEn: "Lesotho", NameRu: "Lesotho"},
			"gw": &Country{Id: 5149, Code: "gw", NameEn: "Guinea-Bissau", NameRu: "Guinea-Bissau"},
			"yt": &Country{Id: 10722, Code: "yt", NameEn: "Mayotte", NameRu: "Mayotte"},
			"fk": &Country{Id: 10723, Code: "fk", NameEn: "Falkland Islands (Malvinas)", NameRu: "Falkland Islands (Malvinas)"},
			"er": &Country{Id: 10724, Code: "er", NameEn: "Eritrea", NameRu: "Eritrea"},
			"va": &Country{Id: 10725, Code: "va", NameEn: "Holy See (Vatican City State)", NameRu: "Holy See (Vatican City State)"},
			"tk": &Country{Id: 10726, Code: "tk", NameEn: "Tokelau", NameRu: "Tokelau"},
			"ws": &Country{Id: 10727, Code: "ws", NameEn: "Samoa", NameRu: "Samoa"},
			"sb": &Country{Id: 48255, Code: "sb", NameEn: "Solomon Islands", NameRu: "Solomon Islands"},
			"ms": &Country{Id: 48256, Code: "ms", NameEn: "Montserrat", NameRu: "Montserrat"},
			"to": &Country{Id: 48257, Code: "to", NameEn: "Tonga", NameRu: "Tonga"},
			"ki": &Country{Id: 48259, Code: "ki", NameEn: "Kiribati", NameRu: "Kiribati"},
			"st": &Country{Id: 56121, Code: "st", NameEn: "Sao Tome and Principe", NameRu: "Sao Tome and Principe"},
		},
	}

	return c
}