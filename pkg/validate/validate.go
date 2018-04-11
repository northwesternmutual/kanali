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
