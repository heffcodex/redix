package redix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_Namespace(t *testing.T) {
	require.Equal(t, Namespace("foo"), (&Client{ns: "foo"}).Namespace())
}

func TestClient_KeyPrefix(t *testing.T) {
	assert.Equal(t, "", (&Client{}).KeyPrefix())
	assert.Equal(t, "foo:", (&Client{ns: "foo"}).KeyPrefix())
}

func TestClient_Key(t *testing.T) {
	assert.Equal(t, "", (&Client{ns: ""}).Key(""))
	assert.Equal(t, "bar", (&Client{ns: ""}).Key("bar"))
	assert.Equal(t, "foo:", (&Client{ns: "foo"}).Key(""))
	assert.Equal(t, "foo:bar", (&Client{ns: "foo"}).Key("bar"))
}
