package somfy

import (
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/signalduino"
	"strings"
	"time"
)

type Button uint16

const ButtonMy = 0x1
const ButtonUp = 0x2
const ButtonDown = 0x4
const ButtonProg = 0x8

const PosUp = 100
const PosDown = 0

type UpdateFunc func(d *Device)

type Device struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Position int    `json:"position"`

	Address         uint32 `json:"-"`
	RollingCode     uint16 `json:"-"`
	EncryptionKey   byte   `json:"-"`
	ClosingDuration int    `json:"-"`

	updateFuncs []UpdateFunc
}

func (d *Device) GetId() string {
	return d.Id
}

func (d *Device) GetName() string {
	return d.Name
}

func (d *Device) GetPosition() int {
	return d.Position
}

func (d *Device) GetAddress() uint32 {
	return d.Address
}

func (d *Device) GetRollingCode() uint16 {
	return d.RollingCode
}

func (d *Device) GetEncryptionKey() byte {
	return d.EncryptionKey
}

func (d *Device) GetClosingDuration() int {
	return d.ClosingDuration
}

func (d *Device) Down(sig *signalduino.Signalduino) {
	d.send(sig, ButtonDown)
	d.Position = PosDown
	d.update()
}

func (d *Device) Up(sig *signalduino.Signalduino) {
	d.send(sig, ButtonUp)
	d.Position = PosUp
	d.update()
}

func (d *Device) My(sig *signalduino.Signalduino) {
	d.send(sig, ButtonMy)
	d.Position = 50
	d.update()
}

func (d *Device) Prog(sig *signalduino.Signalduino) {
	d.send(sig, ButtonProg)
	d.Position = 50
	d.update()
}

func (d *Device) SetPosition(sig *signalduino.Signalduino, pos int) {
	if pos >= PosUp {
		d.Up(sig)
		return
	}

	if pos <= PosDown {
		d.Down(sig)
		return
	}

	if pos == d.Position {
		return
	}

	delta := pos - d.Position

	var button Button = ButtonUp
	if delta < 0 {
		button = ButtonDown
		delta = delta * -1
	}

	duration := d.calcDuration(delta)

	logrus.Infof("[Device %s] setPos: from %d to %d => delta: %d; direction: %d; duration: %d", d.Id, d.Position, pos, delta, button, duration)

	if duration > 1500 {
		d.send(sig, button)
		logrus.Infof("before sleep")
		time.Sleep(time.Millisecond * duration)
		logrus.Infof("after sleep")
		d.send(sig, ButtonMy)
	} else {
		logrus.Info("position-change to small -> doing nothing")
	}

	// set end Position
	d.Position = pos
	d.update()
}

func (d *Device) calcDuration(delta int) time.Duration {
	return time.Duration(float32(delta)*(float32(d.ClosingDuration)/float32(100))) * 1000
}

func (d *Device) update() {
	for _, fn := range d.updateFuncs {
		fn(d)
	}
}

func (d *Device) OnUpdate(fn UpdateFunc) {
	d.updateFuncs = append(d.updateFuncs, fn)
}

func (d *Device) send(sig *signalduino.Signalduino, btn Button) {
	//increase rollingCode
	d.RollingCode += 1
	d.EncryptionKey += 1
	if d.EncryptionKey > 0xAF {
		d.EncryptionKey = 0xA0
	}

	repetition := 6
	frame := GetFrame(d, btn)
	var cmd = fmt.Sprintf("SC;R=%d;SR;P0=-2560;P1=2560;P3=-640;D=10101010101010113;SM;C=645;D=%s;F=10AB85550A;",
		repetition, strings.ToUpper(hex.EncodeToString(frame)))
	logrus.Infof("SEND: %s", cmd)
	sig.Send(cmd)
	// wait a bit so that the command is processed
	time.Sleep(time.Millisecond * 1000)
}
