package forgery

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	mrand "math/rand"
)

func (f *Forgery) Forge(serverCert *x509.Certificate) (*tls.Certificate, error) {
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(mrand.Int63()),
		Subject:               serverCert.Subject,
		DNSNames:              serverCert.DNSNames,
		NotBefore:             serverCert.NotBefore,
		NotAfter:              serverCert.NotAfter,
		KeyUsage:              serverCert.KeyUsage,
		ExtKeyUsage:           serverCert.ExtKeyUsage,
		IsCA:                  serverCert.IsCA,
		BasicConstraintsValid: serverCert.BasicConstraintsValid,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, f.cert, f.serverPrivateKey.Public(), f.key)
	if err != nil {
		return nil, fmt.Errorf("forging certificate: %w", err)
	}
	certPem := new(bytes.Buffer)
	if err := pem.Encode(certPem, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return nil, fmt.Errorf("encoding certificate: %w", err)
	}
	keyPem := new(bytes.Buffer)
	if err = pem.Encode(keyPem, &pem.Block{Type: "PRIVATE KEY", Bytes: f.serverKeyBytes}); err != nil {
		return nil, fmt.Errorf("encoding private key: %w", err)
	}
	cert, err := tls.X509KeyPair(certPem.Bytes(), keyPem.Bytes())
	if err != nil {
		return nil, fmt.Errorf("creating tls cert: %w", err)
	}
	return &cert, nil
}
