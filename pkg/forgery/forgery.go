package forgery

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

type Forgery struct {
	cert             *x509.Certificate
	key              interface{}
	serverPrivateKey *ecdsa.PrivateKey
	serverKeyBytes   []byte
}

func New(certFile, keyFile string) (*Forgery, error) {
	cert, key, err := loadCertKey(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	//using the same private key for all servers
	privateKey, keyBytes, err := generateServerPrivateKey()
	if err != nil {
		return nil, err
	}

	return &Forgery{cert: cert, key: key, serverPrivateKey: privateKey, serverKeyBytes: keyBytes}, nil
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
