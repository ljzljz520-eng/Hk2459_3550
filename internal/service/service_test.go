package service

import (
	"context"
	"os"
	"paymentconsole/internal/model"
	"paymentconsole/internal/store"
	"testing"
)

func testSvc(t *testing.T) (*Service, string) {
	p := t.TempDir() + "/x.db"
	st, e := store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	return New(st), p
}
func TestBusiness04Regression(t *testing.T) {
	s, p := testSvc(t)
	defer os.Remove(p)
	c := model.NewChannel("card", model.Card)
	if e := s.CreateChannel(context.Background(), c); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Pay(context.Background(), model.PaymentRequest{ID: "1", ChannelID: "card", AmountCents: 100}); e != nil {
		t.Fatal(e)
	}
	if e := s.UpdateChannel(context.Background(), "card", false, 0, 0, "", "test"); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Pay(context.Background(), model.PaymentRequest{ID: "2", ChannelID: "card", AmountCents: 100}); e == nil {
		t.Fatal("closed channel accepted")
	}
}
func TestCreateChannel(t *testing.T) {
	s, _ := testSvc(t)
	if e := s.CreateChannel(context.Background(), model.NewChannel("wallet", model.Wallet)); e != nil {
		t.Fatal(e)
	}
}
func TestPreview(t *testing.T) {
	s, _ := testSvc(t)
	s.CreateChannel(context.Background(), model.NewChannel("x", model.Card))
	if _, e := s.Preview(context.Background(), "x", false); e != nil {
		t.Fatal(e)
	}
}
