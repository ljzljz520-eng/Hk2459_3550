package store

import (
	"paymentconsole/internal/model"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := t.TempDir() + "/db"
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	c := model.NewChannel("card", model.Card)
	if e = s.SaveChannel(c); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if _, e = s.GetChannel("card"); e != nil {
		t.Fatal(e)
	}
}
func TestStoreCounts(t *testing.T) {
	s, _ := Open(t.TempDir() + "/db")
	defer s.Close()
	s.SaveRecord(model.Record{ID: "r"})
	if s.Count(recordBucket) != 1 {
		t.Fatal("count")
	}
}
