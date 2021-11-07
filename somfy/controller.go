package somfy

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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
	s, err := signalduino.Open(ctx.Config.SignalduinoAddress)
	if err != nil {
		return nil, fmt.Errorf("error opening signalduino on address '%s': %v", ctx.Config.SignalduinoAddress, err)
	}
	s.Version()

	devices, err := loadDevices(ctx.Config.DevicesFile)
	if err != nil {
		return nil, fmt.Errorf("error loading devices from '%s': %v", ctx.Config.DevicesFile, err)
	}

	ctrl := &Controller{
		sig:         s,
		devices:     devices,
		devicesFile: ctx.Config.DevicesFile,
	}
	go ctrl.listen(ctx.CommandChannel)

	return ctrl, nil
}

func (c *Controller) Close() {
	logrus.Debugf("closing controller...")
	if err := c.sig.Close(); err != nil {
		logrus.Errorf("error closing signalduino: %v", err)
	}
}

func loadDevices(file string) ([]*Device, error) {
	yamlFile, err := ioutil.ReadFile(file)
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

	// write to file
	f, err := os.Create(c.devicesFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.Name(), d, 0644)
	if err != nil {
		return err
	}

	path, _ := filepath.Abs(f.Name())

	logrus.Infof("saved config to: %s", path)

	return f.Close()
}

func (c *Controller) SendCmd(dc core.DeviceCmd) {
	d, err := c.GetDevice(dc.Device)
	if err != nil {
		logrus.Warn(err)
	}

	switch dc.Cmd {
	case CmdUp:
		d.Up(c.sig)
		break
	case CmdDown:
		d.Down(c.sig)
		break
	case CmdMy:
		d.My(c.sig)
		break
	case CmdProg:
		d.Prog(c.sig)
		break
	case CmdPosition:
		d.SetPosition(c.sig, dc.Pos)
		break
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
