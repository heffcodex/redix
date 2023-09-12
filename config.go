package redix

import (
	"crypto/tls"
	"os"
	"strings"

	"github.com/heffcodex/redix/internal"
)

const (
	KeyDelimiter = ":"
)

type Namespace string

func (ns Namespace) String() string {
	return string(ns)
}

func (ns Namespace) Append(parts ...string) Namespace {
	var slice []string

	if ns != "" {
		slice = []string{string(ns)}
	}

	return Namespace(strings.Join(append(slice, parts...), KeyDelimiter))
}

type Config struct {
	ClientConfig `mapstructure:",squash"` //nolint: tagliatelle // squash does not need any tag
	Namespace    Namespace                `json:"namespace" yaml:"namespace" mapstructure:"namespace"`
}

type ClientConfig struct {
	DSN  string           `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	Cert ClientConfigCert `json:"cert,omitempty" yaml:"cert,omitempty" mapstructure:"cert,omitempty"`
}

type ClientConfigCert struct {
	Env  string `json:"env,omitempty" yaml:"env,omitempty" mapstructure:"env,omitempty"`
	File string `json:"file,omitempty" yaml:"file,omitempty" mapstructure:"file,omitempty"`
	Data []byte `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data,omitempty"`
}

func (c *ClientConfigCert) setupTLS(cfg *tls.Config) error {
	if cfg == nil {
		return nil
	}

	cfg.InsecureSkipVerify = true

	switch {
	case c.Env != "":
		return internal.SetupTLSFile(cfg, os.Getenv(c.Env))
	case c.File != "":
		return internal.SetupTLSFile(cfg, c.File)
	case len(c.Data) > 0:
		return internal.SetupTLSData(cfg, c.Data)
	}

	return nil
}
