package signalduino

import (
	"bufio"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/sirupsen/logrus"
	"io"
)

type Signalduino struct {
	port io.ReadWriteCloser
}

func Open() (*Signalduino, error) {
	options := serial.OpenOptions{
		PortName:              "COM3",
		BaudRate:              57600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       0,
		InterCharacterTimeout: 10,
		ParityMode:            serial.PARITY_NONE,
		RTSCTSFlowControl:     true,
		Rs485RxDuringTx:       true,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		return nil, fmt.Errorf("serial.Open: %v", err)
	}

	s := &Signalduino{port: port}
	//go s.logReads()

	return s, nil
}

func (s *Signalduino) logReads() {
	logrus.Info("reading from signalduino...")
	reader := bufio.NewReader(s.port)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				logrus.Errorf("Error reading from serial port: ", err)
			}
		}
		logrus.Infof("Rx: %s", string(line))
	}
}

func (s *Signalduino) Version() {
	b := []byte{80, 10}
	n, err := s.port.Write(b)
	if err != nil {
		logrus.Errorf("port.Write: %v", err)
	}

	logrus.Infof("Wrote %d bytes.", n)
}

func (s *Signalduino) Close() error {
	return s.port.Close()
}

func (s *Signalduino) Send(cmd string) {
	cmd = cmd + "\n"
	b := []byte(cmd)
	n, err := s.port.Write(b)
	if err != nil {
		logrus.Errorf("port.Write: %v", err)
	}
	logrus.Infof("Wrote %d bytes.", n)
}
