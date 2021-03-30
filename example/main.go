package main

import (
	"github.com/joho/godotenv"
	ha "github.com/signorecello/homegopher"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()

	conn := ha.Connection{
		Prefix: os.Getenv("PREFIX"),
		Host: os.Getenv("HOST"),
		Path: os.Getenv("HOST_PATH"),
		Port: os.Getenv("PORT"),
		Authorization: os.Getenv("AUTHORIZATION"),
	}
	HA := ha.NewConnection(conn)

	health := HA.GetHealth()
	log.Println(health)

}