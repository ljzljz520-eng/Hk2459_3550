package report

import (
	"paymentconsole/internal/model"
	"sort"
)

type ChannelReport struct {
	ID       string
	Enabled  bool
	Payments int
	Volume   int64
	Failure  int
}

func BuildReports(ch []model.Channel, p []model.PaymentResult) []ChannelReport {
	m := map[string]*ChannelReport{}
	for _, c := range ch {
		m[c.ID] = &ChannelReport{ID: c.ID, Enabled: c.Enabled}
	}
	for _, x := range p {
		r := m[x.ChannelID]
		if r == nil {
			r = &ChannelReport{ID: x.ChannelID}
			m[x.ChannelID] = r
		}
		r.Payments++
		if x.Accepted {
			r.Volume += x.FeeCents
		} else {
			r.Failure++
		}
	}
	out := make([]ChannelReport, 0, len(m))
	for _, r := range m {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func FailureRate(r ChannelReport) float64 {
	if r.Payments == 0 {
		return 0
	}
	return float64(r.Failure) / float64(r.Payments)
}
func EnabledReports(rs []ChannelReport) []ChannelReport {
	out := []ChannelReport{}
	for _, r := range rs {
		if r.Enabled {
			out = append(out, r)
		}
	}
	return out
}
func TotalVolume(rs []ChannelReport) int64 {
	var n int64
	for _, r := range rs {
		n += r.Volume
	}
	return n
}
