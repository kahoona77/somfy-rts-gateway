package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/api"
	"somfyRtsGateway/core"
	"somfyRtsGateway/homekit"
	"somfyRtsGateway/somfy"
)

func main() {
	ctx := core.InitApp()

	ctrl, err := somfy.NewController(ctx)
	if err != nil {
		logrus.Errorf("error creating somfy-controller: %v", err)
	}
	defer ctrl.Close()

	homekit.StartHomeKitBridge(ctx.Config.HomekitConfigPath, ctrl)

	e := echo.New()
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	root := e.Group(ctx.Config.BasePath)
	root.POST("/:device/:cmd", api.SomfyCmd)

	//
	//d.Send(s, somfy.ButtonUp)

	// Listen and server on 0.0.0.0:8080
	logrus.Infof("listening on :%s/%s", ctx.Config.Port, ctx.Config.BasePath)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ctx.Config.Port)))
}
