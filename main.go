package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/signalduino"
	"somfyRtsGateway/somfy"
)

var buildVersion string

func main() {
	ctx := core.InitApp(buildVersion)
	defer ctx.Close()

	e := echo.New()
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	//root := e.Group(ctx.AppConfig.BasePath)
	//root.GET("/", views.Index)

	s, err := signalduino.Open()
	if err != nil {
		logrus.Errorf("error opening signalduino: %v", err)
	}
	// Make sure to close it later.
	//defer s.Close()

	s.Version()

	d := &somfy.Device{
		RollingCode:   165,
		Address:       3,
		EncryptionKey: 166,
		Name:          "Oben3",
	}

	d.Send(s, somfy.ButtonUp)

	// Listen and server on 0.0.0.0:8080
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ctx.AppConfig.Port)))
}
