package core

import (
	"os"
)

//AppConfig the Emerald config
type AppConfig struct {
	Port               string
	BasePath           string
	DevicesFile        string
	SignalduinoAddress string
}

//LoadConfiguration loads the configuration file
func LoadConfiguration() AppConfig {
	conf := AppConfig{}

	conf.Port = os.Getenv("PORT")
	if conf.Port == "" {
		conf.Port = "8080"
	}

	conf.BasePath = os.Getenv("BASE_PATH")
	if conf.BasePath == "" {
		conf.BasePath = "/somfy"
	}

	conf.DevicesFile = os.Getenv("DEVICES_CONFIG")
	if conf.DevicesFile == "" {
		conf.DevicesFile = "somfy.yaml"
	}

	conf.SignalduinoAddress = os.Getenv("SIGNALDUINO_ADDRESS")
	if conf.SignalduinoAddress == "" {
		conf.SignalduinoAddress = "COM3"
	}

	return conf
}
