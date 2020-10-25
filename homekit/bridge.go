package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

func StartHomeKitBridge(ctx *core.Ctx, ctrl *somfy.Controller) {
	bridge := accessory.NewBridge(accessory.Info{Name: "Somfy-RTS-Bridge", ID: 700001})

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
	logrus.Infof("started home-kit bridge")
}
