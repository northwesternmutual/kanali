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
	"strings"
	"sync"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type ApiKeyBindingStoreInterface interface {
	Set(apiKeyBinding *v2.ApiKeyBinding)
	Update(old, new *v2.ApiKeyBinding)
	Get(namespace, binding, key, target string) (*v2.Rule, v2.Rate)
	Delete(apiKeyBinding *v2.ApiKeyBinding) error
	Clear()
	IsEmpty() bool
	ApiKeyBindingStoreExpansion
}

type apiKeyBindingFactory struct {
	mutex            sync.RWMutex
	apiKeyBindingMap map[string]map[string]map[string]structuredKey
}

type structuredKey struct {
	key         v2.Key
	subpathTree *subpathNode
}

type subpathNode struct {
	children map[string]*subpathNode
	value    *v2.Path
}

var (
	apiKeyBindingStore = &apiKeyBindingFactory{sync.RWMutex{}, map[string]map[string]map[string]structuredKey{}}
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
	if _, ok := s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace]; !ok {
		s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace] = map[string]map[string]structuredKey{}
	}

	// binding is the second level
	if _, ok := s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace][apiKeyBinding.ObjectMeta.Name]; !ok {
		s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace][apiKeyBinding.ObjectMeta.Name] = map[string]structuredKey{}
	}

	// keys are the third level
	for _, key := range apiKeyBinding.Spec.Keys {
		s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace][apiKeyBinding.ObjectMeta.Name][key.Name] = structuredKey{
			key:         key,
			subpathTree: generateSubpathTree(key),
		}
	}
}

// Get retrieves the highest priority rule given:
//   1. namespace name
//   2. binding name
//   3. api key name
//   4. target path
// O(n), n => number of path segments in target path
func (s *apiKeyBindingFactory) Get(namespace, binding, key, target string) (*v2.Rule, v2.Rate) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result, ok := s.apiKeyBindingMap[namespace][binding][key]
	if !ok {
		return nil, v2.Rate{}
	}
	return result.getHighestPriorityRule(target), result.key.Rate
}

// Delete will remove an ApiKeyBinding
// O(1)
func (s *apiKeyBindingFactory) Delete(apiKeyBinding *v2.ApiKeyBinding) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if apiKeyBinding == nil {
		return nil
	}
	_, ok := s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace][apiKeyBinding.ObjectMeta.Name]
	if !ok {
		return errors.New("ApiKeyBinding not found")
	}
	delete(s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace], apiKeyBinding.ObjectMeta.Name)
	if len(s.apiKeyBindingMap[apiKeyBinding.ObjectMeta.Namespace]) == 0 {
		delete(s.apiKeyBindingMap, apiKeyBinding.ObjectMeta.Namespace)
	}
	return nil
}

func generateSubpathTree(key v2.Key) *subpathNode {
	root := &subpathNode{}

	for _, subpath := range key.Subpaths {
		if subpath.Path[0] == '/' {
			root.doSetSubpath(strings.Split(subpath.Path[1:], "/"), subpath)
		} else {
			root.doSetSubpath(strings.Split(subpath.Path, "/"), subpath)
		}
	}

	return root
}

func (n *subpathNode) doSetSubpath(pathSegments []string, subpath v2.Path) {
	if n.children == nil {
		n.children = map[string]*subpathNode{}
	}
	if n.children[pathSegments[0]] == nil {
		n.children[pathSegments[0]] = &subpathNode{}
	}
	if len(pathSegments) < 2 {
		n.children[pathSegments[0]].value = &subpath
	} else {
		n.children[pathSegments[0]].doSetSubpath(pathSegments[1:], subpath)
	}
}

func (k structuredKey) getHighestPriorityRule(path string) *v2.Rule {
	subpath := k.subpathTree.getSubpath(path)
	if subpath == nil {
		return &k.key.DefaultRule
	}
	return &subpath.Rule
}

func (n *subpathNode) getSubpath(path string) *v2.Path {
	if len(n.children) == 0 || path == "" {
		return nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	for _, part := range strings.Split(path, "/") {
		if n.children[part] == nil {
			break
		} else {
			n = n.children[part]
		}
	}
	return n.value
}
