package vegeta

import (
	"net"
	"net/http"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type AttackFunc func(*AttackOpts) <-chan *vegeta.Result

func DefaultAttackFunc(opts *AttackOpts) <-chan *vegeta.Result {
	atk := vegeta.NewAttacker()
	tr := vegeta.NewStaticTargeter(opts.Target)
	return atk.Attack(tr, opts.Rate, opts.Duration, opts.Name)
}

type headers struct{ http.Header }
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
	Headers     headers
	Laddr       localAddr
	Keepalive   bool
	Resolvers   csl
}
