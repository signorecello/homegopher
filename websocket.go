package homegopher

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"time"
)

type Result struct {
	ID      int      `json:"id"`
	Type    string   `json:"type"`
	Success bool     `json:"success"`
	Result  struct{} `json:"result"`
}

type StateChangedEvent struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Event struct {
		EventType string       `json:"event_type"`
		Data      StateChanged `json:"data"`
		TimeFired time.Time    `json:"time_fired"`
		Origin    string       `json:"origin"`
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
	AuthDone  chan bool
}

func (h HAWS) checkLive() {
	if h.KeepAlive.Unix() < time.Now().Unix() {
		h.Conn.Close()

		log.Println("reconnecting")
		NewWS(h.Timeout, h.URL, h.Auth)
	}
}

func NewWS(timeout time.Duration, url string, auth string) {
	dialer := websocket.Dialer{}

	c, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	haws := HAWS{
		URL:       url,
		Auth:      auth,
		Conn:      c,
		Timeout:   timeout,
		KeepAlive: time.Now().Add(timeout),
		AuthDone:  make(chan bool),
	}

	event := make(chan json.RawMessage)
	go haws.listen(event)

	if !<-haws.AuthDone {
		haws.subscribe("state_changed")
	}

	for e := range event {
		var ge struct {
			ID    int    `json:"id"`
			Type  string `json:"type"`
			Event struct {
				EventType string `json:"event_type"`
			}
		}
		_ = json.Unmarshal(e, &ge)
		haws.routeEvent(ge.Event.EventType, e)
	}
}

func (h HAWS) routeEvent(eventType string, event json.RawMessage) {
	switch eventType {
	case "state_changed":
		var sce StateChangedEvent
		_ = json.Unmarshal(event, &sce)

		split := strings.Split(sce.Event.Data.EntityID, ".")
		domain := split[0]
		entity := split[1]

		//log.Println(entity)
		switch domain {
		case "light":
			select {
			case lightSubs[entity] <- sce:
			default:
			}
		case "sensor":
			select {
			case sensorSubs[entity] <- sce:
			default:
			}
		case "binary_sensor":
			select {
			case bSensorSubs[entity] <- sce:
			default:
			}
		case "switch":
			select {
			case swSubs[entity] <- sce:
			default:
			}
		default:
		}
	default:
	}
}

func (h HAWS) subscribe(et string) {
	e := struct {
		ID        int    `json:"id"`
		Type      string `json:"type"`
		EventType string `json:"event_type"`
	}{
		ID:        1,
		Type:      "subscribe_events",
		EventType: et,
	}

	err := h.Conn.WriteJSON(e)
	if err != nil {
		log.Println("write:", err)
	}
}

func (h HAWS) authenticate() {
	auth := struct {
		Type        string `json:"type"`
		AccessToken string `json:"access_token"`
	}{
		Type:        "auth",
		AccessToken: h.Auth,
	}

	err := h.Conn.WriteJSON(auth)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (h HAWS) listen(event chan json.RawMessage) {
	func() {
		for {
			var v json.RawMessage
			err := h.Conn.ReadJSON(&v)

			if err != nil {
				log.Println("read:", err)
				h.checkLive()
				return
			}

			h.checkLive()

			var t struct {
				Type string `json:"type"`
			}

			_ = json.Unmarshal(v, &t)

			switch t.Type {
			case "auth_required":
				var ar struct {
					Type      string `json:"type"`
					HAVersion string `json:"ha_version"`
				}
				_ = json.Unmarshal(v, &ar)
				h.authenticate()
				break
			case "auth_invalid":
				log.Println("Auth failed")
				close(h.AuthDone)
			case "auth_ok":
				log.Println("Auth OK")
				close(h.AuthDone)
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
