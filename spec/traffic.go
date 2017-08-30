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
	"time"
)

type trafficByAPIKey map[string][]time.Time
type trafficByAPIProxy map[string]trafficByAPIKey
type trafficByNamespace map[string]trafficByAPIProxy

// TrafficFactory is factory that implements a concurrency safe store for Kanali traffic
type TrafficFactory struct {
	mutex      sync.RWMutex
	trafficMap trafficByNamespace
}

// TrafficStore holds all API traffic that Kanali has discovered
// in a cluster. It should not be mutated directly!
var TrafficStore *TrafficFactory

func init() {
	TrafficStore = &TrafficFactory{sync.RWMutex{}, make(trafficByNamespace)}
}

// Clear will remove all entries from the traffic store
func (s *TrafficFactory) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for k := range s.trafficMap {
		delete(s.trafficMap, k)
	}
}

// Set takes a traffic point and either adds it to the store
func (s *TrafficFactory) Set(obj interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.doSet(obj, time.Now())
}

func (s *TrafficFactory) doSet(obj interface{}, currTime time.Time) error {
	kgram, ok := obj.(string)
	if !ok {
		return errors.New("parameter not of type string")
	}
	nSpace, pName, keyName, err := decodeKanaliGram(kgram, ",")
	if err != nil {
		return err
	}
	if _, ok := s.trafficMap[nSpace]; !ok {
		s.trafficMap[nSpace] = make(trafficByAPIProxy)
	}
	if _, ok := s.trafficMap[nSpace][pName]; !ok {
		s.trafficMap[nSpace][pName] = make(trafficByAPIKey)
	}
	if _, ok := s.trafficMap[nSpace][pName][keyName]; !ok {
		s.trafficMap[nSpace][pName][keyName] = make([]time.Time, 0)
	}
	s.trafficMap[nSpace][pName][keyName] = append(s.trafficMap[nSpace][pName][keyName], time.Now())
	return nil
}

// IsEmpty reports whether the traffic store is empty
func (s *TrafficFactory) IsEmpty() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.trafficMap) == 0
}

// IsQuotaViolated will see whether a quota limit has been reached
func (s *TrafficFactory) IsQuotaViolated(binding APIKeyBinding, keyName string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, key := range binding.Spec.Keys {
		if key.Name != keyName {
			continue
		}
		if key.Quota == 0 {
			return false
		}
		result, err := s.contains(binding, keyName)
		if err != nil || !result {
			return false
		}
		return len(s.trafficMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName]) >= key.Quota
	}
	return true
}

// IsRateLimitViolated wee see whether a rate limit has been reached
func (s *TrafficFactory) IsRateLimitViolated(binding APIKeyBinding, keyName string, currTime time.Time) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	for _, key := range binding.Spec.Keys {
		if key.Name != keyName {
			continue
		}
		if key.Rate == nil {
			return false
		}
		if key.Rate.Amount == 0 {
			return false
		}
		result, err := s.contains(binding, keyName)
		if err != nil || !result {
			return false
		}
		return getTrafficVolume(s.trafficMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName], key.Rate.Unit, currTime, 0, len(s.trafficMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName])) >= key.Rate.Amount
	}
	return true
}

// Contains reports whether the traffic store has any traffic for a given proxy/name combination
func (s *TrafficFactory) contains(params ...interface{}) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if len(params) != 2 {
		return false, errors.New("expecting two parameters")
	}
	binding, ok := params[0].(APIKeyBinding)
	if !ok {
		return false, errors.New("expected the first parameter to be of type spec.APIKeyBinding")
	}
	keyName, ok := params[1].(string)
	if !ok {
		return false, errors.New("expected the second parameter to be of type string")
	}
	if _, ok := s.trafficMap[binding.ObjectMeta.Namespace]; !ok {
		return false, nil
	}
	if _, ok := s.trafficMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName]; !ok {
		return false, nil
	}
	if _, ok := s.trafficMap[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName]; !ok {
		return false, nil
	}
	return true, nil
}

// Delete removes all traffic for a given namespace, proxy, and key combination
// TODO
func (s *TrafficFactory) Delete(obj interface{}) (interface{}, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return nil, nil
}

// Get retrieves a set of traffic points for a unique namespace/proxy/key combination
// TODO
func (s *TrafficFactory) Get(params ...interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return nil, nil
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

func decodeKanaliGram(gram, delimiter string) (string, string, string, error) {
	arr := strings.Split(gram, delimiter)
	if len(arr) != 3 {
		return "", "", "", errors.New("kgram must have 3")
	}
	return arr[0], arr[1], arr[2], nil
}
