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
