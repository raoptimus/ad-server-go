package tchttp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"tc/detect"
	"tc/openrtbex"
	"time"
)

const SESSION_TTL = time.Hour * 16

type HttpContext struct {
	W http.ResponseWriter
	R *http.Request

	Debug     string
	IsDebug   bool
	SessionId string
	Ip        string
	Ua        string
	Lang      string
	LangIsRu  bool

	Device   openrtbex.Device
	Os       openrtbex.Os
	Browser  openrtbex.Browser
	Operator openrtbex.Operator

	GeoId   int
	City    detect.City
	Country detect.Country
	Region  detect.Region

	ErrCode int
}

func NewHttpContext(w http.ResponseWriter, r *http.Request) (h *HttpContext, err error) {
	h = &HttpContext{
		W:       w,
		R:       r,
		ErrCode: http.StatusInternalServerError,
	}
	w.Header().Set("Server", "12AdServer")

	err = h.setSession()
	return
}

func (s *HttpContext) ValidUa() error {
	if s.Ua == "" || !strings.Contains(s.Ua, "(") || len(s.Ua) < 10 {
		return errors.New("UserAgent is invalidate")
	}
	ua := strings.ToLower(s.Ua)
	if strings.Contains(strings.ToLower(ua), "bot") || strings.Contains(strings.ToLower(ua), "spider") || strings.Contains(strings.ToLower(ua), "curl") {
		return errors.New("This is bot")
	}
	return nil
}

func (s *HttpContext) ValidLang() error {
	if s.Lang == "" {
		return errors.New("Language is empty")
	}
	if !strings.Contains(s.Lang, "-") {
		return errors.New("Language is not correct")
	}
	return nil
}

func (s *HttpContext) RedirectJs(url string) {
	s.W.Header().Add("Content-Type", "text/html;charset=utf-8")
	io.WriteString(s.W, "<html><header></header><body>")
	io.WriteString(s.W, fmt.Sprintf("<script>document.location.href='%s';</script>", url))
	io.WriteString(s.W, fmt.Sprintf("<noscript><meta http-equiv=\"refresh\" content=\"1;url=%s\" /></noscript>", url))
	io.WriteString(s.W, "</body></html>")
}

func (s *HttpContext) WriteError(err error, redirectUrl string) {
	if s.IsDebug {
		s.W.Header().Add("Content-Type", "text/html;charset=utf-8")
		s.W.Header().Add("DEBUG-IP", s.Ip)
		io.WriteString(s.W, "<pre>"+err.Error()+"</pre>")
		io.WriteString(s.W, "<pre>redirect: "+redirectUrl+"</pre>")
		return
	}
	if redirectUrl != "" {
		s.RedirectJs(redirectUrl)
	} else {
		http.Error(s.W, "", s.ErrCode)
	}
}

func (s *HttpContext) setSession() error {
	s.R.ParseForm()
	s.Debug = s.R.Form.Get("debug")
	s.IsDebug = s.Debug != ""
	ip, err := s.getRealIp(s.Debug)
	if err != nil {
		return err
	}
	s.Ip = ip
	s.Ua = s.R.UserAgent()
	s.SessionId = s.getSessionId(s.Ua, s.Ip)
	s.Lang = s.R.Header.Get("Accept-Language")
	s.LangIsRu = strings.Contains(s.Lang, "ru")

	return nil
}

func (s *HttpContext) getRealIp(debugIp string) (ip string, err error) {
	var ips string

	if debugIp != "" {
		ips = debugIp
	} else {
		ips = s.R.RemoteAddr

		proxy := s.R.Header.Get("X-Forwarded-For")
		if proxy != "" {
			ips = proxy
		}
	}

	nip := net.ParseIP(ips)
	if nip == nil {
		err = errors.New("Ip is not valid")
		return
	}

	ip = nip.String()

	if nip.IsInterfaceLocalMulticast() {
		//todo get output addr
		ip = "..."
	}
	return
}

func (s *HttpContext) SetCookie(key, val string) (exists bool) {
	cookie, err := s.R.Cookie(key)
	if err != nil {
		return true
	}
	expire := time.Now().UTC().Add(time.Duration(SESSION_TTL))
	cookie = &http.Cookie{
		Name:    key,
		Value:   val,
		Path:    "/",
		Domain:  s.R.URL.Host,
		Expires: expire,
		MaxAge:  expire.Second(),
	}
	http.SetCookie(s.W, cookie)
	return false
}

func (s *HttpContext) getSessionId(ua, ip string) string {
	var sid string
	cookie, err := s.R.Cookie("sid")
	if err == nil {
		sid = cookie.Value
	}

	if len(sid) != 32 {
		//generate new session id
		k := ua + ip
		h := md5.New()
		h.Write([]byte(k))
		sid = hex.EncodeToString(h.Sum(nil))

		expire := time.Now().Add(time.Duration(SESSION_TTL))
		cookie = &http.Cookie{
			Name:    "sid",
			Value:   sid,
			Path:    "/",
			Domain:  s.R.URL.Host,
			Expires: expire,
			MaxAge:  expire.Second(),
		}
		http.SetCookie(s.W, cookie)
	}
	return sid
}
