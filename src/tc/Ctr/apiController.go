package main

import (
	"errors"
	"github.com/bsm/openrtb"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"tc/bootstrap/config"
	"tc/openrtbex"
)

type (
	ApiController struct {
		storage *Storage
	}
)

func NewApiController(storage *Storage) *ApiController {
	c := &ApiController{
		storage: storage,
	}
	go c.listen()
	return c
}
func (s *ApiController) GetCtrOld(req map[string]interface{}, ctr *float64) error {
	siteId := req["SiteId"].(int)
	adId := req["AdId"].(int)
	adCodeTypeId := req["AdCodeTypeId"].(int)
	countryId := req["CountryId"].(int)

	log.Println(siteId, adId, adCodeTypeId, countryId)

	*ctr = s.storage.GetCtrByAd(adId)

	return nil
}
func (s *ApiController) GetCtr(req *openrtb.Request, ctr *float64) error {
	reqExt := req.Ext["RequestExt"].(openrtbex.RequestExt)
	siteId, err := strconv.Atoi(*req.Site.Id)
	if err != nil {
		return err
	}
	if len(req.Imp) == 0 {
		return errors.New("Imp not set")
	}
	adId, err := strconv.Atoi(*req.Imp[0].Id)
	if err != nil {
		return err
	}
	*ctr = s.storage.GetCtrByAd(adId)
	log.Println(reqExt, siteId, adId)
	return nil
}

func (s *ApiController) listen() {
	n, a := config.String("CtrApiNet", ""), config.String("CtrApiAddr", "")
	if n == "" || a == "" {
		log.Fatalln(errors.New("ApiController.listen error: net or addr is empty"))
	}
	server, err := net.Listen(n, a)
	if err != nil {
		log.Fatalln(err)
	}

	openrtbex.GobRegisterExt()
	rpc.Register(s)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go func(conn net.Conn) {
			recovery()
			defer conn.Close()
			//		    conn.SetDeadline(time.Now().Add(600 * time.Millisecond))
			rpc.ServeConn(conn)
		}(conn)
	}

}
