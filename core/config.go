package core

import (
	"os"
)

//AppConfig the Emerald config
type AppConfig struct {
	Port         string
	BasePath     string
	DbFile       string
	BuildVersion string
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
		conf.BasePath = "/emerald"
	}

	conf.DbFile = os.Getenv("DB_FILE")
	if conf.DbFile == "" {
		conf.DbFile = "emerald.db"
	}

	return conf
}
