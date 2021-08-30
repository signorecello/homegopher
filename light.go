package homegopher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Light struct {
	ID     string
	State  State
	Client Connection
}

var lightSubs = make(map[string]chan StateChangedEvent)

func (l *Light) GetState() State {
	conn := l.Client
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s:%s%s/states/%s.%s", conn.Prefix, conn.Host, conn.Port, conn.Path, "light", l.ID), nil)
	req.Header.Set("Authorization", conn.Authorization)

	res, _ := conn.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state
}

func (l *Light) Change(data ServiceCall) State {
	body := data.ServiceData

	if (ServiceCall{}.ServiceData) == body {
		body.EntityID = fmt.Sprintf("%s.%s", "light", l.ID)
	}

	reqBody, _ := json.Marshal(body)

	client := l.Client

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s:%s%s/services/%s/%s", client.Prefix, client.Host, client.Port, client.Path, "light", data.Service), bytes.NewReader(reqBody))
	req.Header.Set("Authorization", client.Authorization)
	req.Header.Set("Content-Type", "application/json")

	res, _ := client.Client.Do(req)

	var state State
	dec := json.NewDecoder(res.Body)
	_ = dec.Decode(&state)

	return state

}

func (l *Light) Listen() chan StateChangedEvent {
	if lightSubs[l.ID] == nil {
		lightSubs[l.ID] = make(chan StateChangedEvent)
	}

	return lightSubs[l.ID]
}
