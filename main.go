package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"

	"somfyRtsGateway/core"
	"somfyRtsGateway/homekit"
	"somfyRtsGateway/mqtt"
	"somfyRtsGateway/somfy"
	"somfyRtsGateway/web"
	"somfyRtsGateway/web/views"
)

func main() {
	ctx := core.InitApp()
	ctrl := initController(ctx)
	defer ctrl.Close()

	ctx.Controller = ctrl

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if ctx.Config.HomekitEnabled {
		if err := homekit.StartHomeKitBridge(rootCtx, ctx, ctrl); err != nil {
			logrus.Errorf("HomeKit bridge failed: %v", err)
		}
	}

	if ctx.Config.MqttEnabled {
		bridge := mqtt.NewBridge(ctx.Config, ctx.CommandChannel, ctrl.Devices())
		if err := bridge.Start(); err != nil {
			logrus.Errorf("MQTT bridge start failed: %v", err)
		} else {
			defer bridge.Stop()
		}
	}

	e := newServer(ctx, ctrl)
	startServer(ctx, e)
	waitForShutdown(rootCtx, e)
}

func initController(ctx *core.Ctx) *somfy.Controller {
	ctrl, err := somfy.NewController(ctx)
	if err != nil {
		logrus.Errorf("error creating somfy-controller: %v", err)
	}
	if ctrl == nil {
		logrus.Fatal("somfy-controller is not available")
	}
	return ctrl
}

func newServer(ctx *core.Ctx, ctrl *somfy.Controller) *echo.Echo {
	e := echo.New()
	e.Renderer = web.NewTemplate(ctx)
	e.Debug = true
	e.Use(middleware.CORS())
	e.Use(core.CreateCtx(ctx))

	root := e.Group(ctx.Config.BasePath)
	root.GET("/", somfy.ListDevices(ctrl))
	root.GET("/:device", somfy.GetDevice(ctrl))
	root.POST("/:device/:cmd", somfy.Cmd)
	root.StaticFS("/static", web.StaticFS())
	root.GET("/web", views.Index)
	root.POST("/web/:device/:cmd", views.Cmd)

	return e
}

func startServer(ctx *core.Ctx, e *echo.Echo) {
	go func() {
		logrus.Infof("listening on :%s/%s", ctx.Config.Port, ctx.Config.BasePath)
		if err := e.Start(fmt.Sprintf(":%s", ctx.Config.Port)); err != nil {
			logrus.Info("shutting down...")
		}
	}()
}

func waitForShutdown(ctx context.Context, e *echo.Echo) {
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		e.Logger.Fatal(err)
	}
}
