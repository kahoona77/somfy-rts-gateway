package core

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitApp() *Ctx {
	formatter := &logrus.TextFormatter{}
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)

	conf := LoadConfiguration()

	cmdChan := make(chan DeviceCmd, 50)
	ctx := &Ctx{Config: &conf, CommandChannel: cmdChan}

	return ctx
}
