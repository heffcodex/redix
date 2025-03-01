package internal

import (
	"crypto/tls"
	"os"
	"testing"

	"github.com/heffcodex/testcerts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupTLSFile(t *testing.T) {
	t.Parallel()

	t.Run("empty filename", func(t *testing.T) {
		t.Parallel()
		require.ErrorContains(t, SetupTLSFile(new(tls.Config), ""), "read file")
	})

	t.Run("invalid filename", func(t *testing.T) {
		t.Parallel()
		require.ErrorContains(t, SetupTLSFile(new(tls.Config), "__no.file"), "read file")
	})

	t.Run("valid cert", func(t *testing.T) {
		t.Parallel()

		cfg := &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // test
		}

		certFile, _, err := testcerts.GenerateCertsToTempFile(os.TempDir())
		require.NoError(t, err)
		require.NoError(t, SetupTLSFile(cfg, certFile))

		assert.False(t, cfg.InsecureSkipVerify)
		assert.NotNil(t, cfg.RootCAs)
	})
}

func TestSetupTLSData(t *testing.T) {
	t.Parallel()

	t.Run("no data", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, SetupTLSData(new(tls.Config), nil), ErrInvalidCertData)
	})

	t.Run("invalid data", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, SetupTLSData(new(tls.Config), []byte("WTF?!")), ErrInvalidCertData)
	})

	t.Run("valid data", func(t *testing.T) {
		t.Parallel()

		cfg := &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // test
		}

		cert, _, err := testcerts.GenerateCerts()
		require.NoError(t, err)
		require.NoError(t, SetupTLSData(cfg, cert))

		assert.False(t, cfg.InsecureSkipVerify)
		assert.NotNil(t, cfg.RootCAs)
	})
}
