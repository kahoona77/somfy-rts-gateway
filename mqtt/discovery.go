package mqtt

import "somfyRtsGateway/core"

type discoveryConfig struct {
	Name                string   `json:"name"`
	UniqueId            string   `json:"unique_id"`
	DeviceClass         string   `json:"device_class"`
	CommandTopic        string   `json:"command_topic"`
	StateTopic          string   `json:"state_topic"`
	AvailabilityTopic   string   `json:"availability_topic"`
	PayloadAvailable    string   `json:"payload_available"`
	PayloadNotAvailable string   `json:"payload_not_available"`
	PayloadOpen         string   `json:"payload_open"`
	PayloadClose        string   `json:"payload_close"`
	PayloadStop         string   `json:"payload_stop"`
	StateOpen           string   `json:"state_open"`
	StateOpening        string   `json:"state_opening"`
	StateClosed         string   `json:"state_closed"`
	StateClosing        string   `json:"state_closing"`
	StateStopped        string   `json:"state_stopped"`
	Optimistic          bool     `json:"optimistic"`
	Device              haDevice `json:"device"`
}

type haDevice struct {
	Identifiers  []string `json:"identifiers"`
	Name         string   `json:"name"`
	Manufacturer string   `json:"manufacturer"`
	Model        string   `json:"model"`
}

func newDiscoveryConfig(device core.Device, baseTopic string) discoveryConfig {
	id := device.GetId()
	return discoveryConfig{
		Name:                device.GetName(),
		UniqueId:            "rollo_" + id,
		DeviceClass:         "shutter",
		CommandTopic:        baseTopic + "/" + id + "/set",
		StateTopic:          baseTopic + "/" + id + "/state",
		AvailabilityTopic:   baseTopic + "/bridge/availability",
		PayloadAvailable:    "online",
		PayloadNotAvailable: "offline",
		PayloadOpen:         "OPEN",
		PayloadClose:        "CLOSE",
		PayloadStop:         "STOP",
		StateOpen:           "open",
		StateOpening:        "opening",
		StateClosed:         "closed",
		StateClosing:        "closing",
		StateStopped:        "stopped",
		Optimistic:          true,
		Device: haDevice{
			Identifiers:  []string{"somfy-rts-gateway"},
			Name:         "Rollo-Steuerung",
			Manufacturer: "self-made",
			Model:        "go-signalduino-somfy",
		},
	}
}
