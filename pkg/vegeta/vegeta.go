package vegeta

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"vegeta-server/models"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	vegeta "github.com/tsenart/vegeta/lib"
)

func tlsConfig(insecure bool, key, cert string, rootCerts []string) (*tls.Config, error) {
	c := tls.Config{InsecureSkipVerify: insecure} // nolint: gosec
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		log.WithError(err).Error("Vegeta TLS config failed")
		return nil, errors.Wrap(err, "Vegeta TLS config failed")
	}
	c.Certificates = append(c.Certificates, certificate)
	c.BuildNameToCertificate()

	if len(rootCerts) > 0 {
		c.RootCAs = x509.NewCertPool()
		for _, rootCert := range rootCerts {
			if !c.RootCAs.AppendCertsFromPEM([]byte(rootCert)) {
				log.WithError(err).Error("Vegeta TLS config failed")
				return nil, errors.Wrap(err, "Vegeta TLS config failed")
			}
		}
	}
	return &c, nil
}

func attackWithOpts(opts *AttackOpts) (*vegeta.Attacker, <-chan *vegeta.Result) {
	var c *tls.Config

	if opts.Cert != "" && opts.Key != "" {
		tlsConfig, err := tlsConfig(opts.Insecure, opts.Key, opts.Cert, opts.RootCerts)
		if err != nil {
			return nil, nil
		}
		c = tlsConfig
	}

	atk := vegeta.NewAttacker(
		vegeta.Redirects(opts.Redirects),
		vegeta.Timeout(opts.Timeout),
		vegeta.Workers(opts.Workers),
		vegeta.KeepAlive(opts.Keepalive),
		vegeta.Connections(opts.Connections),
		vegeta.HTTP2(opts.HTTP2),
		vegeta.H2C(opts.H2c),
		vegeta.MaxBody(opts.MaxBody),
		vegeta.TLSConfig(c),
		vegeta.LocalAddr(*opts.Laddr.IPAddr),
	)

	tr := vegeta.NewStaticTargeter(opts.Target)

	return atk, atk.Attack(tr, opts.Rate, opts.Duration, opts.Name)
}

// Attack implements the AttackFunc type for a vegeta based attacker
func Attack(name string, params models.AttackParams, quit chan struct{}) (io.Reader, error) {
	opts, err := NewAttackOptsFromAttackParams(name, params)
	if err != nil {
		log.WithError(err).Error("vegeta attack failed")
		return nil, errors.Wrap(err, "vegeta attack failed")
	}

	atk, result := attackWithOpts(opts)
	if result == nil {
		err := fmt.Errorf("empty channel returned")
		log.WithError(err).Error("vegeta attack failed")
		return nil, errors.Wrap(err, "vegeta attack failed")
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
				log.WithError(err).Error("Vegeta attack failed")
				return nil, errors.Wrap(err, "failed to encode result, vegeta attack failed")
			}
		case <-quit:
			atk.Stop()
			return nil, nil
		}
	}

	return buf, nil
}
