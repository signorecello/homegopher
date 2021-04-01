package haclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sensor struct {
	ID     string
	State  SensorState
	Client Connection
}

type SensorAttributes struct {
	Test string `json:"test"`
}

type SensorState struct {
	EntityID    string           `json:"entity_id"`
	LastChanged time.Time        `json:"last_changed"`
	State       string           `json:"state"`
	Attributes  SensorAttributes `json:"attributes"`
	LastUpdated time.Time        `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
}

func (s Sensor) GetState() SensorState {
	conn := s.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state SensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (s Sensor) SetState(newState string, attributes SensorAttributes) SensorState {
	body := struct {
		State      string           `json:"state"`
		Attributes SensorAttributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := s.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state SensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

type SensorStateChanged struct {
	EntityID    string
	SensorState SensorState
}

var sensorSubs = make(map[string]chan SensorStateChanged)

func ListenSS(entityID string) chan SensorStateChanged {
	if sensorSubs[entityID] == nil {
		sensorSubs[entityID] = make(chan SensorStateChanged)
	}

	return sensorSubs[entityID]
}
