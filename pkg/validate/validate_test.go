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

package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/fake"
	"github.com/northwesternmutual/kanali/pkg/log"
)

func TestIsValidResource(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	validator := New(context.Background(), fake.NewSimpleClientset())

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiProxy",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group: "foo/bar",
	}, nil))

	assert.Nil(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group: "foo.bar.io",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiProxy",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiKey",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiKeyBinding",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "MockTarget",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiProxyList",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiKeyList",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "ApiKeyBindingList",
	}, nil))

	assert.Error(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "MockTargetList",
	}, nil))

	assert.Nil(t, validator.IsValidResource(metav1.GroupVersionKind{
		Group:   "kanali.io",
		Version: "v2",
		Kind:    "Foo",
	}, nil))
}
