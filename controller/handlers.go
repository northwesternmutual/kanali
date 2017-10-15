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
	"fmt"

	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/spec"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

var apiProxyHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		defer logger.Sync()
		proxy, ok := obj.(*spec.APIProxy)
		if !ok {
			logger.Error("received malformed APIProxy from k8s apiserver")
		} else {
			spec.ProxyStore.Set(*proxy)
			logger.Debug(fmt.Sprintf("added ApiProxy %s in %s namespace", proxy.ObjectMeta.Name, proxy.ObjectMeta.Namespace))
		}
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		defer logger.Sync()
		oldProxy, ok := old.(*spec.APIProxy)
		if !ok {
			logger.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		newProxy, ok := new.(*spec.APIProxy)
		if !ok {
			logger.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		if err := spec.ProxyStore.Update(*oldProxy, *newProxy); err != nil {
			logger.Error(err.Error())
		} else {
			logger.Debug(fmt.Sprintf("updated ApiProxy %s in %s namespace", newProxy.ObjectMeta.Name, newProxy.ObjectMeta.Namespace))
		}
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		defer logger.Sync()
		proxy, ok := obj.(*spec.APIProxy)
		if !ok {
			logger.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		if result, _ := spec.ProxyStore.Delete(*proxy); result != nil {
			result := result.(spec.APIProxy)
			logger.Debug(fmt.Sprintf("deleted ApiProxy %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
		}
	},
}

var apiKeyHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		key, ok := obj.(*spec.APIKey)
		if !ok {
			logger.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*key).Decrypt(); err != nil {
			logger.Error(err.Error())
			return
		}
		spec.KeyStore.Set(*key)
		logger.Debug(fmt.Sprintf("added ApiKey %s", key.ObjectMeta.Name))
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		newKey, ok := new.(*spec.APIKey)
		if !ok {
			logger.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		oldKey, ok := old.(*spec.APIKey)
		if !ok {
			logger.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*newKey).Decrypt(); err != nil {
			logger.Error(err.Error())
			return
		}
		spec.KeyStore.Update(*oldKey, *newKey)
		logger.Debug(fmt.Sprintf("updated ApiKey %s", newKey.ObjectMeta.Name))
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		key, ok := obj.(*spec.APIKey)
		if !ok {
			logger.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*key).Decrypt(); err != nil {
			logger.Error(err.Error())
			return
		}
		result, _ := spec.KeyStore.Delete(*key)
		if result != nil {
			result := result.(spec.APIKey)
			logger.Debug(fmt.Sprintf("deleted ApiKey %s", result.ObjectMeta.Name))
		}
	},
}

var apiKeyBindingHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		binding, ok := obj.(*spec.APIKeyBinding)
		if !ok {
			logger.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		spec.BindingStore.Set(*binding)
		logger.Debug(fmt.Sprintf("added ApiKeyBinding %s in %s namespace", binding.ObjectMeta.Name, binding.ObjectMeta.Namespace))
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		newBinding, ok := new.(*spec.APIKeyBinding)
		if !ok {
			logger.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		oldBinding, ok := old.(*spec.APIKeyBinding)
		if !ok {
			logger.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		spec.BindingStore.Update(*newBinding, *oldBinding)
		logger.Debug(fmt.Sprintf("updated ApiKeyBinding %s in %s namespace", newBinding.ObjectMeta.Name, newBinding.ObjectMeta.Namespace))
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		binding, ok := obj.(*spec.APIKeyBinding)
		if !ok {
			logger.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		result, _ := spec.BindingStore.Delete(*binding)
		if result != nil {
			result := result.(spec.APIKeyBinding)
			logger.Debug(fmt.Sprintf("deleted ApiKeyBinding %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
		}
	},
}

var secretHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		secret, ok := obj.(*v1.Secret)
		if !ok {
			logger.Error("received malformed Secret from k8s apiserver")
			return
		}
		spec.SecretStore.Set(*secret)
		logger.Debug(fmt.Sprintf("added Secret %s in %s namespace", secret.ObjectMeta.Name, secret.ObjectMeta.Namespace))
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		oldSecret, ok := old.(*v1.Secret)
		if !ok {
			logger.Error("received malformed Secret from k8s apiserver")
			return
		}
		newSecret, ok := new.(*v1.Secret)
		if !ok {
			logger.Error("received malformed Secret from k8s apiserver")
			return
		}
		spec.SecretStore.Update(*oldSecret, *newSecret)
		logger.Debug(fmt.Sprintf("updated Secret %s in %s namespace", newSecret.ObjectMeta.Name, newSecret.ObjectMeta.Namespace))
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		secret, ok := obj.(*v1.Secret)
		if !ok {
			logger.Error("received malformed Secret from k8s apiserver")
			return
		}
		result, _ := spec.SecretStore.Delete(*secret)
		if result != nil {
			result := result.(v1.Secret)
			logger.Debug(fmt.Sprintf("deleted Secret %s in %s namespace", result.ObjectMeta.Name, result.ObjectMeta.Namespace))
		}
	},
}

var serviceHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		service, ok := obj.(*v1.Service)
		if !ok {
			logger.Error("received malformed Service from k8s apiserver")
			return
		}
		spec.ServiceStore.Set(spec.CreateService(*service))
		logger.Debug(fmt.Sprintf("added Service %s in %s namespace", service.ObjectMeta.Name, service.ObjectMeta.Namespace))
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		oldService, ok := old.(*v1.Service)
		if !ok {
			logger.Error("received malformed Service from k8s apiserver")
			return
		}
		newService, ok := new.(*v1.Service)
		if !ok {
			logger.Error("received malformed Service from k8s apiserver")
			return
		}
		spec.ServiceStore.Update(spec.CreateService(*oldService), spec.CreateService(*newService))
		logger.Debug(fmt.Sprintf("updated Service %s in %s namespace", newService.ObjectMeta.Name, newService.ObjectMeta.Namespace))
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		service, ok := obj.(*v1.Service)
		if !ok {
			logger.Error("received malformed Service from k8s apiserver")
			return
		}
		result, _ := spec.ServiceStore.Delete(spec.CreateService(*service))
		if result != nil {
			result := result.(spec.Service)
			logger.Debug(fmt.Sprintf("deleted Service %s in %s namespace", result.Name, result.Namespace))
		}
	},
}

var configMapHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		cm, ok := obj.(*v1.ConfigMap)
		if !ok {
			logger.Error("received malformed ConfigMap from k8s apiserver")
			return
		}
		spec.MockResponseStore.Set(*cm)
		logger.Debug(fmt.Sprintf("added ConfigMap %s in %s namespace", cm.ObjectMeta.Name, cm.ObjectMeta.Namespace))
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		oldConfigMap, ok := old.(*v1.ConfigMap)
		if !ok {
			logger.Error("received malformed ConfigMap from k8s apiserver")
			return
		}
		newConfigMap, ok := new.(*v1.ConfigMap)
		if !ok {
			logger.Error("received malformed ConfigMap from k8s apiserver")
			return
		}
		spec.MockResponseStore.Update(*oldConfigMap, *newConfigMap)
		logger.Debug(fmt.Sprintf("updated ConfigMap %s in %s namespace", newConfigMap.ObjectMeta.Name, newConfigMap.ObjectMeta.Namespace))
	},
	DeleteFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		cm, ok := obj.(*v1.ConfigMap)
		if !ok {
			logger.Error("received malformed ConfigMap from k8s apiserver")
			return
		}
		spec.MockResponseStore.Delete(*cm)
		logger.Debug(fmt.Sprintf("deleted ConfigMap %s in %s namespace", cm.ObjectMeta.Name, cm.ObjectMeta.Namespace))
	},
}

var endpointsHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		logger := logging.WithContext(nil)
		endpoints, ok := obj.(*v1.Endpoints)
		if !ok {
			logger.Error("received malformed Endpoints from k8s apiserver")
			return
		}
		if endpoints.ObjectMeta.Name == "kanali" {
			logger.Debug(fmt.Sprintf("adding Endpoints kanali in %s namespace", endpoints.ObjectMeta.Namespace))
			spec.KanaliEndpoints = endpoints
		}
	},
	UpdateFunc: func(old, new interface{}) {
		logger := logging.WithContext(nil)
		_, ok := old.(*v1.Endpoints)
		if !ok {
			logger.Error("received malformed Endpoints from k8s apiserver")
			return
		}
		newEndpoints, ok := new.(*v1.Endpoints)
		if !ok {
			logger.Error("received malformed Endpoints from k8s apiserver")
			return
		}
		if newEndpoints.ObjectMeta.Name == "kanali" {
			logger.Debug(fmt.Sprintf("updated Endpoints kanali in %s namespace", newEndpoints.ObjectMeta.Namespace))
			spec.KanaliEndpoints = newEndpoints
		}
	},
	DeleteFunc: func(obj interface{}) {},
}
