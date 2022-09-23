package homekit

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/log"
	"github.com/sirupsen/logrus"
)

func StartHomeKitBridge(ctx *core.Ctx, ctrl *somfy.Controller) {
	log.Debug.Enable()
	bridge := accessory.NewBridge(accessory.Info{Name: "SOMFY-RTS-BRIDGE"})
	bridge.A.Id = 1

	accessories := make([]*accessory.A, len(ctrl.GetDevices()))
	for i, device := range ctrl.GetDevices() {
		cover := NewWindowCovering(device, ctx)
		accessories[i] = cover.A
	}

	// Store the data in the "./db" directory.
	fs := hap.NewFsStore(ctx.Config.HomekitConfigPath)

	// Create the hap server.
	server, err := hap.NewServer(fs, bridge.A, accessories...)
	if err != nil {
		logrus.Errorf("error creating home-kit server")
	}
	server.Pin = ctx.Config.HomekitPin
	server.Addr = fmt.Sprintf(":%s", ctx.Config.HomekitPort)

	// Setup a listener for interrupts and SIGTERM signals
	// to stop the server.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	hapCtx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		// Stop delivering signals.
		signal.Stop(c)
		// Cancel the context to stop the server.
		cancel()
	}()

	// Run the server.
	go server.ListenAndServe(hapCtx)

	logrus.Info("started home-kit bridge")
}
