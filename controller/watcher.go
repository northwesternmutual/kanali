// Copyright (c) 2017 Northwestern Mutual.
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

package controller

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/spec"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) Watch(ctx context.Context) error {
	logrus.Debug("starting watch on k8s resources")

	doWatchResource(ctx, c.RESTClient, "apiproxies", &spec.APIProxy{}, fields.Everything(), apiProxyHandlerFuncs)
	doWatchResource(ctx, c.RESTClient, "apikeys", &spec.APIKey{}, fields.Everything(), apiKeyHandlerFuncs)
	doWatchResource(ctx, c.RESTClient, "apikeybindings", &spec.APIKeyBinding{}, fields.Everything(), apiKeyBindingHandlerFuncs)
	doWatchResource(ctx, c.ClientSet.Core().RESTClient(), "secrets", &v1.Secret{}, fields.OneTermEqualSelector("type", "kubernetes.io/tls"), secretHandlerFuncs)
	doWatchResource(ctx, c.ClientSet.Core().RESTClient(), "services", &v1.Service{}, fields.Everything(), serviceHandlerFuncs)
	doWatchResource(ctx, c.ClientSet.Core().RESTClient(), "endpoints", &v1.Endpoints{}, fields.Everything(), endpointsHandlerFuncs)

	<-ctx.Done()
	return ctx.Err()
}

func doWatchResource(ctx context.Context, restClient rest.Interface, resourcePath string, obj runtime.Object, fieldSelector fields.Selector, handlerFuncs cache.ResourceEventHandlerFuncs) {
	logrus.Debugf("attempting to watch %s", resourcePath)

	source := cache.NewListWatchFromClient(
		restClient,
		resourcePath,
		v1.NamespaceAll,
		fieldSelector,
	)
	_, ctlr := cache.NewInformer(
		source,
		obj,
		5*time.Minute,
		handlerFuncs,
	)

	go ctlr.Run(ctx.Done())
}
