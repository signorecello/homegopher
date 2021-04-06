package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Switch struct {
	ID     string
	State  State
	Client Connection
}

var swSubs = make(map[string]chan StateChangedEvent)

func (s *Switch) GetState() State {
	conn := s.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "switch", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Switch) SetState(newState string, attributes Attributes) State {
	body := struct {
		State      string     `json:"state"`
		Attributes Attributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := s.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "switch", s.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s *Switch) Change(data ServiceCall) int {
	body := data.ServiceData
	if (ServiceCall{}.ServiceData) == body {
		body.EntityID = fmt.Sprintf("%s.%s", "switch", s.ID)
	}

	reqBody, _ := json.Marshal(body)

	client := s.Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, "switch", data.Service), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Client.Do(req)
	if err != nil {
		log.Println("Error: ", err)
	}

	return res.StatusCode
}

func (sw *Switch) Listen() chan StateChangedEvent {
	if swSubs[sw.ID] == nil {
		swSubs[sw.ID] = make(chan StateChangedEvent)
	}

	return swSubs[sw.ID]
}
