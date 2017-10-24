package apiproxy

import (
	"fmt"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	informers "github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/kanali/v2"
	"github.com/northwesternmutual/kanali/pkg/logging"
	"github.com/northwesternmutual/kanali/pkg/store"
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
		store.ApiProxyStore.Set(*proxy)
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
	if err := store.ApiProxyStore.Update(*oldProxy, *newProxy); err != nil {
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
	if result, _ := store.ApiProxyStore.Delete(*proxy); result != nil {
		result := result.(v2.ApiProxy)
		logger.Debug(fmt.Sprintf("deleted ApiProxy %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
	}
}
