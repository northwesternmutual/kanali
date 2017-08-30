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
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/utils"
	"k8s.io/kubernetes/pkg/api"
)

type mock []Route

type routeNode struct {
	Children map[string]*routeNode
	Value    *Route
}

// Route represents the details for a mock response route
type Route struct {
	Route  string      `json:"route"`
	Code   int         `json:"code"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// MockResponseFactory is factory that implements a concurrency safe store for Kubernetes config maps
type MockResponseFactory struct {
	mutex        sync.RWMutex
	mockRespTree map[string]map[string]map[string]*routeNode
}

// MockResponseStore holds all Kubernetes secrets that Kanali has discovered
// in a cluster. It should not be mutated directly!
var MockResponseStore *MockResponseFactory

func init() {
	MockResponseStore = &MockResponseFactory{sync.RWMutex{}, map[string]map[string]map[string]*routeNode{}}
}

// Clear will remove all configmaps from the store
func (s *MockResponseFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.mockRespTree {
		delete(s.mockRespTree, k)
	}
}

// Set takes an APIProxy and either adds it to the store
// or updates it
func (s *MockResponseFactory) Set(obj interface{}) error {
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

	logrus.Debugf("adding mock response %s", cm.ObjectMeta.Name)

	for _, route := range m {
		if !utils.IsValidHTTPMethod(route.Method) {
			logrus.Warnf("route %s in ConfigMap %s contains an invalid HTTP method", route.Route, cm.ObjectMeta.Name)
			continue
		}
		copyRoute := route
		if _, ok := s.mockRespTree[cm.ObjectMeta.Namespace]; !ok {
			s.mockRespTree[cm.ObjectMeta.Namespace] = map[string]map[string]*routeNode{
				cm.ObjectMeta.Name: {
					route.Method: {},
				},
			}
		}
		if _, ok := s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name]; !ok {
			s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name] = map[string]*routeNode{
				route.Method: {},
			}
		}
		if _, ok := s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name][route.Method]; !ok {
			s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name][route.Method] = &routeNode{}
		}
		if route.Route == "" || route.Route == "/" {
			s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name][route.Method].Value = &copyRoute
		} else if route.Route[0] == '/' {
			s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name][route.Method].doSetRoute(strings.Split(route.Route[1:], "/"), &copyRoute)
		} else {
			s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name][route.Method].doSetRoute(strings.Split(route.Route, "/"), &copyRoute)
		}
	}
	return nil
}

func (n *routeNode) doSetRoute(pathSegments []string, r *Route) {
	if n.Children == nil {
		n.Children = map[string]*routeNode{}
	}
	if n.Children[pathSegments[0]] == nil {
		n.Children[pathSegments[0]] = &routeNode{}
	}
	if len(pathSegments) < 2 {
		n.Children[pathSegments[0]].Value = r
	} else {
		n.Children[pathSegments[0]].doSetRoute(pathSegments[1:], r)
	}
}

// Get retrieves a particual mock response from the store. If not found, nil is returned.
func (s *MockResponseFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 4 {
		return nil, errors.New("four parameters expected")
	}
	namespace, ok := params[0].(string)
	if !ok {
		return nil, errors.New("expecting namespace")
	}
	name, ok := params[1].(string)
	if !ok {
		return nil, errors.New("expecting name")
	}
	path, ok := params[2].(string)
	if !ok {
		return nil, errors.New("expecting path")
	}
	method, ok := params[3].(string)
	if !ok {
		return nil, errors.New("expecting method")
	}

	tree, ok := s.mockRespTree[namespace][name][method]
	if !ok {
		return nil, nil
	}

	if path == "" || path == "/" {
		return tree.Value, nil
	}

	if len(tree.Children) == 0 {
		return nil, nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	rootNode := tree
	for _, part := range strings.Split(path, "/") {
		if rootNode.Children[part] == nil {
			break
		}
		rootNode = rootNode.Children[part]
	}
	if rootNode.Value == nil {
		return nil, nil
	}
	return *rootNode.Value, nil
}

// Delete will remove a particular secret from the store
func (s *MockResponseFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	cm, ok := obj.(api.ConfigMap)
	if !ok {
		return nil, errors.New("obj was not a ConfigMap")
	}
	if _, ok := s.mockRespTree[cm.ObjectMeta.Namespace][cm.ObjectMeta.Name]; !ok {
		return nil, nil
	}
	delete(s.mockRespTree[cm.ObjectMeta.Namespace], cm.ObjectMeta.Name)
	if len(s.mockRespTree[cm.ObjectMeta.Namespace]) == 0 {
		delete(s.mockRespTree, cm.ObjectMeta.Namespace)
	}
	return nil, nil
}

// IsEmpty reports whether the configmap store is empty
func (s *MockResponseFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.mockRespTree) == 0
}
