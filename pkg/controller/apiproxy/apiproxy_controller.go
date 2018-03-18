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
	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/tags"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
)

type apiProxyController struct{}

func NewController() cache.ResourceEventHandler {
	return &apiProxyController{}
}

func (ctlr *apiProxyController) OnAdd(obj interface{}) {
	logger := log.WithContext(nil)
	proxy, ok := obj.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
	} else {
		store.ApiProxyStore().Set(proxy)
		logger.With(
			zap.String(tags.KanaliProxyName, proxy.GetName()),
			zap.String(tags.KanaliProxyNamespace, proxy.GetNamespace()),
		).Debug("added ApiProxy")
	}
}

func (ctlr *apiProxyController) OnUpdate(old interface{}, new interface{}) {
	logger := log.WithContext(nil)
	oldProxy, okOld := old.(*v2.ApiProxy)
	newProxy, okNew := new.(*v2.ApiProxy)
	if !okOld || !okNew {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	store.ApiProxyStore().Update(oldProxy, newProxy)
	logger.With(
		zap.String(tags.KanaliProxyName, newProxy.GetName()),
		zap.String(tags.KanaliProxyNamespace, newProxy.GetNamespace()),
	).Debug("updated ApiProxy")
}

func (ctlr *apiProxyController) OnDelete(obj interface{}) {
	logger := log.WithContext(nil)
	proxy, ok := obj.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	if result := store.ApiProxyStore().Delete(proxy); result != nil {
		logger.With(
			zap.String(tags.KanaliProxyName, proxy.GetName()),
			zap.String(tags.KanaliProxyNamespace, proxy.GetNamespace()),
		).Debug("deleted ApiProxy")
	}
}
