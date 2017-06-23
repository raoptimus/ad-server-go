//
// author ra, resmus@gmail.com
// copyright 2014
//
package detect

import (
	"regexp"
	"sync"
	"tc/openrtbex"
)

type (
	DeviceDetector struct {
		cache           *deviceCache
		browserPatterns []struct {
			reg *regexp.Regexp
			val openrtbex.Browser
		}
		osPatterns []struct {
			reg *regexp.Regexp
			val openrtbex.Os
		}
		devicePatterns []struct {
			reg *regexp.Regexp
			val openrtbex.Device
		}
	}
	deviceItem struct {
		device  openrtbex.Device
		os      openrtbex.Os
		browser openrtbex.Browser
	}
	deviceCache struct {
		sync.RWMutex
		items map[string]*deviceItem
	}
)

//todo slice for pos
func NewDeviceDetector() *DeviceDetector {
	return &DeviceDetector{
		cache: &deviceCache{
			items: make(map[string]*deviceItem),
		},
		browserPatterns: []struct {
			reg *regexp.Regexp
			val openrtbex.Browser
		}{
			{regexp.MustCompile(`MSIE\s([0-9]{1,}[.0-9]{0,})|Trident\/.*rv:([0-9]{1,}[.0-9]{0,})`), openrtbex.BrowserIE},
			{regexp.MustCompile(`Safari`), openrtbex.BrowserSafari},
			{regexp.MustCompile(`Firefox`), openrtbex.BrowserFirefox},
			{regexp.MustCompile(`Chrome`), openrtbex.BrowserChrome},
			{regexp.MustCompile(`Opera\sMini`), openrtbex.BrowserOperaMini},
			{regexp.MustCompile(`Opera\sMobi`), openrtbex.BrowserOperaMobile},
			{regexp.MustCompile(`Opera\/{1,}[\.0-9]{1,}`), openrtbex.BrowserOpera},
		},
		osPatterns: []struct {
			reg *regexp.Regexp
			val openrtbex.Os
		}{
			{regexp.MustCompile(`iOs|iPod|iPhone|iPad`), openrtbex.OsIOs},
			{regexp.MustCompile(`Android`), openrtbex.OsAndroid},
			{regexp.MustCompile(`Symbian|SymbOS|Series60|Series40|SYB-[0-9]+|S60`), openrtbex.OsSymbian},
			{regexp.MustCompile(`Windows\sCE.*?(PPC|Smartphone|Mobile|[0-9]{3}x[0-9]{3})|Window\sMobile|Windows\sPhone\s[0-9.]+|WCE`), openrtbex.OsWindows},
			{regexp.MustCompile(`Windows\sPhone\sOS|XBLWP7|ZuneWP7`), openrtbex.OsWindows},
			{regexp.MustCompile(`blackberry|rim\stablet\sos`), openrtbex.OsBlackBerry},
			{regexp.MustCompile(`Macintosh`), openrtbex.OsMacintosh},
			{regexp.MustCompile(`Linux`), openrtbex.OsLinux},
		},
		devicePatterns: []struct {
			reg *regexp.Regexp
			val openrtbex.Device
		}{
			// mobile
			{regexp.MustCompile(`iPhone|iPod|Mobile`), openrtbex.DeviceMobile},
			{regexp.MustCompile(`BlackBerry|rim[0-9]+`), openrtbex.DeviceMobile},
			// tablet
			{regexp.MustCompile(`Tablet`), openrtbex.DeviceTablet},
		},
	}
}

func (s *DeviceDetector) Detect(ua string) (device openrtbex.Device, os openrtbex.Os, browser openrtbex.Browser) {
	s.cache.RLock()
	item, ok := s.cache.items[ua]
	s.cache.RUnlock()

	if ok {
		return item.device, item.os, item.browser
	}

	browser = s.detectBrowser(ua)

	if browser == openrtbex.BrowserIE {
		os = openrtbex.OsWindows
	} else {
		os = s.detectOs(ua)
	}

	if browser == openrtbex.BrowserOperaMini || browser == openrtbex.BrowserOperaMobile {
		device = openrtbex.DeviceMobile
	} else {
		device = s.detectDevice(ua)

		//TODO
		if device == openrtbex.DeviceDesktop && os == openrtbex.OsAndroid {
			device = openrtbex.DeviceTablet
		}
	}

	item = &deviceItem{
		device:  device,
		os:      os,
		browser: browser,
	}

	s.cache.Lock()
	s.cache.items[ua] = item
	s.cache.Unlock()

	return device, os, browser
}

func (s *DeviceDetector) detectDevice(ua string) openrtbex.Device {
	for _, p := range s.devicePatterns {
		if p.reg.MatchString(ua) {
			return p.val
		}
	}

	return openrtbex.DeviceDesktop
}

func (s *DeviceDetector) detectOs(ua string) openrtbex.Os {
	for _, p := range s.osPatterns {
		if p.reg.MatchString(ua) {
			return p.val
		}
	}

	return openrtbex.OsUnknown
}

func (s *DeviceDetector) detectBrowser(ua string) openrtbex.Browser {
	for _, p := range s.browserPatterns {
		if p.reg.MatchString(ua) {
			return p.val
		}
	}

	return openrtbex.BrowserUnknown
}
