package policy

import (
	"paymentconsole/internal/model"
	"testing"
)

func TestPolicyPreview(t *testing.T) {
	p := New()
	a := model.NewChannel("a", model.Card)
	b := a
	b.Enabled = false
	if !p.RequiresPreview(a, b) {
		t.Fatal("preview")
	}
}
