package haclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Sensor struct {
	ID string
	State SensorState
	Client Connection
}

type SensorState struct {
	EntityID    string    `json:"entity_id"`
	LastChanged time.Time `json:"last_changed"`
	State       string    `json:"state"`
	//Attributes  struct {} `json:"attributes"`
	LastUpdated time.Time `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
}

func (s Sensor) GetState() SensorState {
	conn := s.Client
	req, _ := http.NewRequest("GET",fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "sensor", s.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state SensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}
