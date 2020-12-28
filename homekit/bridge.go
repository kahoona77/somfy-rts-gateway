package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/log"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

func StartHomeKitBridge(ctx *core.Ctx, ctrl *somfy.Controller) {
	log.Debug.Enable()
	bridge := accessory.NewBridge(accessory.Info{Name: "SOMFY-RTS-BRIDGE", ID: 1})

	accessories := make([]*accessory.Accessory, len(ctrl.GetDevices()))
	for i, device := range ctrl.GetDevices() {
		cover := NewWindowCovering(device, ctx)
		accessories[i] = cover.Accessory
	}

	config := hc.Config{Pin: ctx.Config.HomekitPin, Port: ctx.Config.HomekitPort, StoragePath: ctx.Config.HomekitConfigPath}
	t, err := hc.NewIPTransport(config, bridge.Accessory, accessories...)
	if err != nil {
		logrus.Errorf("error creating home-kit transport")
	}

	hc.OnTermination(func() {
		<-t.Stop()
	})

	go t.Start()
	logrus.Info("started home-kit bridge")
}
