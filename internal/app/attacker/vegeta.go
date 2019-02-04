package attacker

import (
	"github.com/tsenart/vegeta/lib"
	"net"
	"time"
)

type AttackFunc func(*AttackOpts) <-chan *vegeta.Result

func DefaultAttackFn(opts *AttackOpts) <-chan *vegeta.Result {
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

type localAddr struct{ *net.IPAddr }
type csl []string

// AttackOpts aggregates the attack function command options
type AttackOpts struct {
	Target      vegeta.Target
	Name        string
	Body        string
	Cert        string
	Key         string
	RootCerts   csl
	HTTP2       bool
	H2c         bool
	Insecure    bool
	Duration    time.Duration
	Timeout     time.Duration
	Rate        vegeta.Rate
	Workers     uint64
	Connections int
	Redirects   int
	MaxBody     int64
	Laddr       localAddr
	Keepalive   bool
	Resolvers   csl
}
