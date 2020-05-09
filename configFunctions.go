package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const configPath = "./config.json"

var ConfigObj Config

// Reads and parses the config file,
func init() {
	var err error
	var data []byte
	ConfigObj = Config{}

	if _, err = os.Stat(configPath); err == nil {
		data, err = ioutil.ReadFile(configPath)
		if err != nil {
			log.Fatalf("reading config file: %v", err)
		}
		err = json.Unmarshal(data, &ConfigObj)
		if err != nil {
			log.Fatalf("unmarshalling config file: %v", err)
		}
	} else {
		log.Fatalf("config file not found: %v", err)
	}
	log.Println("config reading was successful")
}

type Config struct {
	EmailSettings EmailSettings `json:"emailsettings"`
	HttpPort      int           `json:"httpport"`
	HttpsPort     int           `json:"httpsport"`
}
type EmailSettings struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}
