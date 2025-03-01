package redix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamespace_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "foo", Namespace("foo").String())
}

func TestNamespace_Append(t *testing.T) {
	t.Parallel()

	assert.Equal(t, Namespace(""), Namespace("").Append(""))
	assert.Equal(t, Namespace("foo"), Namespace("").Append("foo"))
	assert.Equal(t, Namespace("foo:"), Namespace("foo").Append(""))
	assert.Equal(t, Namespace("foo:bar"), Namespace("foo").Append("bar"))
}
