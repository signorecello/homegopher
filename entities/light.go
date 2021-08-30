package entities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/signorecello/homegopher/ha"
	"github.com/signorecello/homegopher/state"
	"github.com/signorecello/homegopher/events"
	"github.com/signorecello/homegopher/service"
)


type LightServiceCall struct {
	Service string `json:"service"`
	ServiceOpts service.ServiceOpts `json:"service_data"`
}

func (lsc *LightServiceCall) GetServiceOpts() service.ServiceOpts {
	return lsc.ServiceOpts
}
func (lsc *LightServiceCall) SetServiceOpts(so service.ServiceOpts) service.ServiceOpts {
	if so != nil {
		lsc.ServiceOpts = so
	} else {
		lsc.ServiceOpts = &LightOpts{}
	}
	return lsc.ServiceOpts
}
func (lsc *LightServiceCall) GetService() string {
	return lsc.Service
}
func (lsc *LightServiceCall) SetService(service string) {
	lsc.Service = service
}




type LightOpts struct {
	EntityID string `json:"entity_id"`
	Kelvin string `json:"kelvin,omitempty"`
	Brightness string `json:"brightness,omitempty"`
}

func (lo *LightOpts) SetEntityID(entity string) {
	lo.EntityID = entity
}



type Light struct {
	ID     string
	Domain string
	State  state.State
	Client ha.Connection
	LightSubs map[string]chan events.StateChangedEvent
}

func (l *Light) GetDomain() string {
	return l.Domain
}
func (l *Light) GetID() string {
	return l.ID
}
func (l *Light) GetClient() ha.Connection {
	return l.Client
}



func (l *Light) GetState() state.State {
	conn := l.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "light", l.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}


func (l *Light) TurnOn(opts LightOpts) state.State {
	state := Change(l, &LightServiceCall{
		Service: "turn_on",
	}, &opts);
	return state
}


func (l *Light) TurnOff() state.State {
	state := Change(l, &LightServiceCall{
		Service: "turn_off",
	}, nil);
	return state
}


func (l *Light) Listen() chan events.StateChangedEvent {
	if ha.LightSubs[l.ID] == nil {
		ha.LightSubs[l.ID] = make(chan events.StateChangedEvent)
	}

	return ha.LightSubs[l.ID]
}


func NewLight(ID string, c ha.Connection) Light {
	return Light{ID: ID, Client: c, Domain: "light"}
}

