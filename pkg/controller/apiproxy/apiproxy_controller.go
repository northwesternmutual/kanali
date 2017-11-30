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

package apiproxy

import (
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	tags "github.com/northwesternmutual/kanali/pkg/tags"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

type ApiProxyController struct {
	apiproxies informers.ApiProxyInformer
}

func NewApiProxyController(apiproxies informers.ApiProxyInformer) *ApiProxyController {

	ctlr := &ApiProxyController{}

	ctlr.apiproxies = apiproxies
	apiproxies.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    ctlr.apiProxyAdd,
			UpdateFunc: ctlr.apiProxyUpdate,
			DeleteFunc: ctlr.apiProxyDelete,
		},
		5*time.Minute,
	)

	return ctlr
}

func (ctlr *ApiProxyController) Run(stopCh <-chan struct{}) {
	ctlr.apiproxies.Informer().Run(stopCh)
}

func (ctlr *ApiProxyController) apiProxyAdd(obj interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	proxy, ok := obj.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
	} else {
		store.ApiProxyStore().Set(proxy)
		logger.With(
			zap.String(tags.KanaliProxyName, proxy.ObjectMeta.Name),
			zap.String(tags.KanaliProxyNamespace, proxy.ObjectMeta.Namespace),
		).Debug("added ApiProxy")
	}
}

func (ctlr *ApiProxyController) apiProxyUpdate(old interface{}, new interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	oldProxy, ok := old.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	newProxy, ok := new.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	store.ApiProxyStore().Update(oldProxy, newProxy)
	logger.With(
		zap.String(tags.KanaliProxyName, newProxy.ObjectMeta.Name),
		zap.String(tags.KanaliProxyNamespace, newProxy.ObjectMeta.Namespace),
	).Debug("updated ApiProxy")
}

func (ctlr *ApiProxyController) apiProxyDelete(obj interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	proxy, ok := obj.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	if result := store.ApiProxyStore().Delete(proxy); result != nil {
		logger.With(
			zap.String(tags.KanaliProxyName, proxy.ObjectMeta.Name),
			zap.String(tags.KanaliProxyNamespace, proxy.ObjectMeta.Namespace),
		).Debug("deleted ApiProxy")
	}
}
