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
	"sync"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

// ApiKeyFactory is factory that implements a concurrency safe store for Kanali ApiKey resources
type ApiKeyFactory struct {
	mutex  sync.RWMutex
	keyMap map[string]v2.ApiKey
}

var (
	// ApiKeyStore holds all discovered ApiKey resources in an efficient data structure.
	// This variable should not be mutated directly!
	ApiKeyStore = &ApiKeyFactory{sync.RWMutex{}, map[string]v2.ApiKey{}}
)

// Clear will remove all ApiKey resources
// O(n), n => the cartesian product of all ApiKey resources and ApiKey revisions
func (s *ApiKeyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.keyMap {
		delete(s.keyMap, k)
	}
}

// Update will update an ApiKey resource
// O(n), n => number of revisions in the new ApiKey
func (s *ApiKeyFactory) Update(old, new interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, oldOk := old.(v2.ApiKey)
	if !oldOk {
		return errors.New("old ApiKey was not the expected type")
	}
	newKey, newOk := new.(v2.ApiKey)
	if !newOk {
		return errors.New("new ApiKey was not the expected type")
	}
	return s.set(newKey)
}

// Set adds an ApiKey resource
// O(n), n => number of revisions in the ApiKey
func (s *ApiKeyFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key, ok := obj.(v2.ApiKey)
	if !ok {
		return errors.New("ApiKey was not the expected type")
	}
	return s.set(key)
}

func (s *ApiKeyFactory) set(key v2.ApiKey) error {
	for _, revision := range key.Spec.Revisions {
		s.keyMap[revision.Data] = key
	}
	return nil
}

// Get retrieves an ApiKey if present
// O(1)
func (s *ApiKeyFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 1 {
		return nil, errors.New("only ApiKey data is needed")
	}
	data, ok := params[0].(string)
	if !ok {
		return nil, errors.New("ApiKey data was not the expected type")
	}
	key, ok := s.keyMap[data]
	if !ok {
		return nil, nil
	}
	return key, nil
}

// Delete will remove an ApiKey
// O(n), n => number of revisions in the ApiKey
func (s *ApiKeyFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key, ok := obj.(v2.ApiKey)
	if !ok {
		return nil, errors.New("ApiKey was not the expected type")
	}
	if len(key.Spec.Revisions) < 1 {
		return nil, errors.New("ApiKey must have at least one revision")
	}
	// Each ApiKey at each revision will be the same
	result, ok := s.keyMap[key.Spec.Revisions[0].Data]
	if !ok {
		return nil, nil
	}
	for _, revision := range key.Spec.Revisions {
		delete(s.keyMap, revision.Data)
	}
	return result, nil
}

// IsEmpty reports whether the ApiKeyStore is empty
// O(1)
func (s *ApiKeyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.keyMap) == 0
}
