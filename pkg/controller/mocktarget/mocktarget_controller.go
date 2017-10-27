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

package mocktarget

import (
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"k8s.io/client-go/tools/cache"
)

type MockTargetController struct {
	mocktargets informers.MockTargetInformer
}

func NewMockTargetController(mocktargets informers.MockTargetInformer) *MockTargetController {

	ctlr := &MockTargetController{}

	ctlr.mocktargets = mocktargets
	mocktargets.Informer().AddEventHandlerWithResyncPeriod(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    ctlr.mockTargetAdd,
			UpdateFunc: ctlr.mockTargetUpdate,
			DeleteFunc: ctlr.mockTargetDelete,
		},
		5*time.Minute,
	)

	return ctlr
}

func (ctlr *MockTargetController) Run(stopCh <-chan struct{}) {
	ctlr.mocktargets.Informer().Run(stopCh)
}

func (ctlr *MockTargetController) mockTargetAdd(obj interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	target, ok := obj.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
	} else {
		store.MockTargetStore.Set(*target)
		logger.Debug(fmt.Sprintf("added MockTarget %s in %s namespace", target.ObjectMeta.Name, target.ObjectMeta.Namespace))
	}
}

func (ctlr *MockTargetController) mockTargetUpdate(old interface{}, new interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	oldTarget, ok := old.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
		return
	}
	newTarget, ok := new.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
		return
	}
	if err := store.MockTargetStore.Update(*oldTarget, *newTarget); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Debug(fmt.Sprintf("updated MockTarget %s in %s namespace", newTarget.ObjectMeta.Name, newTarget.ObjectMeta.Namespace))
	}
}

func (ctlr *MockTargetController) mockTargetDelete(obj interface{}) {
	logger := logging.WithContext(nil)
	defer logger.Sync()
	target, ok := obj.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
		return
	}
	if result, _ := store.MockTargetStore.Delete(*target); result != nil {
		result := result.(*v2.MockTarget)
		logger.Debug(fmt.Sprintf("deleted MockTarget %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
	}
}
