package mqtt

import (
	"testing"
)

type fakeDevice struct {
	id              string
	name            string
	closingDuration int
}

func (f *fakeDevice) GetId() string           { return f.id }
func (f *fakeDevice) GetName() string         { return f.name }
func (f *fakeDevice) GetPosition() int        { return 0 }
func (f *fakeDevice) GetAddress() uint32      { return 0 }
func (f *fakeDevice) GetRollingCode() uint16  { return 0 }
func (f *fakeDevice) GetEncryptionKey() byte  { return 0 }
func (f *fakeDevice) GetClosingDuration() int { return f.closingDuration }

func TestNewDiscoveryConfig_Topics(t *testing.T) {
	d := &fakeDevice{id: "wohnzimmer", name: "Rollo Wohnzimmer"}
	cfg := newDiscoveryConfig(d, "covers")

	if cfg.CommandTopic != "covers/wohnzimmer/set" {
		t.Errorf("CommandTopic = %q", cfg.CommandTopic)
	}
	if cfg.StateTopic != "covers/wohnzimmer/state" {
		t.Errorf("StateTopic = %q", cfg.StateTopic)
	}
	if cfg.AvailabilityTopic != "covers/bridge/availability" {
		t.Errorf("AvailabilityTopic = %q", cfg.AvailabilityTopic)
	}
}

func TestNewDiscoveryConfig_Payloads(t *testing.T) {
	d := &fakeDevice{id: "o1", name: "Oben 1"}
	cfg := newDiscoveryConfig(d, "covers")

	checks := map[string]string{
		"PayloadOpen":         cfg.PayloadOpen,
		"PayloadClose":        cfg.PayloadClose,
		"PayloadStop":         cfg.PayloadStop,
		"PayloadAvailable":    cfg.PayloadAvailable,
		"PayloadNotAvailable": cfg.PayloadNotAvailable,
	}
	expected := map[string]string{
		"PayloadOpen":         "OPEN",
		"PayloadClose":        "CLOSE",
		"PayloadStop":         "STOP",
		"PayloadAvailable":    "online",
		"PayloadNotAvailable": "offline",
	}
	for field, got := range checks {
		if got != expected[field] {
			t.Errorf("%s = %q, want %q", field, got, expected[field])
		}
	}
}

func TestNewDiscoveryConfig_States(t *testing.T) {
	d := &fakeDevice{id: "o1", name: "Oben 1"}
	cfg := newDiscoveryConfig(d, "covers")

	checks := map[string]string{
		"StateOpen":    cfg.StateOpen,
		"StateOpening": cfg.StateOpening,
		"StateClosed":  cfg.StateClosed,
		"StateClosing": cfg.StateClosing,
		"StateStopped": cfg.StateStopped,
	}
	expected := map[string]string{
		"StateOpen":    "open",
		"StateOpening": "opening",
		"StateClosed":  "closed",
		"StateClosing": "closing",
		"StateStopped": "stopped",
	}
	for field, got := range checks {
		if got != expected[field] {
			t.Errorf("%s = %q, want %q", field, got, expected[field])
		}
	}
}

func TestNewDiscoveryConfig_Metadata(t *testing.T) {
	d := &fakeDevice{id: "u3", name: "Unten 3"}
	cfg := newDiscoveryConfig(d, "covers")

	if cfg.UniqueId != "rollo_u3" {
		t.Errorf("UniqueId = %q", cfg.UniqueId)
	}
	if cfg.Name != "Unten 3" {
		t.Errorf("Name = %q", cfg.Name)
	}
	if !cfg.Optimistic {
		t.Error("Optimistic should be true")
	}
	if cfg.DeviceClass != "shutter" {
		t.Errorf("DeviceClass = %q", cfg.DeviceClass)
	}
	if len(cfg.Device.Identifiers) == 0 || cfg.Device.Identifiers[0] != "somfy-rts-gateway" {
		t.Errorf("Device.Identifiers = %v", cfg.Device.Identifiers)
	}
}
