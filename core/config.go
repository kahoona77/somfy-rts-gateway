package core

import (
	"fmt"
	"os"
)

// AppConfig the Emerald config
type AppConfig struct {
	Port               string
	BasePath           string
	DevicesFile        string
	SignalduinoAddress string
	HomekitConfigPath  string
	HomekitPort        string
	HomekitPin         string
	HomekitEnabled     bool

	MqttEnabled         bool
	MqttBroker          string
	MqttUsername        string
	MqttPassword        string
	MqttClientId        string
	MqttBaseTopic       string
	MqttDiscoveryPrefix string
}

func (c *AppConfig) AbsolutePath(path string) string {
	return fmt.Sprintf("%s/%s", c.BasePath, path)
}

// LoadConfiguration loads the configuration file
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
		conf.SignalduinoAddress = "/dev/ttyUSB0"
	}

	conf.HomekitConfigPath = os.Getenv("HOMEKIT_CONFIG_PATH")
	if conf.HomekitConfigPath == "" {
		conf.HomekitConfigPath = "./db"
	}

	conf.HomekitPort = os.Getenv("HOMEKIT_CONFIG_PORT")
	if conf.HomekitPort == "" {
		conf.HomekitPort = "19123"
	}

	conf.HomekitPin = os.Getenv("HOMEKIT_CONFIG_PIN")
	if conf.HomekitPin == "" {
		conf.HomekitPin = "12344321"
	}

	conf.HomekitEnabled = os.Getenv("HOMEKIT_ENABLED") != "false"

	conf.MqttEnabled = os.Getenv("MQTT_ENABLED") == "true"

	conf.MqttBroker = os.Getenv("MQTT_BROKER")
	if conf.MqttBroker == "" {
		conf.MqttBroker = "tcp://localhost:1883"
	}

	conf.MqttUsername = os.Getenv("MQTT_USERNAME")
	conf.MqttPassword = os.Getenv("MQTT_PASSWORD")

	conf.MqttClientId = os.Getenv("MQTT_CLIENT_ID")
	if conf.MqttClientId == "" {
		conf.MqttClientId = "somfy-rts-gateway"
	}

	conf.MqttBaseTopic = os.Getenv("MQTT_BASE_TOPIC")
	if conf.MqttBaseTopic == "" {
		conf.MqttBaseTopic = "covers"
	}

	conf.MqttDiscoveryPrefix = os.Getenv("MQTT_DISCOVERY_PREFIX")
	if conf.MqttDiscoveryPrefix == "" {
		conf.MqttDiscoveryPrefix = "homeassistant"
	}

	return conf
}
