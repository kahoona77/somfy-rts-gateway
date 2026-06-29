package homekit

import (
	"context"
	"fmt"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/log"
	"github.com/sirupsen/logrus"

	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

func StartHomeKitBridge(ctx context.Context, appCtx *core.Ctx, ctrl *somfy.Controller) error {
	log.Debug.Enable()
	bridge := accessory.NewBridge(accessory.Info{Name: "SOMFY-RTS-BRIDGE"})
	bridge.A.Id = 1

	accessories := make([]*accessory.A, len(ctrl.GetDevices()))
	for i, device := range ctrl.GetDevices() {
		cover := NewWindowCovering(device, appCtx)
		accessories[i] = cover.A
	}

	fs := hap.NewFsStore(appCtx.Config.HomekitConfigPath)

	server, err := hap.NewServer(fs, bridge.A, accessories...)
	if err != nil {
		return fmt.Errorf("error creating HomeKit server: %w", err)
	}
	server.Pin = appCtx.Config.HomekitPin
	server.Addr = fmt.Sprintf(":%s", appCtx.Config.HomekitPort)

	go server.ListenAndServe(ctx)

	logrus.Info("started HomeKit bridge")
	return nil
}
