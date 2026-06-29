package homekit

import (
	"time"

	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
	"github.com/sirupsen/logrus"

	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

type Cover struct {
	*accessory.A
	WindowCovering *service.WindowCovering
	device         *somfy.Device
	cmdChan        chan core.DeviceCmd
}

func (c *Cover) OnTargetPositionUpdate(pos int) {
	logrus.Infof("client changed target-position of %s to %d", c.device.Id, pos)

	cmd := somfy.CmdMy
	positionState := characteristic.PositionStateStopped

	switch {
	case pos == 0:
		cmd = somfy.CmdDown
		positionState = characteristic.PositionStateDecreasing
	case pos == 100:
		cmd = somfy.CmdUp
		positionState = characteristic.PositionStateIncreasing
	case pos > c.device.Position:
		cmd = somfy.CmdPosition
		positionState = characteristic.PositionStateIncreasing
	case pos < c.device.Position:
		cmd = somfy.CmdPosition
		positionState = characteristic.PositionStateDecreasing
	}

	c.WindowCovering.PositionState.SetValue(positionState)

	if positionState != characteristic.PositionStateStopped {
		time.AfterFunc(time.Duration(c.device.ClosingDuration)*time.Second, func() {
			c.WindowCovering.PositionState.SetValue(characteristic.PositionStateStopped)
		})
	}

	c.cmdChan <- core.DeviceCmd{
		Device: c.device.Id,
		Cmd:    cmd,
		Pos:    pos,
	}
}

func (c *Cover) OnDeviceUpdate(device *somfy.Device) {
	logrus.Infof("device update %s to position %d", device.Id, device.Position)
	c.WindowCovering.CurrentPosition.SetValue(device.Position)
}

func NewWindowCovering(device *somfy.Device, ctx *core.Ctx) *Cover {
	cover := Cover{device: device, cmdChan: ctx.CommandChannel}
	cover.A = accessory.New(accessory.Info{
		Name:         device.Name,
		Manufacturer: "Somfy",
		Model:        "Cover",
		Firmware:     "somfy-rts-gateway",
	}, accessory.TypeWindowCovering)
	cover.A.Id = 100 + uint64(device.Address)

	cover.WindowCovering = service.NewWindowCovering()
	cover.WindowCovering.PositionState.SetValue(characteristic.PositionStateStopped)
	cover.WindowCovering.CurrentPosition.SetValue(device.Position)
	cover.WindowCovering.TargetPosition.SetValue(device.Position)

	device.OnUpdate(cover.OnDeviceUpdate)
	cover.WindowCovering.TargetPosition.OnValueRemoteUpdate(cover.OnTargetPositionUpdate)

	cover.AddS(cover.WindowCovering.S)

	return &cover
}
