package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Sensor struct {
	ID     string
	State  State
	Client Connection
}

var sensorSubs = make(map[string]chan StateChangedEvent)

func (s *Sensor) GetState() State {
	conn := s.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Sensor) SetState(newState string, attributes Attributes) State {
	body := struct {
		State      string     `json:"state"`
		Attributes Attributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := s.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Sensor) Listen() chan StateChangedEvent {
	if sensorSubs[s.ID] == nil {
		sensorSubs[s.ID] = make(chan StateChangedEvent)
	}

	return sensorSubs[s.ID]
}
