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
	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/spec"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

var apiProxyHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		proxy, ok := obj.(*spec.APIProxy)
		if !ok {
			logrus.Error("received malformed APIProxy from k8s apiserver")
			return
		}
		if err := spec.ProxyStore.Set(*proxy); err != nil {
			logrus.Errorf("could not add ApiProxy: %s", err.Error())
		}
	},
	UpdateFunc: func(old, new interface{}) {
		oldProxy, ok := old.(*spec.APIProxy)
		if !ok {
			logrus.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		newProxy, ok := new.(*spec.APIProxy)
		if !ok {
			logrus.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		if err := spec.ProxyStore.Update(*oldProxy, *newProxy); err != nil {
			logrus.Errorf("could not update ApiProxy: %s", err.Error())
		}
	},
	DeleteFunc: func(obj interface{}) {
		proxy, ok := obj.(*spec.APIProxy)
		if !ok {
			logrus.Error("received malformed ApiProxy from k8s apiserver")
			return
		}
		if _, err := spec.ProxyStore.Delete(*proxy); err != nil {
			logrus.Errorf("could not delete ApiProxy: %s", err.Error())
		}
	},
}

var apiKeyHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		key, ok := obj.(*spec.APIKey)
		if !ok {
			logrus.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*key).Decrypt(); err != nil {
			logrus.Errorf("error decrypting ApiKey %s", key.ObjectMeta.Name)
			return
		}
		if err := spec.KeyStore.Set(*key); err != nil {
			logrus.Errorf("could not add ApiKey: %s", err.Error())
		}
	},
	UpdateFunc: func(old, new interface{}) {
		key, ok := old.(*spec.APIKey)
		if !ok {
			logrus.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*key).Decrypt(); err != nil {
			logrus.Errorf("error decrypting ApiKey %s", key.ObjectMeta.Name)
			return
		}
		if err := spec.KeyStore.Set(*key); err != nil {
			logrus.Errorf("could not update ApiKey: %s", err.Error())
		}
	},
	DeleteFunc: func(obj interface{}) {
		key, ok := obj.(*spec.APIKey)
		if !ok {
			logrus.Error("received malformed ApiKey from k8s apiserver")
			return
		}
		if err := (*key).Decrypt(); err != nil {
			logrus.Errorf("error decrypting ApiKey %s", key.ObjectMeta.Name)
			return
		}
		if _, err := spec.KeyStore.Delete(*key); err != nil {
			logrus.Errorf("could not delete ApiKey: %s", err.Error())
		}
	},
}

var apiKeyBindingHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		binding, ok := obj.(*spec.APIKeyBinding)
		if !ok {
			logrus.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		if err := spec.BindingStore.Set(*binding); err != nil {
			logrus.Errorf("could not add ApiKeyBinding: %s", err.Error())
		}
	},
	UpdateFunc: func(old, new interface{}) {
		binding, ok := old.(*spec.APIKeyBinding)
		if !ok {
			logrus.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		if err := spec.BindingStore.Set(*binding); err != nil {
			logrus.Errorf("could not update ApiKeyBinding: %s", err.Error())
		}
	},
	DeleteFunc: func(obj interface{}) {
		binding, ok := obj.(*spec.APIKeyBinding)
		if !ok {
			logrus.Error("received malformed ApiKeyBinding from k8s apiserver")
			return
		}
		if _, err := spec.BindingStore.Delete(*binding); err != nil {
			logrus.Errorf("could not delete ApiKeyBinding: %s", err.Error())
		}
	},
}

var secretHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		secret, ok := obj.(*v1.Secret)
		if !ok {
			logrus.Error("received malformed Secret from k8s apiserver")
			return
		}
		if err := spec.SecretStore.Set(*secret); err != nil {
			logrus.Errorf("could not add Secret: %s", err.Error())
		}
	},
	UpdateFunc: func(old, new interface{}) {
		secret, ok := old.(*v1.Secret)
		if !ok {
			logrus.Error("received malformed Secret from k8s apiserver")
			return
		}
		if err := spec.SecretStore.Set(*secret); err != nil {
			logrus.Errorf("could not update Secret: %s", err.Error())
		}
	},
	DeleteFunc: func(obj interface{}) {
		secret, ok := obj.(*v1.Secret)
		if !ok {
			logrus.Error("received malformed Secret from k8s apiserver")
			return
		}
		if _, err := spec.SecretStore.Delete(*secret); err != nil {
			logrus.Errorf("could not delete Secret: %s", err.Error())
		}
	},
}

var serviceHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		service, ok := obj.(*v1.Service)
		if !ok {
			logrus.Error("received malformed Service from k8s apiserver")
			return
		}
		if err := spec.ServiceStore.Set(spec.CreateService(*service)); err != nil {
			logrus.Errorf("could not add Service: %s", err.Error())
		}
	},
	UpdateFunc: func(old, new interface{}) {
		service, ok := old.(*v1.Service)
		if !ok {
			logrus.Error("received malformed Service from k8s apiserver")
			return
		}
		if err := spec.ServiceStore.Set(spec.CreateService(*service)); err != nil {
			logrus.Errorf("could not update Service: %s", err.Error())
		}
	},
	DeleteFunc: func(obj interface{}) {
		service, ok := obj.(*v1.Service)
		if !ok {
			logrus.Error("received malformed Service from k8s apiserver")
			return
		}
		if _, err := spec.ServiceStore.Delete(spec.CreateService(*service)); err != nil {
			logrus.Errorf("could not delete Service: %s", err.Error())
		}
	},
}

var endpointsHandlerFuncs = cache.ResourceEventHandlerFuncs{
	AddFunc: func(obj interface{}) {
		endpoints, ok := obj.(*v1.Endpoints)
		if !ok {
			logrus.Error("received malformed Endpoints from k8s apiserver")
			return
		}
		if endpoints.ObjectMeta.Name == "kanali" {
			logrus.Debugf("adding Kanali endpoints object")
			spec.KanaliEndpoints = endpoints
		}
	},
	UpdateFunc: func(old, new interface{}) {
		endpoints, ok := old.(*v1.Endpoints)
		if !ok {
			logrus.Error("received malformed Endpoints from k8s apiserver")
			return
		}
		if endpoints.ObjectMeta.Name == "kanali" {
			logrus.Debugf("updating Kanali endpoints object")
			spec.KanaliEndpoints = endpoints
		}
	},
	DeleteFunc: func(obj interface{}) {
		return
	},
}
