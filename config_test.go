package redix

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/heffcodex/testcerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/heffcodex/redix/internal"
)

func TestClientConfigCert_setupTLS(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		require.Panics(t, func() { _ = new(ConfigCert).setupTLS(nil) })
	})

	t.Run("no cert", func(t *testing.T) {
		c := ConfigCert{}
		tc := &tls.Config{}

		require.NoError(t, c.setupTLS(tc))
		assert.True(t, tc.InsecureSkipVerify)
		assert.Empty(t, tc.RootCAs)
	})

	t.Run("env", func(t *testing.T) {
		const envKey = "TEST_CERT_ENV"

		c := ConfigCert{Env: envKey}

		t.Run("empty", func(t *testing.T) {
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.True(t, tc.InsecureSkipVerify)
			assert.Empty(t, tc.RootCAs)
		})

		t.Run("wrong file", func(t *testing.T) {
			tc := &tls.Config{}

			t.Setenv(envKey, "wrong")
			require.ErrorContains(t, c.setupTLS(tc), "read file")
		})

		t.Run("ok", func(t *testing.T) {
			certFile, _, err := testcerts.GenerateCertsToTempFile(os.TempDir())
			require.NoError(t, err)

			tc := &tls.Config{}

			t.Setenv(envKey, certFile)
			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})

	t.Run("file", func(t *testing.T) {
		t.Run("wrong filename", func(t *testing.T) {
			c := ConfigCert{File: "wrong"}
			tc := &tls.Config{}

			require.ErrorContains(t, c.setupTLS(tc), "read file")
		})

		t.Run("wrong contents", func(t *testing.T) {
			filename := os.TempDir() + "/wtf.crt"

			err := os.WriteFile(filename, []byte("wrong"), 0o600)
			require.NoError(t, err)

			c := ConfigCert{File: filename}
			tc := &tls.Config{}

			require.ErrorIs(t, c.setupTLS(tc), internal.ErrInvalidCertData)
		})

		t.Run("ok", func(t *testing.T) {
			cert, _, err := testcerts.GenerateCertsToTempFile(os.TempDir())
			require.NoError(t, err)

			c := ConfigCert{File: cert}
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})

	t.Run("data", func(t *testing.T) {
		t.Run("wrong data", func(t *testing.T) {
			c := ConfigCert{Data: []byte("wrong")}
			tc := &tls.Config{}

			require.ErrorIs(t, c.setupTLS(tc), internal.ErrInvalidCertData)
		})

		t.Run("ok", func(t *testing.T) {
			cert, _, err := testcerts.GenerateCerts()
			require.NoError(t, err)

			c := ConfigCert{Data: cert}
			tc := &tls.Config{}

			require.NoError(t, c.setupTLS(tc))
			assert.False(t, tc.InsecureSkipVerify)
			assert.NotEmpty(t, tc.RootCAs)
		})
	})
}
