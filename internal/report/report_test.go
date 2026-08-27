package report

import (
	"paymentconsole/internal/model"
	"testing"
)

func TestSummary(t *testing.T) {
	if PaymentSummary([]model.PaymentResult{{Accepted: true}}) != "accepted=1 total=1" {
		t.Fatal("summary")
	}
}
