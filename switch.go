package haclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Switch struct {
	ID     string
	State  SwitchState
	Client Connection
}

type SwitchAttributes struct {
	Test string `json:"test"`
}

type SwitchState struct {
	EntityID    string           `json:"entity_id"`
	LastChanged time.Time        `json:"last_changed"`
	State       string           `json:"state"`
	Attributes  SwitchAttributes `json:"attributes"`
	LastUpdated time.Time        `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
}

func (s Switch) GetState() SwitchState {
	conn := s.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "switch", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state SwitchState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s Switch) SetState(newState string, attributes SwitchAttributes) SwitchState {
	body := struct {
		State      string           `json:"state"`
		Attributes SwitchAttributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := s.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "switch", s.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state SwitchState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s Switch) Change(service string) int {
	body := struct {
		EntityID string `json:"entity_id"`
	}{fmt.Sprintf("%s.%s", "switch", s.ID)}
	reqBody, _ := json.Marshal(body)

	client := s.Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, "switch", service), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Client.Do(req)
	if err != nil {
		log.Println("Error: ", err)
	}

	return res.StatusCode
}

type SwitchStateChanged struct {
	EntityID    string
	SwitchState SwitchState
}

var swSubs = make(map[string]chan SwitchStateChanged)

func ListenSWS(entityID string) chan SwitchStateChanged {
	if swSubs[entityID] == nil {
		swSubs[entityID] = make(chan SwitchStateChanged)
	}

	return swSubs[entityID]
}
