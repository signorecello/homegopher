package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BinarySensor struct {
	ID     string
	State  BinarySensorState
	Client Connection
}

type BinarySensorAttributes map[string]interface{}

type BinarySensorState struct {
	EntityID    string                 `json:"entity_id"`
	LastChanged time.Time              `json:"last_changed"`
	State       string                 `json:"state"`
	Attributes  BinarySensorAttributes `json:"attributes"`
	LastUpdated time.Time              `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
}

type BinarySensorStateChanged struct {
	EntityID          string
	BinarySensorState BinarySensorState
}

var bSensorSubs = make(map[string]chan BinarySensorStateChanged)


func (bs *BinarySensor) GetState() BinarySensorState {
	conn := bs.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state BinarySensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) SetState(newState string, attributes BinarySensorAttributes) BinarySensorState {
	body := struct {
		State      string                 `json:"state"`
		Attributes BinarySensorAttributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := bs.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state BinarySensorState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) Listen() chan BinarySensorStateChanged {
	if bSensorSubs[bs.ID] == nil {
		bSensorSubs[bs.ID] = make(chan BinarySensorStateChanged)
	}

	return bSensorSubs[bs.ID]
}
