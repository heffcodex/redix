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

func NewClient(config Config) (*Client, error) {
	opts, err := redis.ParseURL(config.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse DSN: %w", err)
	}

	if config.Name != "" {
		opts.ClientName = config.Name
	}

	if opts.TLSConfig != nil {
		if err = config.Cert.setupTLS(opts.TLSConfig); err != nil {
			return nil, fmt.Errorf("setup TLS: %w", err)
		}
	}

	return &Client{
		UniversalClient: redis.NewClient(opts),
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
