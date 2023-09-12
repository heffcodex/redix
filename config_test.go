package redix

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/heffcodex/testcerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNamespace_String(t *testing.T) {
	require.Equal(t, "foo", Namespace("foo").String())
}

func TestNamespace_Append(t *testing.T) {
	assert.Equal(t, Namespace(""), Namespace("").Append(""))
	assert.Equal(t, Namespace("foo"), Namespace("").Append("foo"))
	assert.Equal(t, Namespace("foo:"), Namespace("foo").Append(""))
	assert.Equal(t, Namespace("foo:bar"), Namespace("foo").Append("bar"))
}

func TestClientConfigCert_setupTLS(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		c := ClientConfigCert{}
		require.NoError(t, c.setupTLS(nil))
	})

	t.Run("no cert", func(t *testing.T) {
		c := ClientConfigCert{}
		tc := &tls.Config{}

		require.NoError(t, c.setupTLS(tc))
		assert.True(t, tc.InsecureSkipVerify)
		assert.Empty(t, tc.RootCAs)
	})

	t.Run("env", func(t *testing.T) {
		const envKey = "TEST_CERT_ENV"

		c := ClientConfigCert{Env: envKey}

		t.Run("empty", func(t *testing.T) {
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.True(t, tc.InsecureSkipVerify)
			assert.Empty(t, tc.RootCAs)
		})

		t.Run("wrong file", func(t *testing.T) {
			tc := &tls.Config{}

			require.NoError(t, os.Setenv(envKey, "wrong"))
			require.Error(t, c.setupTLS(tc))
		})

		t.Run("ok", func(t *testing.T) {
			cert, key, err := testcerts.GenerateCertsToTempFile(os.TempDir())
			require.NoError(t, err)

			defer func() {
				_ = os.Remove(cert)
				_ = os.Remove(key)
			}()

			tc := &tls.Config{}

			require.NoError(t, os.Setenv(envKey, cert))
			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})

	t.Run("file", func(t *testing.T) {
		t.Run("wrong filename", func(t *testing.T) {
			c := ClientConfigCert{File: "wrong"}
			tc := &tls.Config{}

			require.Error(t, c.setupTLS(tc))
		})

		t.Run("wrong contents", func(t *testing.T) {
			filename := os.TempDir() + "/wtf.crt"

			err := os.WriteFile(filename, []byte("wrong"), 0600)
			require.NoError(t, err)

			defer func() { _ = os.Remove(filename) }()

			c := ClientConfigCert{File: filename}
			tc := &tls.Config{}

			require.Error(t, c.setupTLS(tc))
		})

		t.Run("ok", func(t *testing.T) {
			cert, key, err := testcerts.GenerateCertsToTempFile(os.TempDir())
			require.NoError(t, err)

			defer func() {
				_ = os.Remove(cert)
				_ = os.Remove(key)
			}()

			c := ClientConfigCert{File: cert}
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})

	t.Run("data", func(t *testing.T) {
		t.Run("wrong data", func(t *testing.T) {
			c := ClientConfigCert{Data: []byte("wrong")}
			tc := &tls.Config{}

			require.Error(t, c.setupTLS(tc))
		})

		t.Run("ok", func(t *testing.T) {
			cert, _, err := testcerts.GenerateCerts()
			require.NoError(t, err)

			c := ClientConfigCert{Data: cert}
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})
}
