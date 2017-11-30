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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type TrafficStoreInterface interface {
	Set(tp *TrafficPoint)
	Clear()
	IsEmpty() bool
	IsRateLimitViolated(proxy *v2.ApiProxy, binding *v2.ApiKeyBinding, keyName string, currTime time.Time) bool
	TrafficStoreExpansion
}

type trafficByApiKey map[string][]time.Time
type trafficByApiProxy map[string]trafficByApiKey
type trafficByNamespace map[string]trafficByApiProxy

type trafficFactory struct {
	mutex      sync.RWMutex
	trafficMap trafficByNamespace
}

var (
	trafficStore = &trafficFactory{sync.RWMutex{}, make(trafficByNamespace)}
)

func TrafficStore() TrafficStoreInterface {
	return trafficStore
}

// Clear will remove all entries from the traffic store
func (s *trafficFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.trafficMap {
		delete(s.trafficMap, k)
	}
}

// Set takes a traffic point and either adds it to the store
func (s *trafficFactory) Set(tp *TrafficPoint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.doSet(tp)
}

func (s *trafficFactory) doSet(tp *TrafficPoint) {
	if _, ok := s.trafficMap[tp.Namespace]; !ok {
		s.trafficMap[tp.Namespace] = make(trafficByApiProxy)
	}
	if _, ok := s.trafficMap[tp.Namespace][tp.ProxyName]; !ok {
		s.trafficMap[tp.Namespace][tp.ProxyName] = make(trafficByApiKey)
	}
	if _, ok := s.trafficMap[tp.Namespace][tp.ProxyName][tp.KeyName]; !ok {
		s.trafficMap[tp.Namespace][tp.ProxyName][tp.KeyName] = make([]time.Time, 0)
	}
	s.trafficMap[tp.Namespace][tp.ProxyName][tp.KeyName] = append(s.trafficMap[tp.Namespace][tp.ProxyName][tp.KeyName], time.Unix(0, tp.Time))
}

// IsEmpty reports whether the traffic store is empty
func (s *trafficFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.trafficMap) == 0
}

// // IsQuotaViolated will see whether a quota limit has been reached
// func (s *trafficFactory) IsQuotaViolated(proxy v2.ApiProxy, binding v2.ApiKeyBinding, keyName string) bool {
// 	s.mutex.RLock()
// 	defer s.mutex.RUnlock()
// 	for _, key := range binding.Spec.Keys {
// 		if key.Name != keyName {
// 			continue
// 		}
// 		if key.Quota == 0 {
// 			return false
// 		}
// 		result := s.contains(proxy, binding, keyName)
// 		if !result {
// 			return false
// 		}
// 		return len(s.trafficMap[binding.ObjectMeta.Namespace][proxy.ObjectMeta.Name][keyName]) >= key.Quota
// 	}
// 	return true
// }

// IsRateLimitViolated wee see whether a rate limit has been reached
func (s *trafficFactory) IsRateLimitViolated(proxy *v2.ApiProxy, binding *v2.ApiKeyBinding, keyName string, currTime time.Time) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, key := range binding.Spec.Keys {
		if key.Name != keyName {
			continue
		}
		if key.Rate.Amount < 1 {
			return false
		}
		result := s.contains(proxy, binding, keyName)
		if !result {
			return false
		}
		if key.Rate.Unit == "" { // quota
			return len(s.trafficMap[binding.ObjectMeta.Namespace][proxy.ObjectMeta.Name][keyName]) >= key.Rate.Amount
		}
		return getTrafficVolume(s.trafficMap[binding.ObjectMeta.Namespace][proxy.ObjectMeta.Name][keyName], key.Rate.Unit, currTime, 0, len(s.trafficMap[binding.ObjectMeta.Namespace][proxy.ObjectMeta.Name][keyName])) >= key.Rate.Amount
	}
	return true
}

// Contains reports whether the traffic store has any traffic for a given proxy/name combination
func (s *trafficFactory) contains(proxy *v2.ApiProxy, binding *v2.ApiKeyBinding, keyName string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	_, ok := s.trafficMap[binding.ObjectMeta.Namespace][proxy.ObjectMeta.Name][keyName]
	return ok
}

func getTrafficVolume(arr []time.Time, unit string, currTime time.Time, low, high int) int {
	if arr == nil {
		return 0
	}
	if high <= low {
		return len(arr[low:])
	}
	mid := (low + high) / 2
	tMinusOne, err := timeMinusOneUnit(currTime, unit)
	if err != nil {
		return len(arr)
	}
	if compareTime(tMinusOne, arr[mid]) < 0 {
		return getTrafficVolume(arr, unit, currTime, 0, mid)
	}
	return getTrafficVolume(arr, unit, currTime, mid+1, high)
}

func timeMinusOneUnit(t time.Time, unit string) (time.Time, error) {
	var newTime time.Time
	d, err := time.ParseDuration(fmt.Sprintf("-1%s", strings.ToLower(string(unit[0]))))
	if err != nil {
		return newTime, err
	}
	return t.Add(d), nil
}

func compareTime(t1, t2 time.Time) int {
	if t1.Equal(t2) {
		return 0
	}
	if t1.Before(t2) {
		return -1
	}
	return 1
}
