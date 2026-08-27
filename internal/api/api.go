package api

import (
	"context"
	"encoding/json"
	"net/http"
	"paymentconsole/internal/model"
	"paymentconsole/internal/service"
	"paymentconsole/internal/store"
	"time"
)

type Server struct{ s *service.Service }

func New(s *service.Service) *Server { return &Server{s: s} }
func (a *Server) Routes() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/health", a.health)
	m.HandleFunc("/channels", a.channels)
	m.HandleFunc("/pay", a.pay)
	return m
}
func (a *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func (a *Server) channels(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := a.s.Snapshot()
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(v)
		return
	}
	var c model.Channel
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "bad json", 400)
		return
	}
	if e := a.s.CreateChannel(r.Context(), c); e != nil {
		http.Error(w, e.Error(), 400)
		return
	}
	w.WriteHeader(201)
}
func (a *Server) pay(w http.ResponseWriter, r *http.Request) {
	var p model.PaymentRequest
	if json.NewDecoder(r.Body).Decode(&p) != nil {
		http.Error(w, "bad json", 400)
		return
	}
	p.CreatedAt = now()
	v, e := a.s.Pay(r.Context(), p)
	if e != nil {
		http.Error(w, e.Error(), 409)
		return
	}
	json.NewEncoder(w).Encode(v)
}
func now() time.Time { return time.Now() }

var _ = context.Background
var _ = store.ID
