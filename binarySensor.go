package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BinarySensor struct {
	ID     string
	State  State
	Client Connection
}

var bSensorSubs = make(map[string]chan StateChangedEvent)

func (bs *BinarySensor) GetState() State {
	conn := bs.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) SetState(newState string, attributes Attributes) State {
	body := struct {
		State      string     `json:"state"`
		Attributes Attributes `json:"attributes"`
	}{newState, attributes}
	reqBody, _ := json.Marshal(body)

	conn := bs.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "binary_sensor", bs.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (bs *BinarySensor) Listen() chan StateChangedEvent {
	if bSensorSubs[bs.ID] == nil {
		bSensorSubs[bs.ID] = make(chan StateChangedEvent)
	}

	return bSensorSubs[bs.ID]
}
