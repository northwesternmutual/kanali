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
	"strings"
	"sync"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/utils"
)

type ApiProxyStoreInterface interface {
	Set(apiProxy *v2.ApiProxy)
	Update(old, new *v2.ApiProxy)
	Get(path, vhost string) *v2.ApiProxy
	Delete(apiProxy *v2.ApiProxy) *v2.ApiProxy
	Clear()
	IsEmpty() bool
}

type apiProxyFactory struct {
	mutex     sync.RWMutex
	proxyTree *proxyNode
}

type proxyNode struct {
	children map[string]*proxyNode
	value    *proxyNodeValue
}

type proxyNodeValue struct {
	vhosts map[string]*v2.ApiProxy
	global *v2.ApiProxy
}

var (
	apiProxyStore = &apiProxyFactory{sync.RWMutex{}, &proxyNode{}}
)

func ApiProxyStore() ApiProxyStoreInterface {
	return apiProxyStore
}

// Clear will remove all ApiProxy resources
// O(1)
func (s *apiProxyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*(s.proxyTree) = proxyNode{}
}

// Update will update an ApiProxy resource.
// O(n), n => number of path segments in the ApiProxy source path
func (s *apiProxyFactory) Update(old, new *v2.ApiProxy) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.update(old, new)
}

func (s *apiProxyFactory) update(old, new *v2.ApiProxy) {
	normalizeProxyPaths(old)
	normalizeProxyPaths(new)
	existing := s.get(new.Spec.Source.Path, new.Spec.Source.VirtualHost)
	if existing != nil && (new.GetName() != existing.GetName() || new.GetNamespace() != existing.GetNamespace() || new.Spec.Source.VirtualHost != existing.Spec.Source.VirtualHost) {
		return
	}

	s.proxyTree.doSet(strings.Split(new.Spec.Source.Path[1:], "/"), new)
	if old.Spec.Source.Path != new.Spec.Source.Path || old.Spec.Source.VirtualHost != new.Spec.Source.VirtualHost {
		s.proxyTree.delete(strings.Split(old.Spec.Source.Path[1:], "/"), old.Spec.Source.VirtualHost)
	}
}

// Set adds an ApiProxy resource
// O(n), n => number of path segments in the ApiProxy source path
func (s *apiProxyFactory) Set(apiProxy *v2.ApiProxy) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set(apiProxy)
}

func (s *apiProxyFactory) set(apiProxy *v2.ApiProxy) {
	normalizeProxyPaths(apiProxy)
	// edge case
	if apiProxy.Spec.Source.Path == "/" {
		s.proxyTree.value = new(proxyNodeValue)
		if apiProxy.Spec.Source.VirtualHost == "" {
			s.proxyTree.value.global = apiProxy
		} else {
			s.proxyTree.value.vhosts = map[string]*v2.ApiProxy{
				apiProxy.Spec.Source.VirtualHost: apiProxy,
			}
		}
	} else {
		s.proxyTree.doSet(strings.Split(apiProxy.Spec.Source.Path[1:], "/"), apiProxy)
	}
}

func (n *proxyNode) doSet(keys []string, v *v2.ApiProxy) {
	if n.children == nil {
		n.children = map[string]*proxyNode{}
	}
	if n.children[keys[0]] == nil {
		n.children[keys[0]] = &proxyNode{}
	}
	if len(keys) < 2 {
		if n.children[keys[0]].value == nil {
			n.children[keys[0]].value = new(proxyNodeValue)
		}
		if v.Spec.Source.VirtualHost == "" {
			n.children[keys[0]].value.global = v
		} else {
			if n.children[keys[0]].value.vhosts == nil {
				n.children[keys[0]].value.vhosts = map[string]*v2.ApiProxy{}
			}
			n.children[keys[0]].value.vhosts[v.Spec.Source.VirtualHost] = v
		}
	} else {
		n.children[keys[0]].doSet(keys[1:], v)
	}
}

// IsEmpty reports whether the ApiProxyStore is empty
// O(1)
func (s *apiProxyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.proxyTree.children) <= 0
}

// Get retrieves an ApiProxy if it exists given a request path
// and virtual host (if any).
// ϴ(logn), n => cardinality of ApiProxy resources in the store unique to source path
func (s *apiProxyFactory) Get(path, vhost string) *v2.ApiProxy {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.get(path, vhost)
}

func (s *apiProxyFactory) get(path, vhost string) *v2.ApiProxy {
	if len(path) > 0 && path[0] == '/' {
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

	if rootNode.value.global != nil {
		return rootNode.value.global
	}

	return rootNode.value.vhosts[vhost]
}

// Delete will remove an ApiProxy resource
// O(n), n => number of path segments in the ApiProxy source path
func (s *apiProxyFactory) Delete(apiProxy *v2.ApiProxy) *v2.ApiProxy {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.delete(apiProxy)
}

func (s *apiProxyFactory) delete(apiProxy *v2.ApiProxy) *v2.ApiProxy {
	if apiProxy == nil {
		return nil
	}
	normalizeProxyPaths(apiProxy)
	return s.proxyTree.delete(strings.Split(apiProxy.Spec.Source.Path[1:], "/"), apiProxy.Spec.Source.VirtualHost)
}

func (n *proxyNode) delete(segments []string, vhost string) *v2.ApiProxy {
	if len(segments) == 0 || (len(segments) == 1 && segments[0] == "") {
		if n.value != nil && n.value.global != nil {
			tmp := n.value.global
			n.value = nil
			return tmp
		}
		if n.value != nil && n.value.vhosts != nil {
			if val, ok := n.value.vhosts[vhost]; ok {
				tmp := val
				delete(n.value.vhosts, vhost)
				if len(n.value.vhosts) < 1 {
					n.value = nil
				}
				return tmp
			}
		}
		return nil
	}
	result := n.children[segments[0]].delete(segments[1:], vhost)
	if len(n.children[segments[0]].children) == 0 && n.children[segments[0]].value == nil {
		delete(n.children, segments[0])
	}
	return result
}

func normalizeProxyPaths(p *v2.ApiProxy) {
	(*p).Spec.Source.Path = utils.NormalizeURLPath(p.Spec.Source.Path)
	(*p).Spec.Target.Path = utils.NormalizeURLPath(p.Spec.Target.Path)
}
