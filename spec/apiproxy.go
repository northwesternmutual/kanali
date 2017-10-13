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
	"fmt"
	"strings"
	"sync"

	"github.com/northwesternmutual/kanali/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// APIProxyList represents a list of APIProxies
type APIProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIProxy `json:"items"`
}

// APIProxy represents the TPR for an APIProxy
type APIProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              APIProxySpec `json:"spec"`
}

// APIProxySpec represents the data fields for the APIProxy TPR
type APIProxySpec struct {
	Path    string   `json:"path"`
	Target  string   `json:"target,omitempty"`
	Mock    *Mock    `json:"mock,omitempty"`
	Hosts   []Host   `json:"hosts,omitempty"`
	Service Service  `json:"service,omitempty"`
	Plugins []Plugin `json:"plugins,omitempty"`
	SSL     SSL      `json:"ssl,omitempty"`
}

// Mock represents a mock configuration
type Mock struct {
	ConfigMapName string `json:"configMapName,omitempty"`
}

// Plugin defines a plugin which may be version controlled
type Plugin struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// Host represents the name and SSL object to use for SNI
type Host struct {
	Name string `json:"name"`
	SSL  SSL    `json:"ssl"`
}

// SSL defines the secret to use for certificates
type SSL struct {
	SecretName string `json:"secretName"`
}

type proxyNode struct {
	Children map[string]*proxyNode
	Value    *APIProxy
}

// ProxyFactory is factory that implements a concurrency safe store for Kanali ApiProxies
type ProxyFactory struct {
	mutex     sync.RWMutex
	proxyTree *proxyNode
}

// ProxyStore holds all Kanali ApiProxies that Kanali has discovered
// in a cluster. It should not be mutated directly!
var ProxyStore *ProxyFactory

func init() {
	ProxyStore = &ProxyFactory{sync.RWMutex{}, &proxyNode{}}
}

// Clear will remove all proxies from the store
func (s *ProxyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*(s.proxyTree) = proxyNode{}
}

// Update will update an APIProxy and preform necessary clean up of old APIProxy is necessary.
func (s *ProxyFactory) Update(old, new interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	oldProxy, ok := old.(APIProxy)
	if !ok {
		return errors.New("first parameter was not of type APIProxy")
	}
	newProxy, ok := new.(APIProxy)
	if !ok {
		return errors.New("second parameter was not of type APIProxy")
	}
	normalize(&oldProxy)
	normalize(&newProxy)
	return s.update(oldProxy, newProxy)
}

func (s *ProxyFactory) update(old, new APIProxy) error {
	untyped := s.get(new.Spec.Path)
	if untyped != nil {
		typed, ok := untyped.(APIProxy)
		if !ok {
			return errors.New("received interface not of type APIProxy")
		}
		if new.ObjectMeta.Name != typed.ObjectMeta.Name || new.ObjectMeta.Namespace != typed.ObjectMeta.Namespace {
			return errors.New("there exists an APIProxy as the targeted path - APIProxy can not be updated - consider using kanalictl to avoid this error in the future")
		}
	}

	new.Spec.Service.Namespace = new.ObjectMeta.Namespace
	s.proxyTree.doSet(strings.Split(new.Spec.Path[1:], "/"), &new)
	if old.Spec.Path != new.Spec.Path {
		s.proxyTree.delete(strings.Split(old.Spec.Path[1:], "/"))
	}
	return nil
}

// Set creates or updates an APIProxy
func (s *ProxyFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := obj.(APIProxy)
	if !ok {
		return errors.New("parameter was not of type APIProxy")
	}
	p.Spec.Service.Namespace = p.ObjectMeta.Namespace
	normalize(&p)
	s.proxyTree.doSet(strings.Split(p.Spec.Path[1:], "/"), &p)
	return nil
}

func (n *proxyNode) doSet(keys []string, v *APIProxy) {
	if n.Children == nil {
		n.Children = map[string]*proxyNode{}
	}
	if n.Children[keys[0]] == nil {
		n.Children[keys[0]] = &proxyNode{}
	}
	if len(keys) < 2 {
		n.Children[keys[0]].Value = v
	} else {
		n.Children[keys[0]].doSet(keys[1:], v)
	}
}

// IsEmpty reports whether the proxy store is empty
func (s *ProxyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.proxyTree.Children) <= 0
}

// Get retrieves a particual proxy in the store. If not found, nil is returned.
func (s *ProxyFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 1 {
		return nil, errors.New("should only pass the path of the proxy")
	}
	path, ok := params[0].(string)
	if !ok {
		return nil, errors.New("when retrieving a proxy, use the proxy path")
	}
	return s.get(path), nil
}

func (s *ProxyFactory) get(path string) interface{} {
	if len(s.proxyTree.Children) == 0 || path == "" {
		return nil
	}
	if path[0] == '/' {
		path = path[1:]
	}
	rootNode := s.proxyTree
	for i, part := range strings.Split(path, "/") {
		if rootNode.Children[part] == nil || (rootNode.Children[part].Value == nil && i == len(strings.Split(path, "/"))-1) {
			break
		}
		rootNode = rootNode.Children[part]
	}
	if rootNode.Value == nil {
		return nil
	}
	return *rootNode.Value
}

// Delete will remove a particular proxy from the store
func (s *ProxyFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	p, ok := obj.(APIProxy)
	if !ok {
		return nil, errors.New("there's no way this api proxy could've gotten in here")
	}
	normalize(&p)
	result := s.proxyTree.delete(strings.Split(p.Spec.Path[1:], "/"))
	if result == nil {
		return nil, nil
	}
	return *result, nil
}

func (n *proxyNode) delete(segments []string) *APIProxy {
	if len(segments) == 0 {
		tmp := n.Value
		n.Value = nil
		return tmp
	}
	result := n.Children[segments[0]].delete(segments[1:])
	if len(n.Children[segments[0]].Children) == 0 && n.Children[segments[0]].Value == nil {
		delete(n.Children, segments[0])
	}
	return result
}

// GetSSLCertificates retreives the SSL object for a given hostname
func (p APIProxy) GetSSLCertificates(host string) *SSL {
	for _, h := range p.Spec.Hosts {
		if strings.Compare(h.Name, host) == 0 && h.SSL != (SSL{}) {
			return &h.SSL
		}
	}
	return &p.Spec.SSL
}

// GetFileName gets the file name for a plugin.
// This is dynamic base on the plugin version used.
func (p Plugin) GetFileName() string {
	if strings.Compare(p.Version, "") != 0 {
		return fmt.Sprintf("%s_%s",
			p.Name,
			p.Version,
		)
	}
	return p.Name
}

func normalize(p *APIProxy) {
	(*p).Spec.Path = utils.NormalizeURLPath(p.Spec.Path)
	(*p).Spec.Target = utils.NormalizeURLPath(p.Spec.Target)
}
