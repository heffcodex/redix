package redix

import (
	"crypto/tls"
	"os"

	"github.com/heffcodex/redix/internal"
)

type Config struct {
	Name      string     `json:"name,omitempty" yaml:"name,omitempty" mapstructure:"name,omitempty"`
	Namespace Namespace  `json:"namespace,omitempty" yaml:"namespace,omitempty" mapstructure:"namespace,omitempty"`
	DSN       string     `json:"dsn" yaml:"dsn" mapstructure:"dsn"`
	Cert      ConfigCert `json:"cert,omitempty" yaml:"cert,omitempty" mapstructure:"cert,omitempty"`
}

type ConfigCert struct {
	Env  string `json:"env,omitempty" yaml:"env,omitempty" mapstructure:"env,omitempty"`
	File string `json:"file,omitempty" yaml:"file,omitempty" mapstructure:"file,omitempty"`
	Data []byte `json:"data,omitempty" yaml:"data,omitempty" mapstructure:"data,omitempty"`
}

func (c *ConfigCert) setupTLS(cfg *tls.Config) error {
	if certFile, ok := os.LookupEnv(c.Env); ok {
		return internal.SetupTLSFile(cfg, certFile)
	}

	if certFile := c.File; certFile != "" {
		return internal.SetupTLSFile(cfg, certFile)
	}

	if len(c.Data) > 0 {
		return internal.SetupTLSData(cfg, c.Data)
	}

	cfg.InsecureSkipVerify = true

	return nil
}
