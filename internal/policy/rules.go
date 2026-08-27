package policy

import "paymentconsole/internal/model"

type Rule struct {
	Name           string
	Type           model.ChannelType
	Min            int64
	Max            int64
	RequiresNotice bool
}

var Rules = []Rule{
	{Name: "card-basic", Type: model.Card, Min: 1, Max: 1000000, RequiresNotice: false},
	{Name: "wallet-basic", Type: model.Wallet, Min: 1, Max: 500000, RequiresNotice: false},
	{Name: "installment-review", Type: model.Installment, Min: 100, Max: 300000, RequiresNotice: true},
	{Name: "offline-review", Type: model.Offline, Min: 100, Max: 200000, RequiresNotice: true},
	{Name: "card-high", Type: model.Card, Min: 1000001, Max: 10000000, RequiresNotice: true},
	{Name: "wallet-high", Type: model.Wallet, Min: 500001, Max: 5000000, RequiresNotice: true},
	{Name: "installment-high", Type: model.Installment, Min: 300001, Max: 3000000, RequiresNotice: true},
	{Name: "offline-high", Type: model.Offline, Min: 200001, Max: 2000000, RequiresNotice: true},
}

func MatchRule(c model.Channel, amount int64) Rule {
	for _, r := range Rules {
		if r.Type == c.Type && amount >= r.Min && amount <= r.Max {
			return r
		}
	}
	return Rule{Name: "unmatched", Type: c.Type}
}
func NeedsNotice(c model.Channel, amount int64) bool { return MatchRule(c, amount).RequiresNotice }
func RuleNames() []string {
	out := make([]string, 0, len(Rules))
	for _, r := range Rules {
		out = append(out, r.Name)
	}
	return out
}
