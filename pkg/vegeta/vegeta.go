package vegeta

import (
	vegeta "github.com/tsenart/vegeta/lib"
)

func AttackFn(opts *AttackOpts) <-chan *vegeta.Result {
	atk := vegeta.NewAttacker(
		vegeta.Redirects(opts.Redirects),
		vegeta.Timeout(opts.Timeout),
		vegeta.Workers(opts.Workers),
		vegeta.KeepAlive(opts.Keepalive),
		vegeta.Connections(opts.Connections),
		vegeta.HTTP2(opts.HTTP2),
		vegeta.H2C(opts.H2c),
		vegeta.MaxBody(opts.MaxBody),
	)
	tr := vegeta.NewStaticTargeter(opts.Target)
	return atk.Attack(tr, opts.Rate, opts.Duration, opts.Name)
}
