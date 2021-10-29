package forging

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	mrand "math/rand"
	"os"
)

type Forging struct {
	cert *x509.Certificate
	key  interface{}

	serverPrivateKey *ecdsa.PrivateKey
	serverKeyBytes   []byte
}

func New(certFile, keyFile string) (*Forging, error) {
	cert, key, err := loadCertKey(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	//using the same private key for all servers
	privateKey, keyBytes, err := generateServerPrivateKey()
	if err != nil {
		return nil, err
	}

	return &Forging{cert: cert, key: key, serverPrivateKey: privateKey, serverKeyBytes: keyBytes}, nil
}

func (f *Forging) Forge(serverCert *x509.Certificate) (*tls.Certificate, error) {
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

func loadCertKey(certFile, keyFile string) (*x509.Certificate, interface{}, error) {
	certBytes, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("reading certFile: %w", err)
	}
	certDecoded, _ := pem.Decode(certBytes)
	if certDecoded == nil {
		return nil, nil, fmt.Errorf("decoding certBytes: %w", err)
	}
	cert, err := x509.ParseCertificate(certDecoded.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing certDecoded: %w", err)
	}

	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("reading keyFile: %w", err)
	}
	keyDecoded, _ := pem.Decode(keyBytes)
	if keyDecoded == nil {
		return nil, nil, fmt.Errorf("decoding keyBytes: %w", err)
	}
	key, err := x509.ParsePKCS8PrivateKey(keyDecoded.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing keyDecoded: %w", err)
	}

	return cert, key, err
}

func generateServerPrivateKey() (*ecdsa.PrivateKey, []byte, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	keyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, keyBytes, nil
}
