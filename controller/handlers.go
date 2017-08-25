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
	case api.Endpoints:
		if endpoints, ok := obj.(api.Endpoints); ok {
			if endpoints.ObjectMeta.Name == "kanali" {
				spec.KanaliEndpoints = &endpoints
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
	}
}
