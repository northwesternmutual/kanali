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
	"sync"
  "bytes"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type ApiKeyStoreInterface interface {
	Set(apiKey *v2.ApiKey)
	Update(old, new *v2.ApiKey)
	Get(data []byte) *v2.ApiKey
	Delete(apiKey *v2.ApiKey) *v2.ApiKey
	Clear()
	IsEmpty() bool
}

type apiKeyFactory struct {
	mutex  sync.RWMutex
	keyMap map[string]v2.ApiKey
}

var (
	apiKeyStore = &apiKeyFactory{sync.RWMutex{}, map[string]v2.ApiKey{}}
)

func ApiKeyStore() ApiKeyStoreInterface {
	return apiKeyStore
}

// Clear will remove all ApiKey resources
// O(n), n => the cartesian product of all ApiKey resources and ApiKey revisions
func (s *apiKeyFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.keyMap {
		delete(s.keyMap, k)
	}
}

// Update will update an ApiKey resource
// O(nlogn),
//   n => max(x, y)
//   x => number of revisions in old ApiKey
//   y => number of revisions in new ApiKey
func (s *apiKeyFactory) Update(old, new *v2.ApiKey) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldRevisions := mergeSort(old.Spec.Revisions)
	newRevisions := mergeSort(new.Spec.Revisions)

	for len(oldRevisions) > 0 && len(newRevisions) > 0 {
		if bytes.Equal(oldRevisions[0].Data, newRevisions[0].Data) {
			s.keyMap[string(newRevisions[0].Data)] = *new
			oldRevisions = oldRevisions[1:]
			newRevisions = newRevisions[1:]
		} else if string(oldRevisions[0].Data) < string(newRevisions[0].Data) {
			delete(s.keyMap, string(oldRevisions[0].Data))
			oldRevisions = oldRevisions[1:]
		}
	}

	for i := range oldRevisions {
		delete(s.keyMap, string(oldRevisions[i].Data))
	}

	for i := range newRevisions {
		s.keyMap[string(newRevisions[i].Data)] = *new
	}
}

// Set adds an ApiKey resource
// O(n), n => number of revisions in the ApiKey
func (s *apiKeyFactory) Set(key *v2.ApiKey) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set(key)
}

func (s *apiKeyFactory) set(key *v2.ApiKey) {
	for _, revision := range key.Spec.Revisions {
		s.keyMap[string(revision.Data)] = *key
	}
}

// Get retrieves an ApiKey if present
// O(1)
func (s *apiKeyFactory) Get(data []byte) *v2.ApiKey {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	key, ok := s.keyMap[string(data)]
	if !ok {
		return nil
	}
	return &key
}

// Delete will remove an ApiKey
// O(n), n => number of revisions in the ApiKey
func (s *apiKeyFactory) Delete(key *v2.ApiKey) *v2.ApiKey {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if key == nil {
		return nil
	}
	if len(key.Spec.Revisions) < 1 {
		return nil
	}
	// Each ApiKey at each revision will be the same
	result, ok := s.keyMap[string(key.Spec.Revisions[0].Data)]
	if !ok {
		return nil
	}
	for _, revision := range key.Spec.Revisions {
		delete(s.keyMap, string(revision.Data))
	}
	return &result
}

// IsEmpty reports whether the ApiKeyStore is empty
// O(1)
func (s *apiKeyFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.keyMap) == 0
}

func mergeSort(slice []v2.Revision) []v2.Revision {
	if len(slice) < 2 {
		return slice
	}
	mid := (len(slice)) / 2
	return merge(mergeSort(slice[:mid]), mergeSort(slice[mid:]))
}

func merge(left, right []v2.Revision) []v2.Revision {
	size, i, j := len(left)+len(right), 0, 0
	slice := make([]v2.Revision, size, size)

	for k := 0; k < size; k++ {
		if i > len(left)-1 && j <= len(right)-1 {
			slice[k] = right[j]
			j++
		} else if j > len(right)-1 && i <= len(left)-1 {
			slice[k] = left[i]
			i++
		} else if string(left[i].Data) < string(right[j].Data) {
			slice[k] = left[i]
			i++
		} else {
			slice[k] = right[j]
			j++
		}
	}
	return slice
}
