package ha

import (
	"fmt"
	"log"
	"net/http"
	"github.com/signorecello/homegopher/events"
)

type Connection struct {
	Prefix string
	Host string
	Path string
	Port string
	Authorization string
	Client *http.Client
}


func NewConnection(c Connection) Connection {
	prefix := c.Prefix
	if prefix == "" {
		prefix = "https://"
	}

	host := c.Host
	if host == "" {
		host = "localhost"
	}

	path := c.Path
	if path == "" {
		path = "/api"
	}

	port := c.Port
	if port == "" {
		port = "8123"
	}

	if c.Authorization == "" {
		log.Fatal("Need to provide a long-lived token")
	}

	conn := Connection{prefix,host, path, port, "Bearer " + c.Authorization, &http.Client{}}
	conn.GetHealth()

	return conn
}


func (c Connection) GetHealth() int {
	req, _ := http.NewRequest("GET",fmt.Sprintf("%s%s:%s%s/", c.Prefix, c.Host, c.Port, c.Path), nil)
	req.Header.Set("Authorization", c.Authorization)

	res, _ := c.Client.Do(req)

	if res == nil || res.StatusCode != 200 {
		//log.Fatal("Can't establish connection, check connection settings")
		return 500
	} else {
		return 200
	}
}


var LightSubs = make(map[string]chan events.StateChangedEvent)
var BSensorSubs = make(map[string]chan events.StateChangedEvent)
var SensorSubs = make(map[string]chan events.StateChangedEvent)
var SwitchSubs = make(map[string]chan events.StateChangedEvent)
