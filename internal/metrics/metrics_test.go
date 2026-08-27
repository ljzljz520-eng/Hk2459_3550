package metrics

import "testing"

func TestCounter(t *testing.T) {
	c := New()
	c.Inc("x")
	if c.Get("x") != 1 {
		t.Fatal("counter")
	}
}
