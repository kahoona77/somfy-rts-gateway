package signalduino

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
	"io"
	"time"
)

type Signalduino struct {
	port serial.Port
}

func Open(devicePort string) (*Signalduino, error) {
	mode := &serial.Mode{
		BaudRate: 57600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	// Open the port.
	port, err := serial.Open(devicePort, mode)
	if err != nil {
		return nil, fmt.Errorf("serial.Open: %v", err)
	}

	if err := port.SetDTR(true); err != nil {
		return nil, fmt.Errorf("error setting dtr: %v", err)
	}

	s := &Signalduino{port: port}
	s.logReads()

	return s, nil
}

func (s *Signalduino) logReads() {
	logrus.Info("start reading from signalduino...")
	reader := bufio.NewReader(s.port)
	go func() {
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				if err != io.EOF {
					logrus.Errorf("Error reading from serial port: ", err)
				}
			}
			logrus.Infof("Rx: %s", string(line))
		}
	}()
	time.Sleep(time.Second * 2)
}

func (s *Signalduino) Version() {
	s.Send("V")
}

func (s *Signalduino) Ping() {
	s.Send("P")
}

func (s *Signalduino) Close() error {
	return s.port.Close()
}

func (s *Signalduino) Send(cmd string) {
	cmd = cmd + "\n"
	b := []byte(cmd)
	_, err := s.port.Write(b)
	if err != nil {
		logrus.Errorf("port.Write: %v", err)
	}
}
