package core

import "somfyRtsGateway/signalduino"

type Controller interface {
	Signalduino() *signalduino.Signalduino
	Devices() []Device
}

type Device interface {
	GetId() string
	GetName() string
	GetPosition() int
	GetAddress() uint32
	GetRollingCode() uint16
	GetEncryptionKey() byte
	GetClosingDuration() int
}
