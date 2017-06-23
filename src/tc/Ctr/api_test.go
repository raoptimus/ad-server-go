package main

import (
	"errors"
	"fmt"
	"github.com/bsm/openrtb"
	"net/rpc"
	"tc/bootstrap/config"
	"tc/openrtbex"
	"testing"
)

func TestGetCtr(t *testing.T) {
	conn, err := connectionToRpc()
	if err != nil {
		t.Fatal(err.Error())
	}

	defer conn.Close()

	siteId := "7"
	adId := "97627"
	req := openrtb.Request{
		Site: &openrtb.Site{
			Id: &siteId,
		},
		Imp: []openrtb.Impression{
			openrtb.Impression{
				Id: &adId,
			},
		},
		Ext: openrtb.Extensions{
			"RequestExt": openrtbex.RequestExt{},
		},
	}
	var ctr float64
	err = conn.Call("ApiController.GetCtr", &req, &ctr)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(ctr)
	if ctr <= 0 {
		t.Fatal("Return data is empty")
	}
}

func TestGetCtrOld(t *testing.T) {
	conn, err := connectionToRpc()
	if err != nil {
		t.Fatal(err.Error())
	}

	defer conn.Close()

	req := map[string]interface{}{
		"SiteId":       7,
		"AdId":         97627,
		"AdCodeTypeId": 1,
		"CountryId":    1,
	}
	var ctr float64
	err = conn.Call("ApiController.GetCtrOld", &req, &ctr)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(ctr)
	if ctr <= 0 {
		t.Fatal("Return data is empty")
	}
}

func connectionToRpc() (cn *rpc.Client, err error) {
	n, a := config.String("CtrApiNet", ""), config.String("CtrApiAddr", "")
	if n == "" || a == "" {
		err = errors.New("ApiController.listen error: net or addr is empty")
	}
	openrtbex.GobRegisterExt()
	cn, err = rpc.Dial(n, a)
	return
}
