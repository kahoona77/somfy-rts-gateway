package mqtt

import (
	"testing"
	"time"
)

func assertState(t *testing.T, ch <-chan string, expected string) {
	t.Helper()
	select {
	case got := <-ch:
		if got != expected {
			t.Errorf("expected state %q, got %q", expected, got)
		}
	case <-time.After(200 * time.Millisecond):
		t.Errorf("timeout waiting for state %q", expected)
	}
}

func newTestState() (*deviceState, <-chan string) {
	ch := make(chan string, 10)
	ds := newDeviceState(func(state string) { ch <- state })
	return ds, ch
}

func TestOnOpen_PublishesOpeningThenOpen(t *testing.T) {
	ds, ch := newTestState()
	ds.OnOpen(0)
	assertState(t, ch, "opening")
	assertState(t, ch, "open")
}

func TestOnClose_PublishesClosingThenClosed(t *testing.T) {
	ds, ch := newTestState()
	ds.OnClose(0)
	assertState(t, ch, "closing")
	assertState(t, ch, "closed")
}

func TestOnStop_PublishesStopped(t *testing.T) {
	ds, ch := newTestState()
	ds.OnStop()
	assertState(t, ch, "stopped")
}

func TestOnStop_CancelsRunningOpenTimer(t *testing.T) {
	ds, ch := newTestState()
	ds.OnOpen(10) // 10 seconds — would not fire in this test
	assertState(t, ch, "opening")

	ds.OnStop()
	assertState(t, ch, "stopped")

	// "open" must not arrive after stop
	select {
	case got := <-ch:
		t.Errorf("unexpected state after stop: %q", got)
	case <-time.After(100 * time.Millisecond):
		// expected: nothing arrives
	}
}

func TestOnStop_CancelsRunningCloseTimer(t *testing.T) {
	ds, ch := newTestState()
	ds.OnClose(10)
	assertState(t, ch, "closing")

	ds.OnStop()
	assertState(t, ch, "stopped")

	select {
	case got := <-ch:
		t.Errorf("unexpected state after stop: %q", got)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestOnOpen_ReplacesRunningCloseTimer(t *testing.T) {
	ds, ch := newTestState()
	ds.OnClose(10)
	assertState(t, ch, "closing")

	// Reverse direction before timer fires
	ds.OnOpen(0)
	assertState(t, ch, "opening")
	assertState(t, ch, "open")

	// "closed" from the cancelled timer must not arrive
	select {
	case got := <-ch:
		t.Errorf("unexpected state after direction change: %q", got)
	case <-time.After(100 * time.Millisecond):
	}
}
