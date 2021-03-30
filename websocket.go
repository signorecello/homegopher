package haclient

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

type Generic struct {
	Type string `json:"type"`
}

type AuthRequired struct {
	Type      string `json:"type"`
	HAVersion string `json:"ha_version"`
}

type Auth struct {
	Type        string `json:"type"`
	AccessToken string `json:"access_token"`
}

type Result struct {
	ID      int      `json:"id"`
	Type    string   `json:"type"`
	Success bool     `json:"success"`
	Result  struct{} `json:"result"`
}

type SubscribeEvents struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	EventType string `json:"event_type"`
}

type GenericEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		EventType string `json:"event_type"`
	}
}

type GenericStateChangedEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		Data struct {
			EntityID string          `json:"entity_id"`
			NewState json.RawMessage `json:"new_state"`
			OldState json.RawMessage `json:"old_state"`
		} `json:"data"`
		EventType string    `json:"event_type"`
		TimeFired time.Time `json:"time_fired"`
		Origin    string    `json:"origin"`
		Context   struct {
			ID       string      `json:"id"`
			ParentID interface{} `json:"parent_id"`
			UserID   string      `json:"user_id"`
		} `json:"context"`
	} `json:"event"`
}

type HAWS struct {
	URL       string
	Auth      string
	Conn      *websocket.Conn
	Timeout   time.Duration
	KeepAlive time.Time
	Done      chan bool
	WatchFor  StateChanges
}

func (h HAWS) checkLive() {
	if h.KeepAlive.Unix() < time.Now().Unix() {
		h.Conn.Close()

		log.Println("reconnecting")
		NewWS(h.Timeout, h.WatchFor, h.URL, h.Auth)
	}
}

type StateChanges struct {
	Sensor chan SensorState
	Light  chan LightState
}

func NewWS(t time.Duration, sc StateChanges, url string, auth string) {
	dialer := websocket.Dialer{}

	c, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	haws := HAWS{
		URL:       url,
		Auth:      auth,
		Conn:      c,
		Timeout:   t,
		KeepAlive: time.Now().Add(t),
		Done:      make(chan bool),
		WatchFor:  sc,
	}

	event := make(chan json.RawMessage)
	go haws.listen(event)
	if <-haws.Done {
		haws.subscribe("state_changed")
	}

	for e := range event {
		var ge GenericEvent
		_ = json.Unmarshal(e, &ge)
		haws.routeEvent(ge.Event.EventType, e)
	}
}

func (h HAWS) routeEvent(eventType string, event json.RawMessage) {
	switch eventType {
	case "state_changed":
		var sce GenericStateChangedEvent
		_ = json.Unmarshal(event, &sce)

		domain := strings.Split(sce.Event.Data.EntityID, ".")[0]
		switch domain {
		case "light":
			var l LightState
			_ = json.Unmarshal(sce.Event.Data.NewState, &l)
			h.WatchFor.Light <- l
			break
		case "sensor":
			var s SensorState
			_ = json.Unmarshal(sce.Event.Data.NewState, &s)
			h.WatchFor.Sensor <- s
			break
		default:
			break
		}

	}
}

func (h HAWS) subscribe(et string) {
	e := SubscribeEvents{
		ID:        1,
		Type:      "subscribe_events",
		EventType: et,
	}

	err := h.Conn.WriteJSON(e)
	if err != nil {
		log.Println("write:", err)
		return
	}

}

func (h HAWS) authenticate() {
	auth := Auth{Type: "auth", AccessToken: h.Auth}
	err := h.Conn.WriteJSON(auth)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (h HAWS) listen(event chan json.RawMessage) {
	go func() {
		defer close(h.Done)
		for {
			var v json.RawMessage
			err := h.Conn.ReadJSON(&v)

			if err != nil {
				log.Println("read:", err)
				h.checkLive()
				return
			}

			h.checkLive()

			var t Generic
			_ = json.Unmarshal(v, &t)

			switch t.Type {
			case "auth_required":
				var ar AuthRequired
				_ = json.Unmarshal(v, &ar)
				h.authenticate()
				break
			case "auth_invalid":
				log.Println("Auth failed")
				h.Done <- true
			case "auth_ok":
				log.Println("Auth OK")
				h.Done <- true
				break
			case "event":
				h.KeepAlive = time.Now().Add(h.Timeout)
				event <- v
			case "result":
				var r Result
				_ = json.Unmarshal(v, &r)
				if !r.Success {
					log.Println("Failed: ", r)
				}
				break
			default:
				log.Println("Some other message")
				break
			}
		}
	}()
}
