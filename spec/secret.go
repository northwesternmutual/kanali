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
	"k8s.io/client-go/pkg/api/v1"
)

// SecretFactory is factory that implements a concurrency safe store for Kubernetes secrets
type SecretFactory struct {
	mutex     sync.RWMutex
	secretMap map[string]map[string]v1.Secret
}

// SecretStore holds all Kubernetes secrets that Kanali has discovered
// in a cluster. It should not be mutated directly!
var SecretStore *SecretFactory

func init() {
	SecretStore = &SecretFactory{sync.RWMutex{}, map[string]map[string]v1.Secret{}}
}

// Clear will remove all secrets from the store
func (s *SecretFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.secretMap {
		delete(s.secretMap, k)
	}
}

// Update will update a secret
func (s *SecretFactory) Update(old, new interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	secret, ok := old.(v1.Secret)
	if !ok {
		return errors.New("grrr - you're only allowed add secrets to the secrets store.... duh")
	}
	return s.set(secret)
}

// Set takes a Secret and either adds it to the store
// or updates it
func (s *SecretFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	secret, ok := obj.(v1.Secret)
	if !ok {
		return errors.New("grrr - you're only allowed add secrets to the secrets store.... duh")
	}
	return s.set(secret)
}

func (s *SecretFactory) set(secret v1.Secret) error {
	logrus.Infof("Adding new Secret named %s", secret.ObjectMeta.Name)
	if _, ok := s.secretMap[secret.ObjectMeta.Namespace]; ok {
		s.secretMap[secret.ObjectMeta.Namespace][secret.ObjectMeta.Name] = secret
	} else {
		s.secretMap[secret.ObjectMeta.Namespace] = map[string]v1.Secret{
			secret.ObjectMeta.Name: secret,
		}
	}
	return nil
}

// Get retrieves a particual secret in the store. If not found, nil is returned.
func (s *SecretFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return nil, errors.New("should should take 2 params, name and namespace")
	}
	name, ok := params[0].(string)
	if !ok {
		return nil, errors.New("secret name must be of type string")
	}
	namespace, ok := params[1].(string)
	if !ok {
		return nil, errors.New("secret namespace must be of type string")
	}
	secret, ok := s.secretMap[namespace][name]
	if !ok {
		return nil, nil
	}
	return secret, nil
}

// Delete will remove a particular secret from the store
func (s *SecretFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if obj == nil {
		return nil, nil
	}
	secret, ok := obj.(v1.Secret)
	if !ok {
		return nil, errors.New("there's no way this secret could've gotten in here")
	}
	if _, ok = s.secretMap[secret.ObjectMeta.Namespace]; !ok {
		return nil, nil
	}
	oldSecret, ok := s.secretMap[secret.ObjectMeta.Namespace][secret.ObjectMeta.Name]
	if !ok {
		return nil, nil
	}
	logrus.Debugf("deleting secret object %s", oldSecret.ObjectMeta.Name)
	delete(s.secretMap[secret.ObjectMeta.Namespace], secret.ObjectMeta.Name)
	if len(s.secretMap[secret.ObjectMeta.Namespace]) == 0 {
		delete(s.secretMap, secret.ObjectMeta.Namespace)
	}
	return oldSecret, nil
}

// IsEmpty reports whether the secret store is empty
func (s *SecretFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.secretMap) == 0
}

// X509KeyPair creates a tls.Certificate from the tls data in
// a Kubernetes secret of type kubernetes.io/tls
func X509KeyPair(s v1.Secret) (*tls.Certificate, error) {
	pair, err := tls.X509KeyPair(s.Data["tls.crt"], s.Data["tls.key"])
	if err != nil {
		return nil, err
	}
	return &pair, err
}
