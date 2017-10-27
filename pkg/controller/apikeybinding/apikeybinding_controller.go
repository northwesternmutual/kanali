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

package apikeybinding

import (
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"k8s.io/client-go/tools/cache"
)

type ApiKeyBindingController struct {
	apikeybindings informers.ApiKeyBindingInformer
}

func NewApiKeyBindingController(apikeybindings informers.ApiKeyBindingInformer) *ApiKeyBindingController {

	ctlr := &ApiKeyBindingController{}

	ctlr.apikeybindings = apikeybindings
	apikeybindings.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    ctlr.apiKeyBindingAdd,
			UpdateFunc: ctlr.apiKeyBindingUpdate,
			DeleteFunc: ctlr.apiKeyBindingDelete,
		},
		5*time.Minute,
	)

	return ctlr
}

func (ctlr *ApiKeyBindingController) Run(stopCh <-chan struct{}) {
	ctlr.apikeybindings.Informer().Run(stopCh)
}

func (ctlr *ApiKeyBindingController) apiKeyBindingAdd(obj interface{}) {
	logger := logging.WithContext(nil)
	binding, ok := obj.(*v2.ApiKeyBinding)
	if !ok {
		logger.Error("received malformed ApiKeyBinding from k8s apiserver")
		return
	}
	store.ApiKeyBindingStore.Set(*binding)
	logger.Debug(fmt.Sprintf("added ApiKeyBinding %s in %s namespace", binding.ObjectMeta.Name, binding.ObjectMeta.Namespace))
}

func (ctlr *ApiKeyBindingController) apiKeyBindingUpdate(old interface{}, new interface{}) {
	logger := logging.WithContext(nil)
	newBinding, ok := new.(*v2.ApiKeyBinding)
	if !ok {
		logger.Error("received malformed ApiKeyBinding from k8s apiserver")
		return
	}
	oldBinding, ok := old.(*v2.ApiKeyBinding)
	if !ok {
		logger.Error("received malformed ApiKeyBinding from k8s apiserver")
		return
	}
	store.ApiKeyBindingStore.Update(*newBinding, *oldBinding)
	logger.Debug(fmt.Sprintf("updated ApiKeyBinding %s in %s namespace", newBinding.ObjectMeta.Name, newBinding.ObjectMeta.Namespace))
}

func (ctlr *ApiKeyBindingController) apiKeyBindingDelete(obj interface{}) {
	logger := logging.WithContext(nil)
	binding, ok := obj.(*v2.ApiKeyBinding)
	if !ok {
		logger.Error("received malformed ApiKeyBinding from k8s apiserver")
		return
	}
	result, _ := store.ApiKeyBindingStore.Delete(*binding)
	if result != nil {
		result := result.(v2.ApiKeyBinding)
		logger.Debug(fmt.Sprintf("deleted ApiKeyBinding %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
	}
}
