package redix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heffcodex/redix/internal"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	t.Run("dsn", func(t *testing.T) {
		t.Parallel()

		t.Run("empty", func(t *testing.T) {
			t.Parallel()

			_, err := NewClient(Config{DSN: ""})
			require.ErrorContains(t, err, "parse DSN")
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			_, err := NewClient(Config{DSN: "sh!t"})
			require.ErrorContains(t, err, "parse DSN")
		})

		t.Run("plain no cert read", func(t *testing.T) {
			t.Parallel()

			_, err := NewClient(Config{DSN: "redis://localhost:6379", Cert: ConfigCert{Data: []byte("wtf")}})
			require.NoError(t, err)
		})

		t.Run("tls", func(t *testing.T) {
			t.Parallel()

			t.Run("invalid", func(t *testing.T) {
				_, err := NewClient(Config{DSN: "rediss://localhost:6379", Cert: ConfigCert{Data: []byte("wtf")}})
				require.ErrorIs(t, err, internal.ErrInvalidCertData)
			})

			t.Run("valid", func(t *testing.T) {
				_, err := NewClient(Config{DSN: "rediss://localhost:6379"})
				require.NoError(t, err)
			})
		})
	})
}

func TestClient_Namespace(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Namespace("foo"), (&Client{ns: "foo"}).Namespace())
}

func TestClient_KeyPrefix(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", (&Client{}).KeyPrefix())
	assert.Equal(t, "foo:", (&Client{ns: "foo"}).KeyPrefix())
}

func TestClient_Key(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "", (&Client{ns: ""}).Key(""))
	assert.Equal(t, "bar", (&Client{ns: ""}).Key("bar"))
	assert.Equal(t, "foo:", (&Client{ns: "foo"}).Key(""))
	assert.Equal(t, "foo:bar", (&Client{ns: "foo"}).Key("bar"))
}
