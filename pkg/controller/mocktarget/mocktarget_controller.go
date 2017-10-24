package mocktarget

import (
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/store"
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
