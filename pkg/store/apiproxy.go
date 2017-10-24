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

package store

import (
	"errors"
	"strings"
	"sync"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

// ApiProxyFactory is factory that implements a concurrency safe store for Kanali ApiProxy resources
type ApiProxyFactory struct {
	mutex     sync.RWMutex
	proxyTree *proxyNode
}

type proxyNode struct {
	children map[string]*proxyNode
	value    *v2.ApiProxy
}

var (
	// ApiProxyStore holds all discovered ApiProxy resources in an efficient data structure.
	// It should not be mutated directly!
	ApiProxyStore = &ApiProxyFactory{sync.RWMutex{}, &proxyNode{}}
)

// Clear will remove all ApiProxy resources
// O(1)
func (s *ApiProxyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*(s.proxyTree) = proxyNode{}
}

// Update will update an ApiProxy resource.
// O(n), n => number of path segments in the ApiProxy source path
func (s *ApiProxyFactory) Update(old, new interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	oldProxy, ok := old.(v2.ApiProxy)
	if !ok {
		return errors.New("old ApiProxy not expected type")
	}
	newProxy, ok := new.(v2.ApiProxy)
	if !ok {
		return errors.New("new ApiProxy not expected type")
	}
	normalizeProxyPaths(&oldProxy)
	normalizeProxyPaths(&newProxy)
	return s.update(oldProxy, newProxy)
}

func (s *ApiProxyFactory) update(old, new v2.ApiProxy) error {
	untyped := s.get(new.Spec.Source.Path)
	if untyped != nil {
		typed, ok := untyped.(v2.ApiProxy)
		if !ok {
			return errors.New("expected type ApiProxy")
		}
		if new.ObjectMeta.Name != typed.ObjectMeta.Name || new.ObjectMeta.Namespace != typed.ObjectMeta.Namespace {
			return errors.New("ApiProxy with different ObjectMeta exists at this path")
		}
	}

	s.proxyTree.doSet(strings.Split(new.Spec.Source.Path[1:], "/"), &new)
	if old.Spec.Source.Path != new.Spec.Source.Path {
		s.proxyTree.delete(strings.Split(old.Spec.Source.Path[1:], "/"))
	}
	return nil
}

// Set adds an ApiProxy resource
// O(n), n => number of path segments in the ApiProxy source path
func (s *ApiProxyFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	proxy, ok := obj.(v2.ApiProxy)
	if !ok {
		return errors.New("ApiProxy not expected type")
	}
	normalizeProxyPaths(&proxy)
	s.proxyTree.doSet(strings.Split(proxy.Spec.Source.Path[1:], "/"), &proxy)
	return nil
}

func (n *proxyNode) doSet(keys []string, v *v2.ApiProxy) {
	if n.children == nil {
		n.children = map[string]*proxyNode{}
	}
	if n.children[keys[0]] == nil {
		n.children[keys[0]] = &proxyNode{}
	}
	if len(keys) < 2 {
		n.children[keys[0]].value = v
	} else {
		n.children[keys[0]].doSet(keys[1:], v)
	}
}

// IsEmpty reports whether the ApiProxyStore is empty
// O(1)
func (s *ApiProxyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.proxyTree.children) <= 0
}

// Get retrieves an ApiProxy if it exists from a request path
// O(n), n => number of path segments in request path
func (s *ApiProxyFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 1 {
		return nil, errors.New("too many arguments")
	}
	path, ok := params[0].(string)
	if !ok {
		return nil, errors.New("request path was not expected type")
	}
	return s.get(path), nil
}

func (s *ApiProxyFactory) get(path string) interface{} {
	if len(s.proxyTree.children) == 0 || path == "" {
		return nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	rootNode := s.proxyTree
	for i, part := range strings.Split(path, "/") {
		if rootNode.children[part] == nil || (rootNode.children[part].value == nil && i == len(strings.Split(path, "/"))-1) {
			break
		}
		rootNode = rootNode.children[part]
	}
	if rootNode.value == nil {
		return nil
	}
	return *rootNode.value
}

// Delete will remove an ApiProxy resource
// O(n), n => number of path segments in the ApiProxy source path
func (s *ApiProxyFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	p, ok := obj.(v2.ApiProxy)
	if !ok {
		return nil, errors.New("ApiProxy was not expected type")
	}
	normalizeProxyPaths(&p)
	result := s.proxyTree.delete(strings.Split(p.Spec.Source.Path[1:], "/"))
	if result == nil {
		return nil, nil
	}
	return *result, nil
}

func (n *proxyNode) delete(segments []string) *v2.ApiProxy {
	if len(segments) == 0 {
		tmp := n.value
		n.value = nil
		return tmp
	}
	result := n.children[segments[0]].delete(segments[1:])
	if len(n.children[segments[0]].children) == 0 && n.children[segments[0]].value == nil {
		delete(n.children, segments[0])
	}
	return result
}

func normalizeProxyPaths(p *v2.ApiProxy) {
	(*p).Spec.Source.Path = utils.NormalizeURLPath(p.Spec.Source.Path)
	(*p).Spec.Target.Path = utils.NormalizeURLPath(p.Spec.Target.Path)
}
