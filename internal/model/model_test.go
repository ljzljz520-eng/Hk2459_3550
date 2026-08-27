package model

import "testing"

func TestChannelRules(t *testing.T) {
	c := NewChannel("x", Card)
	if !c.CanAccept(10) || c.Fee(100) != 0 {
		t.Fatal("rules")
	}
	c.Enabled = false
	if c.CanAccept(10) {
		t.Fatal("disabled")
	}
}
