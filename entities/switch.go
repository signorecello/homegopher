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




type SwitchServiceCall struct {
	Service string `json:"service"`
	ServiceOpts service.ServiceOpts `json:"service_data"`
}

func (lsc *SwitchServiceCall) GetServiceOpts() service.ServiceOpts {
	return lsc.ServiceOpts
}
func (lsc *SwitchServiceCall) SetServiceOpts(so service.ServiceOpts) service.ServiceOpts {
	if so != nil {
		lsc.ServiceOpts = so
	} else {
		lsc.ServiceOpts = &SwitchOpts{}
	}
	return lsc.ServiceOpts
}
func (lsc *SwitchServiceCall) GetService() string {
	return lsc.Service
}
func (lsc *SwitchServiceCall) SetService(service string) {
	lsc.Service = service
}




type SwitchOpts struct {
	EntityID string `json:"entity_id"`
}

func (lo *SwitchOpts) SetEntityID(entity string) {
	lo.EntityID = entity
}



type Switch struct {
	ID     string
	Domain string
	State  state.State
	Client ha.Connection
	SwitchSubs map[string]chan events.StateChangedEvent
}

func (l *Switch) GetDomain() string {
	return l.Domain
}
func (l *Switch) GetID() string {
	return l.ID
}
func (l *Switch) GetClient() ha.Connection {
	return l.Client
}



func (l *Switch) GetState() state.State {
	conn := l.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "switch", l.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}


func (l *Switch) TurnOn(opts SwitchOpts) state.State {
	state := Change(l, &SwitchServiceCall{
		Service: "turn_on",
	}, &opts);
	return state
}


func (l *Switch) TurnOff() state.State {
	state := Change(l, &SwitchServiceCall{
		Service: "turn_off",
	}, nil);
	return state
}


func (l *Switch) Listen() chan events.StateChangedEvent {
	if ha.SwitchSubs[l.ID] == nil {
		ha.SwitchSubs[l.ID] = make(chan events.StateChangedEvent)
	}

	return ha.SwitchSubs[l.ID]
}


func NewSwitch(ID string, c ha.Connection) Switch {
	return Switch{ID: ID, Client: c, Domain: "switch"}
}
