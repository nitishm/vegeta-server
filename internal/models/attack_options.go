package models

import (
	"github.com/tsenart/vegeta/lib"
	"net"
	"net/http"
	"strings"
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

func NewAttackOptsFromAttackParams(name string, params AttackParams) (*AttackOpts, error) {
	rate := vegeta.Rate{Freq: int(params.Rate), Per: time.Second}

	// Set Duration
	dur, err := time.ParseDuration(params.Duration)
	if err != nil {
		return nil, err
	}

	// Set timeout
	timeout, _ := time.ParseDuration(params.Timeout)

	// Set target headers
	var hdr http.Header
	for _, h := range params.Headers {
		hdr.Add(h.Key, h.Value)
	}

	// Set resolvers
	resolvers := strings.Split(params.Resolvers, ",")

	// TODO: Set Local Address

	// TODO: Set TLS configuration

	// TODO: Set target body

	// Set Target
	tgt := vegeta.Target{
		Method: params.Target.Method,
		URL:    params.Target.URL,
		Header: hdr,
	}

	opts := &AttackOpts{
		Name:      name,
		Target:    tgt,
		Insecure:  params.Insecure,
		Duration:  dur,
		Timeout:   timeout,
		Rate:      rate,
		Redirects: int(params.Redirects),
		MaxBody:   params.MaxBody,
		Keepalive: params.Keepalive,
		Resolvers: resolvers,
	}
	opts.HTTP2 = params.Http2
	opts.H2c = params.H2c
	opts.Workers = uint64(params.Workers)

	return opts, nil
}
