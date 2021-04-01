package main

import (
	"github.com/joho/godotenv"
	ha "github.com/signorecello/homegopher"
	"log"
	"os"
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

	date := HA.NewSensor("some_sensor")
	dateState := date.GetState().State
	log.Println(dateState)

	sw := HA.NewSwitch("some_switch")
	sw.Change("turn_on")

	light := HA.NewLight("some_light")
	light.Change("toggle")


}
