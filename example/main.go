package main

import (
	"github.com/joho/godotenv"
	ha "github.com/signorecello/homegopher"
	"log"
	"os"
	"time"
)

func main() {

	_ = godotenv.Load("../.env")

	conn := ha.Connection{
		Prefix:        os.Getenv("PREFIX"),
		Host:          os.Getenv("HOST"),
		Path:          os.Getenv("HOST_PATH"),
		Port:          os.Getenv("PORT"),
		Authorization: os.Getenv("AUTHORIZATION"),
	}
	HA := ha.NewConnection(conn)

	date := HA.NewSensor("date")
	dateState := date.GetState().State
	log.Println(dateState)

	sw := HA.NewSwitch("some_switch")
	sw.Change("turn_on")

	light := HA.NewLight("some_light")
	light.Change("toggle")

	ss := make(chan ha.SensorState)
	ls := make(chan ha.LightState)

	stateChanges := ha.StateChanges{Sensor: ss, Light: ls}
	go ha.NewWS(
		time.Second*10,
		stateChanges,
		os.Getenv("WSURL"),
		os.Getenv("AUTHORIZATION"),
	)

	go func() {
		for s := range ss {
			log.Println("Sensor state change: ", s.EntityID)
		}
	}()

	go func() {
		for s := range ls {
			log.Println("Sensor state change: ", s.EntityID)
		}
	}()

}
