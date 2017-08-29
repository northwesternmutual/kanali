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

package spec

import (
	"crypto/tls"
	"errors"
	"sync"

	"github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
)

type mock []route

type routeNode {
  Children map[string]*routeNode
	Value    *route
}

type route struct {
	Route  string      `json:"route"`
	Code   int         `json:"code"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// ConfigMapFactory is factory that implements a concurrency safe store for Kubernetes config maps
type ConfigMapFactory struct {
	mutex     sync.RWMutex
	routeTree *routeNode
}

// SecretStore holds all Kubernetes secrets that Kanali has discovered
// in a cluster. It should not be mutated directly!
var ConfigMapStore *ConfigMapFactory

func init() {
	ConfigMapStore = &ConfigMapFactory{sync.RWMutex{}, map[string]map[string]api.ConfigMap{}}
}

// Clear will remove all configmaps from the store
func (s *ConfigMapFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.configmapMap {
		delete(s.configmapMap, k)
	}
}

// Set takes a ConfigMap and either adds it to the store
// or updates it
func (s *ConfigMapFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	cm, ok := obj.(api.ConfigMap)
	if !ok {
		return errors.New("obj was not a ConfigMap")
	}

  mockResponse, ok := cm.Data["response"]
	if !ok {
		logrus.Debugf("ConfigMap %s does not contains a response data field", cm.ObjectMeta.Name)
    return nil
	}

  var m mock
  if err := json.Unmarshal([]byte(mockResponse), &m); err != nil {
    logrus.Debugf("ConfigMap %s does not contains a properly formed response field", cm.ObjectMeta.Name)
    return nil
  }

	logrus.Debugf("adding ConfigMap %s", cm.ObjectMeta.Name)
	if _, ok := s.configmapMap[cm.ObjectMeta.Namespace]; ok {
		s.configmapMap[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name] = m
	} else {
		s.configmapMap[cm.ObjectMeta.Namespace] = map[string]api.ConfigMap{
			cm.ObjectMeta.Name: m,
		}
	}
	return nil
}

// Get retrieves a particual configmap in the store. If not found, nil is returned.
func (s *ConfigMapFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return nil, errors.New("should should take 2 params, name and namespace")
	}
	name, ok := params[0].(string)
	if !ok {
		return nil, errors.New("configmap name must be of type string")
	}
	namespace, ok := params[1].(string)
	if !ok {
		return nil, errors.New("configmap namespace must be of type string")
	}
	cm, ok := s.configmapMap[namespace][name]
	if !ok {
		return nil, nil
	}
	return cm, nil
}

// Delete will remove a particular secret from the store
func (s *ConfigMapFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	cm, ok := obj.(api.ConfigMap)
	if !ok {
		return nil, errors.New("obj was not a ConfigMap")
	}
	if _, ok = s.configmapMap[cm.ObjectMeta.Namespace]; !ok {
		return nil, nil
	}
	oldCm, ok := s.configmapMap[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name]
	if !ok {
		return nil, nil
	}
	delete(s.configmapMap[cm.ObjectMeta.Namespace], cm.ObjectMeta.Name)
	if len(s.configmapMap[cm.ObjectMeta.Namespace]) == 0 {
		delete(s.configmapMap, cm.ObjectMeta.Namespace)
	}
	return oldCm, nil
}

// Contains reports whether the secrets store contains a particular secret
func (s *ConfigMapFactory) Contains(params ...interface{}) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return false, errors.New("containers requires 2 params")
	}
	name, ok := params[0].(string)
	if !ok {
		return false, errors.New("first parameter should be a string")
	}
	namespace, ok := params[1].(string)
	if !ok {
		return false, errors.New("second parameter should be a string")
	}
	if _, ok := s.configmapMap[namespace]; !ok {
		return false, nil
	}
	_, ok = s.configmapMap[namespace][name]
	return ok, nil
}

// IsEmpty reports whether the configmap store is empty
func (s *ConfigMapFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.configmapMap) == 0
}
