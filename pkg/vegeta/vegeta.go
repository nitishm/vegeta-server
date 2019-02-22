package vegeta

import (
	"bytes"
	"fmt"
	"io"
	"vegeta-server/models"

	vegeta "github.com/tsenart/vegeta/lib"
)

func attackWithOpts(opts *AttackOpts) <-chan *vegeta.Result {
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

// Attack implements the AttackFunc type for a vegeta based attacker
func Attack(name string, params models.AttackParams, quit chan struct{}) (io.Reader, error) {
	opts, err := NewAttackOptsFromAttackParams(name, params)
	if err != nil {
		return nil, err
	}

	result := attackWithOpts(opts)
	if result == nil {
		return nil, fmt.Errorf("empty channel returned")
	}

	buf := bytes.NewBuffer(nil)
	enc := vegeta.NewEncoder(buf)
loop:
	for {
		select {
		case r, ok := <-result:
			if !ok {
				break loop
			}
			if err := enc.Encode(r); err != nil {
				return nil, err
			}
		case <-quit:
			return nil, nil
		}
	}

	return buf, nil
}
