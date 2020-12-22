package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/homekit"
	"somfyRtsGateway/somfy"
	"somfyRtsGateway/web"
	"somfyRtsGateway/web/views"
)

func main() {
	ctx := core.InitApp()

	ctrl, err := somfy.NewController(ctx)
	if err != nil {
		logrus.Errorf("error creating somfy-controller: %v", err)
	}
	defer ctrl.Close()
	ctx.Controller = ctrl

	homekit.StartHomeKitBridge(ctx, ctrl)

	e := echo.New()
	e.Renderer = web.NewTemplate(ctx)
	e.Debug = true
	//e.Logger.SetLevel(log.DEBUG)
	//e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	root := e.Group(ctx.Config.BasePath)
	root.GET("/", somfy.ListDevices(ctrl))
	root.GET("/:device", somfy.GetDevice(ctrl))
	root.POST("/:device/:cmd", somfy.Cmd)

	root.Static("/static", "./web/static")
	root.GET("/web", views.Index)
	root.POST("/web/:device/:cmd", views.Cmd)

	// Listen and server on 0.0.0.0:8080
	logrus.Infof("listening on :%s/%s", ctx.Config.Port, ctx.Config.BasePath)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", ctx.Config.Port)))
}
