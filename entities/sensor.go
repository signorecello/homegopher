package entities

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/signorecello/homegopher/ha"
	"github.com/signorecello/homegopher/state"
	"github.com/signorecello/homegopher/events"
)

type Sensor struct {
	ID     string
	State  state.State
	Client ha.Connection
}


func (s *Sensor) GetState() state.State {
	conn := s.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Sensor) SetState(newState string, attributes state.Attributes) state.State {
	body := struct {
		State      string     `json:"state"`
		Attributes state.Attributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := s.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Sensor) Listen() chan events.StateChangedEvent {
	if ha.SensorSubs[s.ID] == nil {
		ha.SensorSubs[s.ID] = make(chan events.StateChangedEvent)
	}

	return ha.SensorSubs[s.ID]
}



func NewSensor(ID string, c ha.Connection) Sensor {
	return Sensor{ID: ID, Client: c}
}
