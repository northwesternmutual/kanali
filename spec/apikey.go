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
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// APIKeyList represents a list of APIKeys
type APIKeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []APIKey `json:"items"`
}

// APIKey represents the TPR for an APIKey
type APIKey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              APIKeySpec `json:"spec"`
}

// DeepCopyObject TODO
func (p *APIKey) DeepCopyObject() runtime.Object {
	return p
}

// DeepCopyObject TODO
func (p *APIKeyList) DeepCopyObject() runtime.Object {
	return p
}

// APIKeySpec represents the data fields for the APIKey TPR
type APIKeySpec struct {
	APIKeyData string `json:"data"`
}

// KeyFactory is factory that implements a concurrency safe store for Kanali APIKeys
type KeyFactory struct {
	mutex  sync.RWMutex
	keyMap map[string]APIKey
}

// KeyStore holds all Kanali APIKeys that Kanali has discovered
// in a cluster. It should not be mutated directly!
var KeyStore *KeyFactory

// APIKeyDecryptionKey references the rsa private key that Kanali
// will use to decrypt the data in an APIKey spec
var APIKeyDecryptionKey *rsa.PrivateKey

func init() {
	KeyStore = &KeyFactory{sync.RWMutex{}, map[string]APIKey{}}
}

// Clear will remove all keys from the store
func (s *KeyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.keyMap {
		delete(s.keyMap, k)
	}
}

// Update will update an APIKeyBinding
func (s *KeyFactory) Update(old, new interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key, ok := old.(APIKey)
	if !ok {
		return errors.New("grrr - you're only allowed add api keys to the api key store.... duh")
	}
	return s.set(key)
}

// Set takes a APIKey and either adds it to the store
// or updates it
func (s *KeyFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key, ok := obj.(APIKey)
	if !ok {
		return errors.New("grrr - you're only allowed add api keys to the api key store.... duh")
	}
	return s.set(key)
}

func (s *KeyFactory) set(key APIKey) error {
	s.keyMap[key.Spec.APIKeyData] = key
	return nil
}

// Get retrieves a particual key in the store. If not found, nil is returned.
func (s *KeyFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 1 {
		return nil, errors.New("should only pass the name of the api key")
	}
	name, ok := params[0].(string)
	if !ok {
		return nil, errors.New("when retrieving a key, use the keys name")
	}
	k, ok := s.keyMap[name]
	if !ok {
		return nil, nil
	}
	return k, nil
}

// Delete will remove a particular key from the store
func (s *KeyFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key, ok := obj.(APIKey)
	if !ok {
		return nil, errors.New("there's no way this api key could've gotten in here")
	}
	actual, ok := s.keyMap[key.Spec.APIKeyData]
	if !ok {
		return nil, nil
	}
	delete(s.keyMap, key.Spec.APIKeyData)
	return actual, nil
}

// IsEmpty reports whether the key store is empty
func (s *KeyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.keyMap) == 0
}

// Decrypt decrypts the data in an APIKey
func (k *APIKey) Decrypt() error {
	cipherText, err := hex.DecodeString(k.Spec.APIKeyData)
	if err != nil {
		return err
	}
	unencryptedAPIKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, APIKeyDecryptionKey, cipherText, []byte("kanali"))
	if err != nil {
		return err
	}
	k.Spec.APIKeyData = string(unencryptedAPIKey)
	return nil
}
