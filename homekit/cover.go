package homekit

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

type Cover struct {
	*accessory.Accessory
	WindowCovering *service.WindowCovering
	device         *somfy.Device
	cmdChan        chan core.DeviceCmd
}

func (c *Cover) OnPositionStateUpdate(pos int) {
	logrus.Infof("client changed position-state of %s to %d", c.device.Id, pos)
}

func (c *Cover) OnTargetPositionUpdate(pos int) {
	logrus.Infof("client changed target-position of %s to %d", c.device.Id, pos)
	cmd := somfy.CmdMy
	switch pos {
	case 0:
		cmd = somfy.CmdDown
		break
	case 100:
		cmd = somfy.CmdUp
		break
	}

	c.cmdChan <- core.DeviceCmd{
		Device: c.device.Id,
		Cmd:    cmd,
	}
	c.WindowCovering.CurrentPosition.SetValue(pos)
}

func (c *Cover) OnCurrentPositionUpdate(pos int) {
	logrus.Infof("client changed current-position of %s to %d", c.device.Id, pos)
}

// NewWindow returns a window which implements model.NewWindow.
func NewWindowCovering(device *somfy.Device, ctx *core.Ctx) *Cover {
	cover := Cover{device: device, cmdChan: ctx.CommandChannel}
	cover.Accessory = accessory.New(accessory.Info{
		Name:         device.Name,
		Manufacturer: "Somfy",
		Model:        "Cover",
		ID:           uint64(device.Address),
	}, accessory.TypeWindowCovering)
	cover.WindowCovering = service.NewWindowCovering()
	cover.AddService(cover.WindowCovering.Service)

	cover.WindowCovering.PositionState.OnValueRemoteUpdate(cover.OnPositionStateUpdate)
	cover.WindowCovering.TargetPosition.OnValueRemoteUpdate(cover.OnTargetPositionUpdate)
	cover.WindowCovering.CurrentPosition.OnValueRemoteUpdate(cover.OnCurrentPositionUpdate)

	return &cover
}
