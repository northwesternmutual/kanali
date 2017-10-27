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
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
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
		logger.Error("received malformed APIProxy from k8s apiserver")
	} else {
		store.ApiProxyStore().Set(proxy)
		logger.Debug(fmt.Sprintf("added ApiProxy %s in %s namespace", proxy.ObjectMeta.Name, proxy.ObjectMeta.Namespace))
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
	if err := store.ApiProxyStore().Update(oldProxy, newProxy); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Debug(fmt.Sprintf("updated ApiProxy %s in %s namespace", newProxy.ObjectMeta.Name, newProxy.ObjectMeta.Namespace))
	}
}

func (ctlr *ApiProxyController) apiProxyDelete(obj interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	proxy, ok := obj.(*v2.ApiProxy)
	if !ok {
		logger.Error("received malformed ApiProxy from k8s apiserver")
		return
	}
	if result, _ := store.ApiProxyStore().Delete(proxy); result != nil {
		logger.Debug(fmt.Sprintf("deleted ApiProxy %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
	}
}
