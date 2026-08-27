package config

import "os"

type Config struct {
	DBPath          string
	Addr            string
	CacheTTLSeconds int
}

func Load() Config {
	p := os.Getenv("PAYMENT_DB")
	if p == "" {
		p = "payment.db"
	}
	a := os.Getenv("PAYMENT_ADDR")
	if a == "" {
		a = ":8080"
	}
	return Config{DBPath: p, Addr: a, CacheTTLSeconds: 30}
}
func (c Config) Validate() error {
	if c.DBPath == "" || c.Addr == "" {
		return os.ErrInvalid
	}
	return nil
}
