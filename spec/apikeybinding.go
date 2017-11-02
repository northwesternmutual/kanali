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
	"errors"
	"regexp"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

// APIKeyBindingList represents a list of APIKeyBindingings
type APIKeyBindingList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	Bindings             []APIKeyBinding `json:"items"`
}

// APIKeyBinding represents the TPR for an APIKeyBinding
type APIKeyBinding struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 APIKeyBindingSpec `json:"spec"`
}

// APIKeyBindingSpec represents the data fields for the APIKeyBinding TPR
type APIKeyBindingSpec struct {
	APIProxyName string `json:"proxy"`
	Keys         []Key  `json:"keys"`
}

// Rate defines rate limit rule
type Rate struct {
	Amount int    `json:"amount,omitempty"`
	Unit   string `json:"unit,omitempty"`
}

// Path represents the fine grained subpath that
// finer permissions will be assined for this apikey
type Path struct {
	Path string `json:"path,omitempty"`
	Rule Rule   `json:"rule,omitempty"`
}

// Key defines an apikey that has some level of permissions
// the the proxy this binding is bound to
type Key struct {
	Name        string  `json:"name"`
	Quota       int     `json:"quota,omitempty"`
	Rate        *Rate   `json:"rate,omitempty"`
	DefaultRule Rule    `json:"defaultRule,omitempty"`
	Subpaths    []*Path `json:"subpaths,omitempty"`
}

// Rule defines the global and granular rules that this
// apikey should be assigned
type Rule struct {
	Global   bool           `json:"global,omitempty"`
	Granular *GranularProxy `json:"granular,omitempty"`
}

// GranularProxy defines the list of HTTP methods that this
// key has access to
type GranularProxy struct {
	Verbs []string `json:"verbs,omitempty"`
}

// BindingFactory is factory that implements a concurrency safe store for Kanali APIKeyBindings
type BindingFactory struct {
	mutex      sync.RWMutex
	bindingMap map[string]map[string]APIKeyBinding
}

// BindingStore holds all Kanali APIKeyBindings that Kanali has discovered
// in a cluster. It should not be mutated directly!
var BindingStore *BindingFactory

func init() {
	BindingStore = &BindingFactory{sync.RWMutex{}, map[string]map[string]APIKeyBinding{}}
}

// Clear will remove all bindings from the store
func (s *BindingFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for b := range s.bindingMap {
		delete(s.bindingMap, b)
	}
}

// IsEmpty reports whether the binding store is empty
func (s *BindingFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.bindingMap) == 0
}

// Update will update an APIKeyBinding
func (s *BindingFactory) Update(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	binding, ok := obj.(APIKeyBinding)
	if !ok {
		return errors.New("grrr - you're only allowed add api key bindings to the api key binding store.... duh")
	}
	return s.set(binding)
}

// Set takes a APIKeyBinding and either adds it to the store
// or updates it
func (s *BindingFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	binding, ok := obj.(APIKeyBinding)
	if !ok {
		return errors.New("grrr - you're only allowed add api key bindings to the api key binding store.... duh")
	}
	return s.set(binding)
}

func (s *BindingFactory) set(binding APIKeyBinding) error {
	logrus.Infof("Adding new APIKeyBinding named %s in namespace %s", binding.ObjectMeta.Name, binding.ObjectMeta.Namespace)

	for _, key := range binding.Spec.Keys {
		for _, subpath := range key.Subpaths {
			if subpath.Path[0] != '/' {
				subpath.Path = "/" + subpath.Path
			}
		}
	}

	if s.bindingMap[binding.ObjectMeta.Namespace] == nil {
		s.bindingMap[binding.ObjectMeta.Namespace] = map[string]APIKeyBinding{
			binding.Spec.APIProxyName: binding,
		}
		return nil
	}
	s.bindingMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName] = binding
	return nil
}

// Get retrieves a particual binding in the store. If not found, nil is returned.
func (s *BindingFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return nil, errors.New("should only pass the proxy name and namespace name")
	}
	proxyName, ok := params[0].(string)
	if !ok {
		return nil, errors.New("proxy name should be a string")
	}
	namespace, ok := params[1].(string)
	if !ok {
		return nil, errors.New("namespace should be a string")
	}
	if _, ok := s.bindingMap[namespace]; !ok {
		return nil, nil
	}
	if val, ok := s.bindingMap[namespace][proxyName]; ok {
		return val, nil
	}
	return nil, nil
}

// Delete will remove a particular binding from the store
func (s *BindingFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	binding, ok := obj.(APIKeyBinding)
	if !ok {
		return nil, errors.New("there's no way this api key binding could've gotten in here")
	}
	val, ok := s.bindingMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName]
	if !ok {
		return nil, nil
	}
	delete(s.bindingMap[binding.ObjectMeta.Namespace], binding.Spec.APIProxyName)
	if len(s.bindingMap[binding.ObjectMeta.Namespace]) == 0 {
		delete(s.bindingMap, binding.ObjectMeta.Namespace)
	}
	return val, nil
}

// GetRule returns the highest priority rule to use
// for the incoming request path
func (k *Key) GetRule(targetPath string) Rule {
	for _, subpath := range k.Subpaths {
		if result, err := regexp.MatchString("^"+subpath.Path, targetPath); err != nil || !result {
			continue
		}
		return subpath.Rule
	}
	return k.DefaultRule
}

// GetAPIKey retrieves a pointer to a Key object for a given
// apikey name
func (b *APIKeyBinding) GetAPIKey(apiKeyName string) *Key {

	for _, key := range b.Spec.Keys {
		if strings.ToLower(key.Name) == strings.ToLower(apiKeyName) {
			return &key
		}
	}
	return nil

}
