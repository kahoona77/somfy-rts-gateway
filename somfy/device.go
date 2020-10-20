package somfy

import (
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/signalduino"
)

type Device struct {
	RollingCode   uint16
	Address       uint32
	EncryptionKey byte
	Name          string
}

func (d *Device) Send(sig *signalduino.Signalduino, btn Button) {
	repetition := 6
	frame := GetFrame(d, btn)
	var cmd = fmt.Sprintf("SC;R=%d;SR;P0=-2560;P1=2560;P3=-640;D=10101010101010113;SM;C=645;D=%s;F=10AB85550A;", repetition, hex.EncodeToString(frame))
	logrus.Infof("SEND: %s", cmd)
	sig.Send(cmd)
}
