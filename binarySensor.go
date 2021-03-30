package haclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BinarySensor struct {
	ID string
	State BinarySensorState
	Client Connection
}

type BinarySensorState struct {
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

func (bs BinarySensor) GetState() BinarySensorState {
	conn := bs.Client
	req, _ := http.NewRequest("GET",fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state BinarySensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}
