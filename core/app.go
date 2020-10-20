package core

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitApp(buildVersion string) *Ctx {
	formatter := &logrus.TextFormatter{}
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Infof("emerald-build-version is '%s'", buildVersion)

	conf := LoadConfiguration()
	conf.BuildVersion = buildVersion

	ctx := &Ctx{AppConfig: &conf}

	return ctx
}
