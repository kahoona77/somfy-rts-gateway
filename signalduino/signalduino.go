package signalduino

import (
	"bufio"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type Signalduino struct {
	port io.ReadWriteCloser
}

func Open(devicePort string) (*Signalduino, error) {
	options := serial.OpenOptions{
		PortName:              devicePort,
		BaudRate:              57600,
		DataBits:              8,
		StopBits:              1,
		ParityMode:            0,
		RTSCTSFlowControl:     true,
		InterCharacterTimeout: 0,
		MinimumReadSize:       4,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return nil, fmt.Errorf("serial.Open: %v", err)
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
					logrus.Errorf("Error reading from serial port: %v", err)
				}
			}
			if len(line) > 0 {
				logrus.Infof("Rx: %s", string(line))
			}
		}
	}()
	time.Sleep(time.Second * 2)
}

func (s *Signalduino) Version() {
	logrus.Infof("SEND V (Version)")
	s.Send("V")
}

func (s *Signalduino) Ping() {
	logrus.Infof("SEND P (Ping)")
	s.Send("P")
}

func (s *Signalduino) Close() error {
	return s.port.Close()
}

func (s *Signalduino) Send(cmd string) {
	logrus.Infof("SEND cmd: %s", cmd)
	cmd = cmd + "\n"
	b := []byte(cmd)
	_, err := s.port.Write(b)
	if err != nil {
		logrus.Errorf("port.Write: %v", err)
	}
}
