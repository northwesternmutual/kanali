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
	"k8s.io/kubernetes/pkg/api"
)

type handlerFuncs interface {
	addFunc(obj interface{})
	updateFunc(obj interface{})
	deleteFunc(obj interface{})
}

type k8sEventHandler struct{}

func (h k8sEventHandler) addFunc(obj interface{}) {
	switch obj.(type) {
	case spec.APIProxy:
		if proxy, ok := obj.(spec.APIProxy); ok {
			err := spec.ProxyStore.Set(proxy)
			if err != nil {
				logrus.Errorf("could not add/modify api proxy. skipping: %s", err.Error())
			}
		}
	case spec.APIKey:
		if key, ok := obj.(spec.APIKey); ok {
			if err := key.Decrypt(); err == nil {
				err := spec.KeyStore.Set(key)
				if err != nil {
					logrus.Errorf("could not add/modify apikey. skipping: %s", err.Error())
				}
			} else {
				logrus.Error("could not decrypt apikey. skipping...")
			}
		}
	case spec.APIKeyBinding:
		if binding, ok := obj.(spec.APIKeyBinding); ok {
			err := spec.BindingStore.Set(binding)
			if err != nil {
				logrus.Errorf("could not add/modify apikey binding. skipping: %s", err.Error())
			}
		}
	case api.Secret:
		if secret, ok := obj.(api.Secret); ok {
			err := spec.SecretStore.Set(secret)
			if err != nil {
				logrus.Errorf("could not add/modify secret. skipping: %s", err.Error())
			}
		}
	case api.Service:
		if service, ok := obj.(api.Service); ok {
			err := spec.ServiceStore.Set(spec.CreateService(service))
			if err != nil {
				logrus.Errorf("could not add/modify service. skipping: %s", err.Error())
			}
		}
	case api.ConfigMap:
		if cm, ok := obj.(api.ConfigMap); ok {
			if err := spec.MockResponseStore.Set(cm); err != nil {
				logrus.Errorf("could not add/modify configmap. skipping: %s", err.Error())
			}
		}
	}
}

func (h k8sEventHandler) updateFunc(obj interface{}) {
	h.addFunc(obj)
}

func (h k8sEventHandler) deleteFunc(obj interface{}) {
	switch obj.(type) {
	case spec.APIProxy:
		if proxy, ok := obj.(spec.APIProxy); ok {
			_, err := spec.ProxyStore.Delete(proxy)
			if err != nil {
				logrus.Errorf("could not delete api proxy. skipping: %s", err.Error())
			}
		}
	case spec.APIKey:
		if key, ok := obj.(spec.APIKey); ok {
			if err := key.Decrypt(); err == nil {
				_, err := spec.KeyStore.Delete(key)
				if err != nil {
					logrus.Errorf("could not delete apikey. skipping: %s", err.Error())
				}
			} else {
				logrus.Warn("could not decrypt apikey. skipping...")
			}
		}
	case spec.APIKeyBinding:
		if binding, ok := obj.(spec.APIKeyBinding); ok {
			_, err := spec.BindingStore.Delete(binding)
			if err != nil {
				logrus.Errorf("could not delete apikey binding. skipping: %s", err.Error())
			}
		}
	case api.Secret:
		if secret, ok := obj.(api.Secret); ok {
			_, err := spec.SecretStore.Delete(secret)
			if err != nil {
				logrus.Errorf("could not delete secret. skipping: %s", err.Error())
			}
		}
	case api.Service:
		if service, ok := obj.(api.Service); ok {
			_, err := spec.ServiceStore.Delete(spec.CreateService(service))
			if err != nil {
				logrus.Errorf("could not delete service. skipping: %s", err.Error())
			}
		}
	case api.ConfigMap:
		if cm, ok := obj.(api.ConfigMap); ok {
			if _, err := spec.MockResponseStore.Delete(cm); err != nil {
				logrus.Errorf("could not delete configmap. skipping: %s", err.Error())
			}
		}
	}
}
