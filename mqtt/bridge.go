package mqtt

import (
	"encoding/json"
	"fmt"
	"strings"

	pahomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"

	"somfyRtsGateway/core"
	"somfyRtsGateway/somfy"
)

type Bridge struct {
	client  pahomqtt.Client
	cfg     *core.AppConfig
	cmdChan chan<- core.DeviceCmd
	devices []core.Device
	states  map[string]*deviceState
}

func NewBridge(cfg *core.AppConfig, cmdChan chan<- core.DeviceCmd, devices []core.Device) *Bridge {
	b := &Bridge{
		cfg:     cfg,
		cmdChan: cmdChan,
		devices: devices,
		states:  make(map[string]*deviceState),
	}

	for _, d := range devices {
		id := d.GetId()
		b.states[id] = newDeviceState(func(state string) {
			b.publishState(id, state)
		})
	}

	return b
}

func (b *Bridge) Start() error {
	opts := pahomqtt.NewClientOptions()
	opts.AddBroker(b.cfg.MqttBroker)
	opts.SetClientID(b.cfg.MqttClientId)

	if b.cfg.MqttUsername != "" {
		opts.SetUsername(b.cfg.MqttUsername)
		opts.SetPassword(b.cfg.MqttPassword)
	}

	availabilityTopic := b.cfg.MqttBaseTopic + "/bridge/availability"
	opts.SetWill(availabilityTopic, "offline", 1, true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(b.onConnect)
	opts.SetConnectionLostHandler(func(_ pahomqtt.Client, err error) {
		logrus.Warnf("MQTT connection lost: %v", err)
	})

	b.client = pahomqtt.NewClient(opts)

	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT connect failed: %w", token.Error())
	}

	logrus.Info("MQTT bridge started")
	return nil
}

func (b *Bridge) Stop() {
	if b.client == nil || !b.client.IsConnected() {
		return
	}
	availabilityTopic := b.cfg.MqttBaseTopic + "/bridge/availability"
	b.client.Publish(availabilityTopic, 1, true, "offline").Wait()
	b.client.Disconnect(500)
	logrus.Info("MQTT bridge stopped")
}

func (b *Bridge) onConnect(client pahomqtt.Client) {
	logrus.Info("MQTT connected")

	availabilityTopic := b.cfg.MqttBaseTopic + "/bridge/availability"
	client.Publish(availabilityTopic, 1, true, "online")

	for _, d := range b.devices {
		b.publishDiscovery(client, d)

		setTopic := b.cfg.MqttBaseTopic + "/" + d.GetId() + "/set"
		client.Subscribe(setTopic, 1, b.onMessage)
	}
}

func (b *Bridge) publishDiscovery(client pahomqtt.Client, d core.Device) {
	cfg := newDiscoveryConfig(d, b.cfg.MqttBaseTopic)
	payload, err := json.Marshal(cfg)
	if err != nil {
		logrus.Errorf("MQTT: failed to marshal discovery config for %s: %v", d.GetId(), err)
		return
	}
	topic := fmt.Sprintf("%s/cover/%s/config", b.cfg.MqttDiscoveryPrefix, d.GetId())
	client.Publish(topic, 1, true, payload)
}

func (b *Bridge) publishState(id, state string) {
	if b.client == nil || !b.client.IsConnected() {
		return
	}
	topic := b.cfg.MqttBaseTopic + "/" + id + "/state"
	b.client.Publish(topic, 1, false, state)
	logrus.Debugf("MQTT state %s → %s", id, state)
}

func (b *Bridge) onMessage(_ pahomqtt.Client, msg pahomqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) < 3 {
		return
	}
	id := parts[len(parts)-2]
	payload := string(msg.Payload())

	d := b.deviceById(id)
	if d == nil {
		logrus.Warnf("MQTT: received command for unknown device %s", id)
		return
	}

	state, ok := b.states[id]
	if !ok {
		return
	}

	logrus.Infof("MQTT command %s → %s", id, payload)

	switch payload {
	case "OPEN":
		b.cmdChan <- core.DeviceCmd{Device: id, Cmd: somfy.CmdUp}
		state.OnOpen(d.GetClosingDuration())
	case "CLOSE":
		b.cmdChan <- core.DeviceCmd{Device: id, Cmd: somfy.CmdDown}
		state.OnClose(d.GetClosingDuration())
	case "STOP":
		b.cmdChan <- core.DeviceCmd{Device: id, Cmd: somfy.CmdMy}
		state.OnStop()
	default:
		logrus.Warnf("MQTT: unknown payload %q for device %s", payload, id)
	}
}

func (b *Bridge) deviceById(id string) core.Device {
	for _, d := range b.devices {
		if d.GetId() == id {
			return d
		}
	}
	return nil
}
