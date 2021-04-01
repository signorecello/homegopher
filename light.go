package haclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Light struct {
	ID     string
	State  LightState
	Client Connection
}

type LightAttributes struct {
	Test string `json:"test"`
}

type LightState struct {
	EntityID    string          `json:"entity_id"`
	LastChanged time.Time       `json:"last_changed"`
	State       string          `json:"state"`
	Attributes  LightAttributes `json:"attributes"`
	LastUpdated time.Time       `json:"last_updated"`
	Context     struct {
		ID       string      `json:"id"`
		ParentID interface{} `json:"parent_id"`
		UserID   string      `json:"user_id"`
	} `json:"context"`
	//
	//Transition int
	//Profile string
	//HsColor [2]float64
	//XYColor [2]float64
	//RGBColor [3]int
	//WhiteValue int
	//ColorTemp int
	//Kelvin int
	//ColorName string
	//Brightness int
	//BrightnessPct int
	//BrightnessStep int
	//Flash string
	//Effect string
}

func (l Light) GetState() LightState {
	conn := l.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "light", l.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state LightState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (l Light) SetState(newState string, attributes LightAttributes) LightState {
	body := struct {
		State      string          `json:"state"`
		Attributes LightAttributes `json:"attributes"`
	}{newState, attributes}

	reqBody, _ := json.Marshal(body)

	conn := l.Client
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "light", l.ID), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state LightState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (l Light) Change(service string) LightState {
	body := struct {
		EntityID string `json:"entity_id"`
	}{fmt.Sprintf("%s.%s", "light", l.ID)}
	reqBody, _ := json.Marshal(body)

	client := l.Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, "light", service), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Client.Do(req)

	var state LightState
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state

}

type LightStateChanged struct {
	EntityID   string
	LightState LightState
}

var lightSubs = make(map[string]chan LightStateChanged)

func ListenLS(entityID string) chan LightStateChanged {
	if lightSubs[entityID] == nil {
		lightSubs[entityID] = make(chan LightStateChanged)
	}

	return lightSubs[entityID]
}
