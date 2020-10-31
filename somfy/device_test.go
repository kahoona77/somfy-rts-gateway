package somfy

import (
	"testing"
)

func TestDevice_calcDuration(t *testing.T) {
	d := &Device{ClosingDuration: 23}
	r := d.calcDuration(50)
	if r != 11000 {
		t.Errorf("duration is not 11000 but %d", r)
	}
}
