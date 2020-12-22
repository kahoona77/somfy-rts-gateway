package core

import (
	"fmt"
	"os"
)

//AppConfig the Emerald config
type AppConfig struct {
	Port               string
	BasePath           string
	DevicesFile        string
	SignalduinoAddress string
	HomekitConfigPath  string
	HomekitPort        string
	HomekitPin         string
}

func (c *AppConfig) AbsolutePath(path string) string {
	return fmt.Sprintf("%s/%s", c.BasePath, path)
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

	conf.HomekitConfigPath = os.Getenv("HOMEKIT_CONFIG_PATH")
	if conf.HomekitConfigPath == "" {
		conf.HomekitConfigPath = "./db"
	}

	conf.HomekitPort = os.Getenv("HOMEKIT_CONFIG_PORT")
	if conf.HomekitPort == "" {
		conf.HomekitPort = "9123"
	}

	conf.HomekitPin = os.Getenv("HOMEKIT_CONFIG_PIN")
	if conf.HomekitPin == "" {
		conf.HomekitPin = "12344321"
	}

	return conf
}
