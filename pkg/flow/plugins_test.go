package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

func TestCombinePath(t *testing.T) {
	plugin := v2.Plugin{
		Name:    "foo",
		Version: "",
	}

	assert.Equal(t, "/foo.so", combinePath("/", plugin))
	plugin.Version = "v1.2.3"
	assert.Equal(t, "/foo_v1.2.3.so", combinePath("/", plugin))
	assert.Equal(t, "foo_v1.2.3.so", combinePath("", plugin))
}

func TestGetPlugin(t *testing.T) {
	plugin, err := getPlugin(context.Background(), v2.Plugin{
		Name: "aPiKeY",
	})
	assert.Nil(t, err)
	assert.Equal(t, internalPlugins["apikey"], plugin)
}
