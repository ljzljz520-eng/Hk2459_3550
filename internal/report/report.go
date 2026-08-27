package report

import (
	"fmt"
	"paymentconsole/internal/model"
)

func PaymentSummary(results []model.PaymentResult) string {
	ok := 0
	for _, r := range results {
		if r.Accepted {
			ok++
		}
	}
	return fmt.Sprintf("accepted=%d total=%d", ok, len(results))
}
func ChannelSummary(channels []model.Channel) map[string]int {
	m := map[string]int{}
	for _, c := range channels {
		if c.Enabled {
			m["enabled"]++
		} else {
			m["disabled"]++
		}
	}
	return m
}
