package validate

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/log"
)

type Validation interface {
	IsValidApiProxy([]byte) error
	IsValidApiProxyList([]byte) error
	IsValidApiKeyBinding([]byte) error
	IsValidApiKeyBindingList([]byte) error
	IsValidApiKey([]byte) error
	IsValidApiKeyList([]byte) error
	IsValidMockTarget([]byte) error
	IsValidMockTargetList([]byte) error
}

type validation struct {
	clientset versioned.Interface
	ctx       context.Context
}

func New(ctx context.Context, i versioned.Interface) *validation {
	return &validation{
		clientset: i,
		ctx:       ctx,
	}
}

func (v *validation) IsValidResource(gvk metav1.GroupVersionKind, data []byte) error {
	logger := log.WithContext(v.ctx)

	gvRaw := fmt.Sprintf("%s/%s",
		gvk.Group,
		gvk.Version,
	)

	gv, err := schema.ParseGroupVersion(gvRaw)
	if err != nil {
		logger.Error(fmt.Sprintf("invalid group version %s", gvRaw))
		return err
	}

	if gv != v2.SchemeGroupVersion {
		logger.Info(fmt.Sprintf("will not perform validation on %s", gvk.String()))
		return nil
	}

	switch gvk.Kind {
	case "ApiProxy":
		return v.IsValidApiProxy(data)
	case "ApiProxyList":
		return v.IsValidApiProxyList(data)
	case "ApiKeyBinding":
		return v.IsValidApiKeyBinding(data)
	case "ApiKeyBindingList":
		return v.IsValidApiKeyBindingList(data)
	case "ApiKey":
		return v.IsValidApiKey(data)
	case "ApiKeyList":
		return v.IsValidApiKeyList(data)
	case "MockTarget":
		return v.IsValidMockTarget(data)
	case "MockTargetList":
		return v.IsValidMockTargetList(data)
	default:
		logger.Info(fmt.Sprintf("will not perform validation on %s", gvk.String()))
	}

	return nil
}
