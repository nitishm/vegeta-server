package vegeta

import (
	"encoding/base64"
	"net"
	"net/http"
	"strings"
	"time"
	"vegeta-server/models"

	"fmt"

	"github.com/pkg/errors"
	vegeta "github.com/tsenart/vegeta/lib"
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

// NewAttackOptsFromAttackParams adapts the models AttackParams to the vegeta specific options.
func NewAttackOptsFromAttackParams(name string, params models.AttackParams) (*AttackOpts, error) {
	rate := vegeta.Rate{Freq: params.Rate, Per: time.Second}

	// Set Duration
	dur, err := time.ParseDuration(params.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse duration")
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

	// Set local address
	laddr, err := net.ResolveIPAddr("ip", params.Laddr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to resolve IP address: %s", params.Laddr))
	}

	bBody, err := base64.StdEncoding.DecodeString(params.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode params.Body")
	}

	// Set Target
	tgt := vegeta.Target{
		Method: params.Target.Method,
		URL:    params.Target.URL,
		Header: hdr,
		Body:   bBody,
	}

	opts := &AttackOpts{
		Name:      name,
		Target:    tgt,
		Duration:  dur,
		Timeout:   timeout,
		Rate:      rate,
		Redirects: int(params.Redirects),
		MaxBody:   params.MaxBody,
		Keepalive: params.Keepalive,
		Resolvers: resolvers,
		Laddr:     struct{ *net.IPAddr }{laddr},
		Cert:      params.Cert,
		Key:       params.Key,
		RootCerts: params.RootCerts,
		Insecure:  params.Insecure,
		HTTP2:     params.HTTP2,
		H2c:       params.H2c,
		Workers:   uint64(params.Workers),
	}

	return opts, nil
}
