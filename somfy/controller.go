package somfy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"somfyRtsGateway/core"
	"somfyRtsGateway/signalduino"
)

type Controller struct {
	sig         *signalduino.Signalduino
	devices     []*Device
	devicesFile string
}

func (c *Controller) Signalduino() *signalduino.Signalduino {
	return c.sig
}

func (c *Controller) Devices() []core.Device {
	devices := make([]core.Device, len(c.devices))
	for i, device := range c.devices {
		devices[i] = device
	}
	return devices
}

func NewController(ctx *core.Ctx) (*Controller, error) {
	devices, err := loadDevices(ctx.Config.DevicesFile)
	if err != nil {
		return nil, fmt.Errorf("error loading devices from %s: %v", ctx.Config.DevicesFile, err)
	}

	ctrl := &Controller{
		devices:     devices,
		devicesFile: ctx.Config.DevicesFile,
	}
	go ctrl.listen(ctx.CommandChannel)

	s, err := signalduino.Open(ctx.Config.SignalduinoAddress)
	if err != nil {
		return ctrl, fmt.Errorf("error opening signalduino on address %s: %v", ctx.Config.SignalduinoAddress, err)
	}
	s.Version()
	ctrl.sig = s

	return ctrl, nil
}

func (c *Controller) Close() {
	logrus.Debugf("closing controller...")
	if c == nil || c.sig == nil {
		return
	}
	if err := c.sig.Close(); err != nil {
		logrus.Errorf("error closing signalduino: %v", err)
	}
}

func loadDevices(file string) ([]*Device, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var devices []*Device
	err = yaml.Unmarshal(yamlFile, &devices)
	return devices, err
}

func (c *Controller) listen(queue chan core.DeviceCmd) {
	for dc := range queue {
		if dc.Cmd == "ping" {
			if c.sig == nil {
				logrus.Errorf("ignoring ping command: signalduino is not available")
				continue
			}
			c.sig.Ping()
			continue
		}

		c.SendCmd(dc)
	}
}

func (c *Controller) save() error {
	d, err := yaml.Marshal(c.devices)
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.devicesFile, d, 0644); err != nil {
		return err
	}

	path, _ := filepath.Abs(c.devicesFile)
	logrus.Infof("saved config to: %s", path)

	return nil
}

func (c *Controller) SendCmd(dc core.DeviceCmd) {
	d, err := c.GetDevice(dc.Device)
	if err != nil {
		logrus.Warn(err)
		return
	}

	if c.sig == nil {
		logrus.Errorf("cannot send command %s to device %s: signalduino is not available", dc.Cmd, dc.Device)
		return
	}

	switch dc.Cmd {
	case CmdUp:
		d.Up(c.sig)
	case CmdDown:
		d.Down(c.sig)
	case CmdMy:
		d.My(c.sig)
	case CmdProg:
		d.Prog(c.sig)
	case CmdPosition:
		d.SetPosition(c.sig, dc.Pos)
	default:
		logrus.Warnf("error unknown command: %s", dc.Cmd)
		return
	}

	if err := c.save(); err != nil {
		logrus.Errorf("error saving device-config: %v", err)
	}
}

func (c *Controller) GetDevices() []*Device {
	return c.devices
}

func (c *Controller) GetDevice(id string) (*Device, error) {
	for _, d := range c.devices {
		if d.Id == id {
			return d, nil
		}
	}
	return nil, fmt.Errorf("did not find device with id %s", id)
}
