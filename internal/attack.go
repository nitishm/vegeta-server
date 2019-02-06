package internal

import (
	"github.com/tsenart/vegeta/lib"
	"net"
	"time"
)

// AttackOpts aggregates the attack function command options
type AttackOpts struct {
	Target      vegeta.Target
	Name        string
	Body        string
	Cert        string
	Key         string
	RootCerts   []string
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
	Laddr       struct{ *net.IPAddr }
	Keepalive   bool
	Resolvers   []string
}
