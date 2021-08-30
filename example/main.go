package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
	"github.com/signorecello/homegopher/ha"
	"github.com/signorecello/homegopher/state"
	"github.com/signorecello/homegopher/entities"
)

func main() {
	_ = godotenv.Load("../.env")

	// you have to initialize the connection first
	conn := ha.Connection{
		Prefix:        os.Getenv("PREFIX"),
		Host:          os.Getenv("HOST"),
		Path:          os.Getenv("HOST_PATH"),
		Port:          os.Getenv("PORT"),
		Authorization: os.Getenv("AUTHORIZATION"),
	}
	HA := ha.NewConnection(conn)

	// example of the creation of a new sensor
	sensor := entities.NewSensor("some_sensor", HA)
	sensorState := sensor.GetState().State
	log.Println(sensorState)

	// setting state requires attributes, it's a bit of a manual list as some sensors have specific attributes
	// just let me know if you need some specific attribute here
	attributes := state.Attributes{}
	sensor.SetState("off", attributes)

	sw := entities.NewSwitch("some_switch", HA)

	sw.TurnOn(entities.SwitchOpts{})

	light := entities.NewLight("some_light", HA)
	light.TurnOn(entities.LightOpts{})

	// don't forget the go keyword before NewWS otherwise the program will hang forever
	go ha.NewWS(
		5*time.Second,
		os.Getenv("WSURL"),
		os.Getenv("AUTHORIZATION"),
	)

	// listening to a specific state channel for a specific sensor
	channel := sensor.Listen()

	// example of how you could listen for it
	go func() {
		time.Sleep(time.Second)
		sensor.SetState("on", attributes)
	}()

	func() {
		val := <-channel
		log.Printf(val.Event.Data.NewState.State)
	}()

}
