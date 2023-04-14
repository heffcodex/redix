package redix

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

const (
	keyDelimiter = ":"
)

var (
	_ redis.UniversalClient = (*XClient)(nil)
	_ Client                = (*XClient)(nil)
)

type Client interface {
	NamespacePrefix() string
	Key(parts ...string) string
}

type XClient struct {
	redis.UniversalClient
	ns string
}

func NewClient(config Config) (*XClient, error) {
	opts, err := redis.ParseURL(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	if opts.TLSConfig != nil {
		opts.TLSConfig.InsecureSkipVerify = true

		if config.Cert != "" {
			pemBytes, err := os.ReadFile(config.Cert)
			if err != nil {
				return nil, fmt.Errorf("read cert: %w", err)
			}

			certPool := x509.NewCertPool()
			if !certPool.AppendCertsFromPEM(pemBytes) {
				return nil, errors.New("fill cert pool")
			}

			opts.TLSConfig.InsecureSkipVerify = false
			opts.TLSConfig.RootCAs = certPool
		}
	}

	client := redis.NewClient(opts)
	wrapped := WrapClient(client, config.XConfig)

	return wrapped, nil
}

func WrapClient(c redis.UniversalClient, extraConfig XConfig) *XClient {
	return &XClient{
		UniversalClient: c,
		ns:              extraConfig.Namespace,
	}
}

func (c *XClient) NamespacePrefix() string {
	return c.key("")
}

func (c *XClient) Key(parts ...string) string {
	return c.key(parts...)
}

func (c *XClient) key(parts ...string) string {
	return strings.Join(append([]string{c.ns}, parts...), keyDelimiter)
}
