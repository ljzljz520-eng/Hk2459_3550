package model

import "errors"

func ValidateChannel(c Channel) error {
	if c.ID == "" || c.Name == "" {
		return errors.New("id and name required")
	}
	if c.LimitCents < 0 {
		return errors.New("limit must be positive")
	}
	if c.FeeBps < 0 || c.FeeBps > 10000 {
		return errors.New("fee out of range")
	}
	return nil
}
func NormalizeNotice(n string) string {
	if len(n) > 240 {
		return n[:240]
	}
	return n
}
func IsRiskyChange(old, next Channel) bool {
	return old.Enabled != next.Enabled || old.LimitCents != next.LimitCents
}
