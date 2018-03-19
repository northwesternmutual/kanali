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

package v2

import (
	"errors"
	"regexp"
	"sync"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type ApiKeyBindingStoreInterface interface {
	Set(apiKeyBinding *v2.ApiKeyBinding)
	Update(old, new *v2.ApiKeyBinding)
	Contains(namespace, binding string) bool
	GetHightestPriorityRule(namespace, binding, key, target string) *v2.Rule
	Delete(apiKeyBinding *v2.ApiKeyBinding) error
	Clear()
	IsEmpty() bool
}

type apiKeyBindingFactory struct {
	mutex            sync.RWMutex
	apiKeyBindingMap map[string]map[string]map[string]v2.Key
}

var (
	apiKeyBindingStore = &apiKeyBindingFactory{sync.RWMutex{}, map[string]map[string]map[string]v2.Key{}}
)

func ApiKeyBindingStore() ApiKeyBindingStoreInterface {
	return apiKeyBindingStore
}

// Clear will remove all ApiKeyBinding resources
// O(n), n => number of ApiKeyBinding resources
func (s *apiKeyBindingFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for b := range s.apiKeyBindingMap {
		delete(s.apiKeyBindingMap, b)
	}
}

// IsEmpty reports whether the ApiKeyBindingStore is empty
// O(1)
func (s *apiKeyBindingFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.apiKeyBindingMap) == 0
}

// Update will update an ApiKeyBinding resource
// O(x * y * z),
//   x => number of ApiKey resources ApiKeyBinding
//   y => number of subpaths defined in each ApiKey rule
//   z => number of path segments in each subpath
func (s *apiKeyBindingFactory) Update(old, new *v2.ApiKeyBinding) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set(new)
}

// Set will add an ApiKeyBinding resource
// O(x * y * z),
//   x => number of ApiKey resources ApiKeyBinding
//   y => number of subpaths defined in each ApiKey rule
//   z => number of path segments in each subpath
func (s *apiKeyBindingFactory) Set(apiKeyBinding *v2.ApiKeyBinding) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set(apiKeyBinding)
}

func (s *apiKeyBindingFactory) set(apiKeyBinding *v2.ApiKeyBinding) {
	// namespace is the first level
	if _, ok := s.apiKeyBindingMap[apiKeyBinding.GetNamespace()]; !ok {
		s.apiKeyBindingMap[apiKeyBinding.GetNamespace()] = map[string]map[string]v2.Key{}
	}

	// binding is the second level
	if _, ok := s.apiKeyBindingMap[apiKeyBinding.GetNamespace()][apiKeyBinding.GetName()]; !ok {
		s.apiKeyBindingMap[apiKeyBinding.GetNamespace()][apiKeyBinding.GetName()] = map[string]v2.Key{}
	}

	keys := make(map[string]v2.Key, len(apiKeyBinding.Spec.Keys))

	// keys are the third level
	for _, key := range apiKeyBinding.Spec.Keys {
		keys[key.Name] = key
	}

	s.apiKeyBindingMap[apiKeyBinding.GetNamespace()][apiKeyBinding.GetName()] = keys
}

func (s *apiKeyBindingFactory) Contains(namespace, binding string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.contains(namespace, binding)
}

func (s *apiKeyBindingFactory) contains(namespace, binding string) bool {
	_, ok := s.apiKeyBindingMap[namespace][binding]
	return ok
}

// GetHightestPriorityRule retrieves the highest priority rule given:
//   1. namespace name
//   2. binding name
//   3. api key name
//   4. target path
// O(n), n => number of path segments in target path
func (s *apiKeyBindingFactory) GetHightestPriorityRule(namespace, binding, key, target string) *v2.Rule {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result, ok := s.apiKeyBindingMap[namespace][binding][key]
	if !ok {
		return nil
	}
	for _, subpath := range result.Subpaths {
		if result, err := regexp.MatchString("^"+subpath.Path, target); err != nil || !result {
			continue
		}
		return &subpath.Rule
	}
	return &result.DefaultRule
}

// Delete will remove an ApiKeyBinding
// O(1)
func (s *apiKeyBindingFactory) Delete(apiKeyBinding *v2.ApiKeyBinding) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if apiKeyBinding == nil {
		return nil
	}
	_, ok := s.apiKeyBindingMap[apiKeyBinding.GetNamespace()][apiKeyBinding.GetName()]
	if !ok {
		return errors.New("ApiKeyBinding not found")
	}
	delete(s.apiKeyBindingMap[apiKeyBinding.GetNamespace()], apiKeyBinding.GetName())
	if len(s.apiKeyBindingMap[apiKeyBinding.GetNamespace()]) == 0 {
		delete(s.apiKeyBindingMap, apiKeyBinding.GetNamespace())
	}
	return nil
}
