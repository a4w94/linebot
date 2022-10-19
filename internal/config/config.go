package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	HttpPort string
)

var Config struct {
	Port uint
	DB   struct {
		Name     string
		Adapter  string
		Host     string
		Port     string
		User     string
		Password string
	}
}

func init() {
	y, err := ioutil.ReadFile("./internal/config/database.yml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal(y, &Config)
	LoadServer()
}

func LoadServer() {
	HttpPort = os.Getenv("PORT")
	//HttpPort = "5000"
	if HttpPort == "" {
		log.Fatal("$PORT must be set")
	} else {
		fmt.Println("Port:", HttpPort)
	}

}
