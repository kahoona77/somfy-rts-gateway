package homekit

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

func StartHomeKitBridge(conf *core.AppConfig, ctrl *somfy.Controller) {
	bridge := accessory.NewBridge(accessory.Info{Name: "Somfy-RTS-Bridge", ID: 700001})

	accessories := make([]*accessory.Accessory, len(ctrl.GetDevices()))
	for i, device := range ctrl.GetDevices() {
		cover := NewWindowCovering(device, ctrl)
		accessories[i] = cover.Accessory
	}

	config := hc.Config{Pin: conf.HomekitPin, Port: conf.HomekitPort, StoragePath: conf.HomekitConfigPath}
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

type Covers struct {
	*accessory.Accessory
	WindowCovering *service.WindowCovering
}

// NewWindow returns a window which implements model.NewWindow.
func NewWindowCovering(device *somfy.Device, ctrl *somfy.Controller) *Covers {
	acc := Covers{}
	acc.Accessory = accessory.New(accessory.Info{
		Name:         device.Name,
		Manufacturer: "Somfy",
		Model:        "Cover",
		ID:           uint64(device.Address),
	}, accessory.TypeWindowCovering)
	acc.WindowCovering = service.NewWindowCovering()
	acc.AddService(acc.WindowCovering.Service)

	// Log to console when client (e.g. iOS app) changes the value of the on characteristic
	acc.WindowCovering.PositionState.OnValueRemoteUpdate(func(pos int) {
		logrus.Info("Client changed position-state to %d", pos)
	})

	acc.WindowCovering.TargetPosition.OnValueRemoteUpdate(func(pos int) {
		logrus.Info("Client changed target-position to %d", pos)
		switch pos {
		case 0:
			ctrl.SendCmd(device.Id, somfy.ButtonDown)
			break
		case 1:
			ctrl.SendCmd(device.Id, somfy.ButtonUp)
			break
		default:
			ctrl.SendCmd(device.Id, somfy.ButtonMy)
		}
		acc.WindowCovering.CurrentPosition.SetValue(pos)
	})

	acc.WindowCovering.CurrentPosition.OnValueRemoteUpdate(func(pos int) {
		logrus.Info("Client changed current-position to %d", pos)
	})

	return &acc
}
