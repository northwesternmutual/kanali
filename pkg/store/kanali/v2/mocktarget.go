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

type MockTargetStoreInterface interface {
	Set(mockTarget *v2.MockTarget) error
	Update(old, new *v2.MockTarget) error
	Get(namespace, name, path, method string) *v2.Route
	Delete(mockTarget *v2.MockTarget) bool
	Clear()
	IsEmpty() bool
}

type mockTargetFactory struct {
	mutex        sync.RWMutex
	mockRespTree map[string]map[string]*routeNode
}

type routeNode struct {
	children map[string]*routeNode
	value    *v2.Route
}

var (
	mockTargetStore = &mockTargetFactory{sync.RWMutex{}, map[string]map[string]*routeNode{}}
)

func MockTargetStore() MockTargetStoreInterface {
	return mockTargetStore
}

// Clear will remove all MockTarget resources
// O(n), n => number of namespaces respresented by the current MockTarget resources
func (s *mockTargetFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.mockRespTree {
		delete(s.mockRespTree, k)
	}
}

// Set updates a MockTarget resource
// O(n * m),
//   n => number of Route resources in the MockTarget resource
//   m => number of path segments in the Route path
func (s *mockTargetFactory) Update(old, new *v2.MockTarget) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set(new)
}

// Set adds a MockTarget resource
// O(n * m),
//   n => number of Route resources in the MockTarget resource
//   m => number of path segments in the Route path
func (s *mockTargetFactory) Set(mt *v2.MockTarget) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.set(mt)
}

func (s *mockTargetFactory) set(mt *v2.MockTarget) error {
	if len(mt.Spec.Routes) < 1 {
		return errors.New("MockTarget must contain at least one route")
	}

	if _, ok := s.mockRespTree[mt.GetNamespace()]; !ok {
		s.mockRespTree[mt.GetNamespace()] = map[string]*routeNode{}
	}
	if _, ok := s.mockRespTree[mt.GetNamespace()][mt.GetName()]; !ok {
		s.mockRespTree[mt.GetNamespace()][mt.GetName()] = &routeNode{}
	}
	s.mockRespTree[mt.GetNamespace()][mt.GetName()] = generateRouteTree(mt)
	return nil
}

func generateRouteTree(mt *v2.MockTarget) *routeNode {
	root := &routeNode{}

	for _, route := range mt.Spec.Routes {
		if route.Path[0] == '/' {
			root.doSetRoute(strings.Split(route.Path[1:], "/"), route)
		} else {
			root.doSetRoute(strings.Split(route.Path, "/"), route)
		}
	}

	return root
}

func (n *routeNode) doSetRoute(pathSegments []string, route v2.Route) {
	if n.children == nil {
		n.children = map[string]*routeNode{}
	}
	if n.children[pathSegments[0]] == nil {
		n.children[pathSegments[0]] = &routeNode{}
	}
	if len(pathSegments) < 2 {
		n.children[pathSegments[0]].value = &route
	} else {
		n.children[pathSegments[0]].doSetRoute(pathSegments[1:], route)
	}
}

// Get retrieves the matching Route resource given:
//   1. namespace name
//   2. MockTarget name
//   3. target path
//   4. http method
// O(n), n => number of path segments in Route path
func (s *mockTargetFactory) Get(namespace, name, path, method string) *v2.Route {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.get(namespace, name, path, method)
}

func (s *mockTargetFactory) get(namespace, name, path, method string) *v2.Route {
	root, ok := s.mockRespTree[namespace][name]
	if !ok {
		return nil
	}

	if path == "" || path == "/" {
		return root.value
	}

	if len(root.children) == 0 {
		return nil
	}

	route := root.doGetRoute(path)
	if route == nil {
		return nil
	}

	for _, m := range route.Methods {
		if strings.EqualFold(m, method) {
			return route
		}
	}
	return nil
}

func (n *routeNode) doGetRoute(path string) *v2.Route {
	if len(n.children) == 0 {
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

// Delete will remove a MockTarget resource
// O(1)
func (s *mockTargetFactory) Delete(mt *v2.MockTarget) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.mockRespTree[mt.GetNamespace()][mt.GetName()]
	if ok {
		delete(s.mockRespTree[mt.GetNamespace()], mt.GetName())
		if len(s.mockRespTree[mt.GetNamespace()]) == 0 {
			delete(s.mockRespTree, mt.GetNamespace())
		}
	}
	return ok
}

// IsEmpty reports whether the MockTargetStore is empty
// O(1)
func (s *mockTargetFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.mockRespTree) == 0
}
