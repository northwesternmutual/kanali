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
	"fmt"
	"strings"
	"time"
)

// TrafficFactory is factory that implements a data structure for Kanali traffic
type TrafficFactory map[string]map[string]map[string][]time.Time

// TrafficStore holds all Kanali traffic. It should not be mutated directly!
var TrafficStore TrafficFactory

func init() {
	TrafficStore = map[string]map[string]map[string][]time.Time{}
}

// Clear will remove all entries from the traffic store
func (store TrafficFactory) Clear() {
	for k := range store {
		delete(store, k)
	}
}

// AddTraffic will add a traffic point to the traffic store
func (store TrafficFactory) AddTraffic(namespace, proxyName, keyName string, currTime time.Time) {

	if _, ok := store[namespace]; !ok {
		store[namespace] = map[string]map[string][]time.Time{
			proxyName: {
				keyName: {currTime},
			},
		}
		return
	}

	if _, ok := store[namespace][proxyName]; !ok {
		store[namespace][proxyName] = map[string][]time.Time{
			keyName: {currTime},
		}
		return
	}

	if _, ok := store[namespace][proxyName][keyName]; !ok {
		store[namespace][proxyName][keyName] = []time.Time{currTime}
		return
	}

	store[namespace][proxyName][keyName] = append(store[namespace][proxyName][keyName], currTime)

}

// IsQuotaViolated will see whether a quota limit has been reached
func (store TrafficFactory) IsQuotaViolated(binding APIKeyBinding, keyName string) bool {

	for _, key := range binding.Spec.Keys {

		if key.Name != keyName {
			continue
		}

		if key.Quota == 0 {
			return false
		}

		if !store.hasTraffic(binding, keyName) {
			return false
		}

		return len(store[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName]) >= key.Quota

	}

	return true

}

// IsRateLimitViolated wee see whether a rate limit has been reached
func (store TrafficFactory) IsRateLimitViolated(binding APIKeyBinding, keyName string, currTime time.Time) bool {

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

		if !store.hasTraffic(binding, keyName) {
			return false
		}

		return getTrafficVolume(store[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName], key.Rate.Unit, currTime, 0, len(store[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName])) >= key.Rate.Amount

	}

	return true

}

func (store TrafficFactory) hasTraffic(binding APIKeyBinding, keyName string) bool {

	if _, ok := store[binding.ObjectMeta.Namespace]; !ok {
		return false
	}

	if _, ok := store[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName]; !ok {
		return false
	}

	if _, ok := store[binding.ObjectMeta.Namespace][binding.Spec.APIProxyName][keyName]; !ok {
		return false
	}

	return true

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
