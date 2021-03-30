package haclient

import (
	"fmt"
	"log"
	"net/http"
)

type Connection struct {
	Prefix string
	Host string
	Path string
	Port string
	Authorization string
	Client *http.Client
	SSL bool
}


func NewConnection(c Connection) Connection {
	prefix := "http://"
	if c.SSL {
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

	conn := Connection{prefix,host, path, port, "Bearer " + c.Authorization, &http.Client{}, c.SSL}
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

func (c Connection) NewSwitch(ID string) Switch {
	return Switch{ID: ID, Client: c}
}

func (c Connection) NewSensor(ID string) Sensor {
	return Sensor{ID: ID, Client: c}
}


func (c Connection) NewBinarySensor(ID string) BinarySensor {
	return BinarySensor{ID: ID, Client: c}
}

func (c Connection) NewLight(ID string) Light {
	return Light{ID: ID, Client: c}
}

