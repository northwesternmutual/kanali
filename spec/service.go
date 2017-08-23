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
	"net/http"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/pkg/api/v1"
)

// Service in an internal representation of a Kubernetes Service
type Service struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	ClusterIP string `json:"clusterIP,omitempty"`
	Port      int64  `json:"port,omitempty"`
	Labels    Labels `json:"labels,omitempty"`
}

// Labels represents labels on a Kubernetes service
type Labels []Label

// Label represents a Kubernetes service label.
// It also represents a Kubernetes as defined in
// a proxy spec
type Label struct {
	Name   string `json:"name,omitempty"`
	Header string `json:"header,omitempty"`
	Value  string `json:"value,omitempty"`
}

type services []Service

// ServiceFactory is factory that implements a concurrency safe store for Kubernetes services
type ServiceFactory struct {
	mutex      sync.RWMutex
	serviceMap map[string]services
}

// ServiceStore holds all Kubernetes services that Kanali has discovered
// in a cluster. It should not be mutated directly!
var ServiceStore *ServiceFactory

func init() {
	ServiceStore = &ServiceFactory{sync.RWMutex{}, map[string]services{}}
}

// Clear will remove all services from the store
func (s *ServiceFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for svc := range s.serviceMap {
		delete(s.serviceMap, svc)
	}
}

// Set takes a Service and either adds it to the store
// or updates it
func (s *ServiceFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	service, ok := obj.(Service)
	if !ok {
		return errors.New("grrr - you're only allowed add services to the services store.... duh")
	}
	logrus.Debugf("adding service object %s", service.Name)
	if s.serviceMap[service.Namespace] == nil {
		s.serviceMap[service.Namespace] = []Service{service}
		return nil
	}
	for i, svc := range s.serviceMap[service.Namespace] {
		if svc.Name == service.Name {
			s.serviceMap[service.Namespace][i] = service
			return nil
		}
	}
	s.serviceMap[service.Namespace] = append(s.serviceMap[service.Namespace], service)
	return nil
}

// IsEmpty reports whether the service store is empty
func (s *ServiceFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.serviceMap) == 0
}

// Get retrieves a particual service in the store. If not found, nil is returned.
func (s *ServiceFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return nil, errors.New("getting a service requires 2 parameters")
	}
	svc, ok := params[0].(Service)
	if !ok {
		return nil, errors.New("first argument should be a service")
	}
	headers, ok := params[1].(http.Header)
	if !ok {
		if params[1] != nil {
			return nil, errors.New("second argument should either be nil or http.Header")
		}
	}
	if _, ok := s.serviceMap[svc.Namespace]; !ok {
		return nil, nil
	}
	for _, item := range s.serviceMap[svc.Namespace] {
		if svc.Name == "" {
			if svc.Labels.isSubset(item.Labels, headers) {
				return item, nil
			}
		} else {
			if item.Name == svc.Name {
				return item, nil
			}
		}
	}
	return nil, nil
}

// Contains reports whether the service store contains a particular service
// TODO
func (s *ServiceFactory) Contains(params ...interface{}) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return false, errors.New("method not yet implemented")
}

// Delete will remove a particular service from the store
func (s *ServiceFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	service, ok := obj.(Service)
	if !ok {
		return nil, errors.New("there's no way this service could've gotten in here")
	}
	if _, ok := s.serviceMap[service.Namespace]; !ok {
		return nil, nil
	}
	i, svc := s.serviceMap[service.Namespace].indexOf(service)
	if i < 0 {
		return nil, nil
	}
	logrus.Debugf("deleting service object %s", service.Name)
	if len(s.serviceMap[service.Namespace]) == 1 {
		delete(s.serviceMap, service.Namespace)
	} else {
		s.serviceMap[service.Namespace] = append(s.serviceMap[service.Namespace][:i], s.serviceMap[service.Namespace][i+1:]...)
	}
	return *svc, nil
}

// CreateService transforms a Kubernetes service
// of type v1.Service into type Service
func CreateService(s v1.Service) Service {
	l := Labels{}
	for k, v := range s.ObjectMeta.Labels {
		l = append(l, Label{Name: k, Value: v})
	}
	return Service{
		Name:      s.ObjectMeta.Name,
		Namespace: s.ObjectMeta.Namespace,
		ClusterIP: s.Spec.ClusterIP,
		Labels:    l,
	}
}

func (one Labels) isSubset(other Labels, headers http.Header) bool {
	for _, item := range one {
		if !other.contains(item, headers) {
			return false
		}
	}
	return true
}

func (one Labels) contains(other Label, headers http.Header) bool {
	for _, item := range one {
		if other.equals(item, headers) {
			return true
		}
	}
	return false
}

func (one Label) equals(other Label, headers http.Header) bool {
	// is the name of the label the same
	if strings.Compare(strings.ToLower(one.Name), strings.ToLower(other.Name)) == 0 {
		// is header specified - only on 'one'
		if len(one.Header) > 0 {
			if len(headers.Get(one.Header)) > 0 {
				return strings.Compare(strings.ToLower(headers.Get(one.Header)), strings.ToLower(other.Value)) == 0
			}
			// the header that we are looking for was not part of the request headers
			// so now we need to check the default values and match against those
			if len(viper.GetString(fmt.Sprintf("headers.%s", one.Header))) > 0 {
				return strings.Compare(strings.ToLower(viper.GetString(fmt.Sprintf("headers.%s", one.Header))), strings.ToLower(other.Value)) == 0
			}
			return false
		}
		return strings.Compare(other.Value, one.Value) == 0
	}
	return false
}

func (a services) indexOf(s Service) (int, *Service) {
	// since the name of a service in a namespace
	// must be unique, that's all we have to look at
	for i, item := range a {
		if strings.Compare(item.Name, s.Name) == 0 {
			return i, &item
		}
	}
	return -1, nil
}
