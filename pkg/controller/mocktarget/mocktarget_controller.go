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

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	store "github.com/northwesternmutual/kanali/pkg/store/kanali/v2"
	"k8s.io/client-go/tools/cache"
)

type mockTargetController struct{}

func NewController() cache.ResourceEventHandler {
	return &mockTargetController{}
}

func (ctlr *mockTargetController) OnAdd(obj interface{}) {
	logger := log.WithContext(nil)
	target, ok := obj.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
	} else {
		if err := store.MockTargetStore().Set(target); err != nil {
			logger.Warn(fmt.Sprintf("failed to add MockTarget %s in %s namespace: %s", target.GetName(), target.GetNamespace(), err))
		} else {
			logger.Debug(fmt.Sprintf("added MockTarget %s in %s namespace", target.GetName(), target.GetNamespace()))
		}
	}
}

func (ctlr *mockTargetController) OnUpdate(old interface{}, new interface{}) {
	logger := log.WithContext(nil)
	oldTarget, okOld := old.(*v2.MockTarget)
	newTarget, okNew := new.(*v2.MockTarget)
	if !okOld || !okNew {
		logger.Error("received malformed MockTarget from k8s apiserver")
		return
	}
	if err := store.MockTargetStore().Update(oldTarget, newTarget); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Debug(fmt.Sprintf("updated MockTarget %s in %s namespace", newTarget.GetName(), newTarget.GetNamespace()))
	}
}

func (ctlr *mockTargetController) OnDelete(obj interface{}) {
	logger := log.WithContext(nil)
	target, ok := obj.(*v2.MockTarget)
	if !ok {
		logger.Error("received malformed MockTarget from k8s apiserver")
		return
	}
	if result := store.MockTargetStore().Delete(target); result {
		logger.Debug(fmt.Sprintf("deleted MockTarget %s in %s namespace", target.GetName(), target.GetNamespace()))
	}
}
