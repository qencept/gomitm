package trusted

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

type Trusted struct {
	pool *x509.CertPool
}

func New(trustedRootCaCerts []string) (*Trusted, error) {
	trusted := &Trusted{}
	if len(trustedRootCaCerts) > 0 {
		trusted.pool = x509.NewCertPool()
		for _, certFile := range trustedRootCaCerts {
			certBytes, err := ioutil.ReadFile(certFile)
			if err != nil {
				return nil, fmt.Errorf("reading %s: %w", certFile, err)
			}
			trusted.pool.AppendCertsFromPEM(certBytes)
		}
	}
	return trusted, nil
}

func (t *Trusted) CertPool() *x509.CertPool {
	return t.pool
}
