package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"somfyRtsGateway/core"
	"somfyRtsGateway/homekit"
	"somfyRtsGateway/somfy"
	"somfyRtsGateway/web"
	"somfyRtsGateway/web/views"
	"time"
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

	// Start server
	go func() {
		logrus.Infof("listening on :%s/%s", ctx.Config.Port, ctx.Config.BasePath)
		if err := e.Start(fmt.Sprintf(":%s", ctx.Config.Port)); err != nil {
			logrus.Info("shutting down...")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(cancelCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
