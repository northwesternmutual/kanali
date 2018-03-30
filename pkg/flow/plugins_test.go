// Copyright (c) 2018 Northwestern Mutual.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
