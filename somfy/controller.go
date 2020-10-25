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

func NewController(ctx *core.Ctx) (*Controller, error) {
	s, err := signalduino.Open(ctx.Config.SignalduinoAddress)
	if err != nil {
		return nil, fmt.Errorf("error opening signalduino: %v", err)
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

		btn, err := getButtonFromCommand(dc.Cmd)
		if err != nil {
			logrus.Warnf("error: %v", err)
			continue
		}

		c.SendCmd(dc.Device, btn)
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

func (c *Controller) SendCmd(device string, btn Button) {
	for _, d := range c.devices {
		if d.Id == device {
			d.Send(c.sig, btn)
			if err := c.save(); err != nil {
				logrus.Errorf("error saving device-config: %v", err)
			}
			return
		}
	}
	logrus.Warnf("did not find device with id %s", device)
}

func (c *Controller) GetDevices() []*Device {
	return c.devices
}
