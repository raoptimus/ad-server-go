package openrtbex

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type YaCur struct {
	sync.RWMutex

	Cur    map[Cur]float32
	CurUrl map[Cur]string
}

func NewYaCur() *YaCur {
	yc := &YaCur{
		Cur: map[Cur]float32{
			CurEuro: float32(60),
			CurUsd:  float32(52),
		},
		CurUrl: map[Cur]string{
			CurEuro: "http://news.yandex.ru/quotes/23.html",
			CurUsd:  "http://news.yandex.ru/quotes/1.html",
		},
	}

	go yc.update()
	return yc
}

func (s *YaCur) GetCur(c Cur) float32 {
	s.RLock()
	cur := s.Cur[c]
	s.RUnlock()

	return cur
}

func (s *YaCur) update() {
	for {
		for c, url := range s.CurUrl {
			cur, err := s.downloadCur(url)

			if err != nil {
				log.Println("YaCur error", err)
			} else {
				s.Lock()
				s.Cur[c] = cur
				s.Unlock()
			}
		}

		time.Sleep(time.Duration(time.Hour * 12))
	}
}

func (s *YaCur) downloadCur(url string) (cur float32, err error) {
	resp, err := http.Get(url)

	if err != nil {
		return float32(0), err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	re := regexp.MustCompile(`<td\s+class=\"b-quote__value\"><span\s+class=\"b-quote__sgn\"></span>([0-9]+,[0-9]+)</td>`)
	all := re.FindAllStringSubmatch(string(body), -1)
	sum := float64(0)

	for i, _ := range all {
		cur := all[i][1]
		cur = strings.Replace(cur, ",", ".", 1)
		f, err := strconv.ParseFloat(cur, 64)

		if err != nil {
			return float32(0), err
		}

		sum += f
	}

	avg := sum / float64(len(all))
	cur = float32(avg)

	return cur, nil
}
