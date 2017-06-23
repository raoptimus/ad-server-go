package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CommandController struct {
}

func NewCommandController() *CommandController {
	res := &CommandController{}
	go res.listenServe()

	return res
}

func (s *CommandController) parseCmd(r *http.Request) int {
	q := r.URL.Query()
	id, _ := strconv.ParseInt(q.Get("id"), 10, 0)
	return int(id)
}

func (s *CommandController) listenServe() {
	http.HandleFunc("/update/journal", func(w http.ResponseWriter, r *http.Request) {

	})
	http.HandleFunc("/get/object", s.actionGetObject)
	http.HandleFunc("/measure/reset", measureResetHandler)
	http.HandleFunc("/measure/set", measureSetHandler)
	http.HandleFunc("/measure/print", measurePrintHandler)

	log.Fatal(http.ListenAndServe(":9090", nil))
}

// http://localhost:9090/get/object?name=ad&id=1
func (m *CommandController) actionGetObject(w http.ResponseWriter, r *http.Request) {
	var obj interface{}
	var err error
	var id int

	q := r.URL.Query()
	name := q.Get("name")
	id, err = strconv.Atoi(q.Get("id"))

	if err == nil {
		switch name {
		case "campaign":
			{
				obj, err = StoreContext.Campaigns.Get(id)
			}
		case "ads":
			{
				obj, err = StoreContext.Campaigns.Get(id)

				if err == nil {
					ads := obj.(*Campaign).Ads
					adList := make([]*Ad, 0)

					for _, ad := range ads {
						adList = append(adList, ad)
					}

					obj = adList
				}
			}
		case "ad":
			{
				obj, err = StoreContext.Campaigns.GetAd(id)
			}
		case "user":
			{
				obj = StoreContext.Users.Get(id)
			}
		}
	}

	var b []byte

	if err == nil {
		b, err = json.Marshal(obj)
	}

	if err != nil {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		io.WriteString(w, err.Error())
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, string(b))
}
