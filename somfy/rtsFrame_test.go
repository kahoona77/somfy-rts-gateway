package somfy

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestGetFrame(t *testing.T) {
	device := &Device{
		RollingCode:   172,
		Address:       3,
		EncryptionKey: 173,
		Id:            "Oben3",
		Name:          "Oben3",
	}

	frame := GetFrame(device, ButtonDown)
	expected := "ADEBEB47444444"
	actual := strings.ToUpper(hex.EncodeToString(frame))
	fmt.Printf("frame: '%s'", actual)
	if expected != actual {
		t.Errorf("Frame not correct")
	}
}
