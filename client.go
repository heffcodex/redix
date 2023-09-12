package redix

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var _ UniversalClient = (*Client)(nil)

type UniversalClient interface {
	redis.UniversalClient

	Namespace() Namespace
	Key(parts ...string) string
}

type Client struct {
	redis.UniversalClient
	ns Namespace
}

func NewClient(name string, config *Config) (*Client, error) {
	opts, err := redis.ParseURL(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	if name != "" {
		opts.ClientName = name
	}

	if err = config.Cert.setupTLS(opts.TLSConfig); err != nil {
		return nil, fmt.Errorf("setup TLS: %w", err)
	}

	rc := redis.NewClient(opts)

	return &Client{
		UniversalClient: rc,
		ns:              config.Namespace,
	}, nil
}

func (c *Client) Namespace() Namespace {
	return c.ns
}

func (c *Client) KeyPrefix() string {
	return c.Key("")
}

func (c *Client) Key(parts ...string) string {
	return c.ns.Append(parts...).String()
}
