package api

import (
	"net/http/httptest"
	"paymentconsole/internal/service"
	"paymentconsole/internal/store"
	"testing"
)

func TestHealth(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	r := httptest.NewRecorder()
	q := httptest.NewRequest("GET", "/health", nil)
	New(service.New(s)).Routes().ServeHTTP(r, q)
	if r.Code != 200 {
		t.Fatal(r.Code)
	}
}
