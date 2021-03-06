package config

import (
	"fmt"
	"log"
	"os"
)

var (
	HttpPort string
)

func init() {

	LoadServer()
}

func LoadServer() {
	HttpPort = os.Getenv("PORT")
	HttpPort = "5000"
	if HttpPort == "" {
		log.Fatal("$PORT must be set")
	} else {
		fmt.Println("Port:", HttpPort)
	}

}
