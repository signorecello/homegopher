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

type BinarySensor struct {
	ID     string
	State  state.State
	Client ha.Connection
}


func (bs *BinarySensor) GetState() state.State {
	conn := bs.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) SetState(newState string, attributes state.Attributes) state.State {
	body := struct {
		State      string     `json:"state"`
		Attributes state.Attributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := bs.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state state.State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) Listen() chan events.StateChangedEvent {
	if ha.BSensorSubs[bs.ID] == nil {
		ha.BSensorSubs[bs.ID] = make(chan events.StateChangedEvent)
	}

	return ha.BSensorSubs[bs.ID]
}


func NewBinarySensor(ID string, c ha.Connection) BinarySensor {
	return BinarySensor{ID: ID, Client: c}
}
