package entities

import (
	"fmt"
	"log"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"github.com/signorecello/homegopher/ha"
	"github.com/signorecello/homegopher/state"
	"github.com/signorecello/homegopher/events"
	"github.com/signorecello/homegopher/service"
)


type Entity interface {
	GetState() state.State
	Listen() chan events.StateChangedEvent
	GetID() string
	GetDomain() string
	GetClient() ha.Connection
}


func Change(e Entity, data service.ServiceCall, opts service.ServiceOpts) state.State {
	body := data.SetServiceOpts(opts)

	body.SetEntityID(fmt.Sprintf("%s.%s", e.GetDomain(), e.GetID()))

	reqBody, _ := json.Marshal(body)
	log.Println("BODY: ", string(reqBody))

	client := e.GetClient()

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, e.GetDomain(), data.GetService()), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Client.Do(req)
	log.Println("RES ", res)
	
	resBody, _ := ioutil.ReadAll(res.Body)
	log.Println(string(resBody))

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state

}
