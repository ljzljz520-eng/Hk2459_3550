package main

import (
	"log"
	"net/http"
	"paymentconsole/internal/api"
	"paymentconsole/internal/config"
	"paymentconsole/internal/service"
	"paymentconsole/internal/store"
)

func main() {
	c := config.Load()
	st, e := store.Open(c.DBPath)
	if e != nil {
		log.Fatal(e)
	}
	defer st.Close()
	svc := service.New(st)
	if e := svc.EnsureDefaults(nil); e != nil {
		log.Fatal(e)
	}
	log.Fatal(http.ListenAndServe(c.Addr, api.New(svc).Routes()))
}
