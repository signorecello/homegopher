package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
)

type Light struct {
	ID     string
	State  State
	Client Connection
}

type LightServiceCall struct {
	Service string `json:"service"`
	ServiceData struct {
		EntityID string `json:"entity_id"`
		Kelvin string `json:"kelvin,omitempty"`
		Brightness string `json:"brightness,omitempty"`
	} `json:"service_data"`
}

type LightOpts struct {
	Kelvin string `json:"kelvin,omitempty"`
	Brightness string `json:"brightness,omitempty"`
}


var lightSubs = make(map[string]chan StateChangedEvent)

func (l *Light) GetState() State {
	conn := l.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "light", l.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}


func (l *Light) TurnOn(opts LightOpts) State {
	state := l.Change(LightServiceCall{
		Service: "turn_on",
	}, &opts);
	return state
}


func (l *Light) TurnOff() State {
	state := l.Change(LightServiceCall{
		Service: "turn_off",
	}, nil);
	return state
}

func (l *Light) Change(data LightServiceCall, opts *LightOpts) State {
	if opts != nil {data.ServiceData.Kelvin = opts.Kelvin}
	body := data.ServiceData

	log.Println(body)

	body.EntityID = fmt.Sprintf("%s.%s", "light", l.ID)

	reqBody, _ := json.Marshal(body)
	log.Println("BODY: ", string(reqBody))

	client := l.Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, "light", data.Service), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Client.Do(req)
	log.Println("RES ", res)
	
	resBody, _ := ioutil.ReadAll(res.Body)
	log.Println(string(resBody))

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state

}

func (l *Light) Listen() chan StateChangedEvent {
	if lightSubs[l.ID] == nil {
		lightSubs[l.ID] = make(chan StateChangedEvent)
	}

	return lightSubs[l.ID]
}
