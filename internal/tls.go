package internal

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

var ErrInvalidCertData = errors.New("invalid certificate data")

func SetupTLSFile(cfg *tls.Config, file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	return SetupTLSData(cfg, data)
}

func SetupTLSData(cfg *tls.Config, data []byte) error {
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		return ErrInvalidCertData
	}

	cfg.InsecureSkipVerify = false
	cfg.RootCAs = certPool

	return nil
}
