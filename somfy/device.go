package somfy

import (
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"somfyRtsGateway/signalduino"
	"strings"
)

type Device struct {
	Name          string
	Address       uint32
	RollingCode   uint16
	EncryptionKey byte
}

func (d *Device) Send(sig *signalduino.Signalduino, btn Button) {
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
}
