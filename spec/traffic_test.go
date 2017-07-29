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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestCompareTime(t *testing.T) {

	assert := assert.New(t)

	t1 := time.Now()
	t2 := t1

	d1, err := time.ParseDuration("-1h")
	if err != nil {
		assert.Fail("parse duration should work")
	}

	t3 := t1.Add(d1)

	d2, err := time.ParseDuration("1h")
	if err != nil {
		assert.Fail("parse duration should work")
	}

	t4 := t1.Add(d2)

	assert.Equal(0, compareTime(t1, t2), "times should be equal")
	assert.Equal(-1, compareTime(t3, t2), "time t3 should be less than t2")
	assert.Equal(1, compareTime(t4, t2), "time t4 should be greater than t2")

}

func TestTimeMinusOneUnit(t *testing.T) {

	assert := assert.New(t)

	t1 := time.Now()

	d1, err := time.ParseDuration("-1h")
	if err != nil {
		assert.Fail("parse duration should work")
	}

	t2 := t1.Add(d1)
	t3, err := timeMinusOneUnit(t1, "hour")

	assert.Equal(t2, t3, "times should be equal")

	d2, err := time.ParseDuration("-1m")
	if err != nil {
		assert.Fail("parse duration should work")
	}

	t4 := t1.Add(d2)
	t5, err := timeMinusOneUnit(t1, "minute")

	assert.Equal(t4, t5, "times should be equal")

	d3, err := time.ParseDuration("-1s")
	if err != nil {
		assert.Fail("parse duration should work")
	}

	t6 := t1.Add(d3)
	t7, err := timeMinusOneUnit(t1, "second")

	assert.Equal(t6, t7, "times should be equal")

	_, err = timeMinusOneUnit(t1, "wrong")

	assert.NotNil(err, "should throw error")

}

func TestGetTrafficVolume(t *testing.T) {

	assert := assert.New(t)
	timeArray := []time.Time{}

	// TEST NORMAL CASES
	for i := 0; i < 2; i++ {
		for j := 0; j < 60; j++ {
			tmpTime, _ := time.Parse("Mon Jan 2 15:04:05:00 -0700 MST 2006", fmt.Sprintf("Sun Jun 12 %02d:%02d:00:00 -0000 CST 2017", i+1, j))
			timeArray = append(timeArray, tmpTime)
		}
	}

	mockCurrTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Sun Jun 12 03:00:00 -0000 CST 2017")

	// this implies that in the last hour, there were 59 requests
	assert.Equal(59, getTrafficVolume(timeArray, "hour", mockCurrTime, 0, len(timeArray)), "volume not what expected")

	timeArray = []time.Time{}

	for i := 0; i < 2; i++ {
		for j := 0; j < 60; j++ {
			for k := 0; k < 60; k++ {
				tmpTime, _ := time.Parse("Mon Jan 2 15:04:05:00 -0700 MST 2006", fmt.Sprintf("Sun Jun 12 %02d:%02d:%02d:00 -0000 CST 2017", i+1, j, k))
				timeArray = append(timeArray, tmpTime)
			}
		}
	}

	// this implies that in the last minute, there were 59 requests
	assert.Equal(59, getTrafficVolume(timeArray, "minute", mockCurrTime, 0, len(timeArray)), "volume not what expected")

	timeArray = []time.Time{}

	for i := 0; i < 2; i++ {
		for j := 0; j < 60; j++ {
			for k := 0; k < 60; k++ {
				for l := 0; l < 100; l += 2 {
					tmpTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", fmt.Sprintf("Sun Jun 12 %02d:%02d:%02d.%02d -0000 CST 2017", i+1, j, k, l))
					timeArray = append(timeArray, tmpTime)
				}
			}
		}
	}

	mockCurrTime, _ = time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Sun Jun 12 03:00:00.01 -0000 CST 2017")

	// this implies that in the last second, there were 59 requests
	assert.Equal(49, getTrafficVolume(timeArray, "second", mockCurrTime, 0, len(timeArray)), "volume not what expected")
	assert.Equal(len(timeArray), getTrafficVolume(timeArray, "frank", mockCurrTime, 0, len(timeArray)), "expected length because of error")

	// TEST EDGE CASES
	assert.Equal(0, getTrafficVolume(nil, "second", mockCurrTime, 0, len(timeArray)), "volume not what expected")

	tmpTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 2:05:00.00 -0000 CST 2017")
	timeArray = []time.Time{tmpTime}

	assert.Equal(1, getTrafficVolume(timeArray, "hour", mockCurrTime, 0, len(timeArray)), "volume not what expected")
	assert.Equal(0, getTrafficVolume(timeArray, "minute", mockCurrTime, 0, len(timeArray)), "volume not what expected")
	assert.Equal(0, getTrafficVolume(timeArray, "second", mockCurrTime, 0, len(timeArray)), "volume not what expected")

}

func TestGetTrafficStore(t *testing.T) {

	// setup
	assert := assert.New(t)
	TrafficStore := TrafficStore
	message := "store is not empty"

	// test
	TrafficStore.Clear()
	assert.Equal(0, len(TrafficStore), message)

}

func TestAddTraffic(t *testing.T) {

	// setup
	assert := assert.New(t)
	TrafficStore := TrafficStore
	currTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 2:05:00.00 -0000 CST 2017")

	// test
	TrafficStore.Clear()
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	TrafficStore.AddTraffic("namespace-one", "proxy-two", "key-one", currTime)
	TrafficStore.AddTraffic("namespace-one", "proxy-three", "key-one", currTime)
	TrafficStore.AddTraffic("namespace-one", "proxy-three", "key-two", currTime)
	assert.Equal(1, len(TrafficStore))
	assert.Equal(3, len(TrafficStore["namespace-one"]))
	assert.Equal(2, len(TrafficStore["namespace-one"]["proxy-three"]))

}

func TestTrafficStoreClear(t *testing.T) {

	// setup
	assert := assert.New(t)
	TrafficStore := TrafficStore
	currTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 2:05:00.00 -0000 CST 2017")

	// test
	TrafficStore.Clear()
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	TrafficStore.Clear()
	assert.Equal(0, len(TrafficStore), "store should be empty")

}

func TestHasTraffic(t *testing.T) {

	assert := assert.New(t)
	TrafficStore := TrafficStore
	currTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 2:05:00.00 -0000 CST 2017")

	// test
	TrafficStore.Clear()
	assert.False(TrafficStore.hasTraffic(getTestAPIKeyBinding(), "key-one"))
	TrafficStore.AddTraffic("namespace-one", "proxy-two", "key-two", currTime)
	assert.False(TrafficStore.hasTraffic(getTestAPIKeyBinding(), "key-one"))
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-two", currTime)
	assert.False(TrafficStore.hasTraffic(getTestAPIKeyBinding(), "key-one"))
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	assert.True(TrafficStore.hasTraffic(getTestAPIKeyBinding(), "key-one"))
	TrafficStore.AddTraffic("namespace-one", "proxy-two", "key-one", currTime)
	assert.True(TrafficStore.hasTraffic(getTestAPIKeyBinding(), "key-one"))

}

func TestIsQuotaViolated(t *testing.T) {

	assert := assert.New(t)
	TrafficStore := TrafficStore
	currTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 2:05:00.00 -0000 CST 2017")
	TrafficStore.Clear()

	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	assert.False(TrafficStore.IsQuotaViolated(getTestAPIKeyBinding(), "key-one"))

	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	assert.True(TrafficStore.IsQuotaViolated(getTestAPIKeyBinding(), "key-one"))

	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", currTime)
	assert.True(TrafficStore.IsQuotaViolated(getTestAPIKeyBinding(), "key-one"))

	assert.True(TrafficStore.IsQuotaViolated(getTestAPIKeyBinding(), "key-frank"))
	testBinding := getTestAPIKeyBinding()
	testBinding.Spec.Keys[0].Quota = 0

	assert.False(TrafficStore.IsQuotaViolated(testBinding, "key-one"))

	TrafficStore.Clear()
	assert.False(TrafficStore.IsQuotaViolated(getTestAPIKeyBinding(), "key-one"))

}

func TestIsRateLimitViolated(t *testing.T) {

	assert := assert.New(t)
	TrafficStore := TrafficStore
	currTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 3:05:54.10 -0000 CST 2017")
	TrafficStore.Clear()
	testBinding := getTestAPIKeyBinding()
	testBinding.Spec.Keys[0].Rate = &Rate{3, "minute"}

	assert.False(TrafficStore.IsRateLimitViolated(testBinding, "key-one", currTime))

	tmpTime, _ := time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 3:05:04.14 -0000 CST 2017")
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", tmpTime)
	assert.False(TrafficStore.IsRateLimitViolated(testBinding, "key-one", currTime))

	tmpTime, _ = time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 3:05:06.01 -0000 CST 2017")
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", tmpTime)
	assert.False(TrafficStore.IsRateLimitViolated(testBinding, "key-one", currTime))

	tmpTime, _ = time.Parse("Mon Jan 2 15:04:05.00 -0700 MST 2006", "Sun Jun 12 3:05:08.99 -0000 CST 2017")
	TrafficStore.AddTraffic("namespace-one", "proxy-one", "key-one", tmpTime)
	assert.True(TrafficStore.IsRateLimitViolated(testBinding, "key-one", currTime))

	testBinding.Spec.Keys[0].Rate = &Rate{}
	assert.False(TrafficStore.IsRateLimitViolated(testBinding, "key-one", currTime))

	assert.True(TrafficStore.IsRateLimitViolated(testBinding, "key-two", currTime))

}

func getTestAPIKeyBinding() APIKeyBinding {
	return APIKeyBinding{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "namespace-one",
		},
		Spec: APIKeyBindingSpec{
			APIProxyName: "proxy-one",
			Keys: []Key{
				{
					Name:        "key-one",
					DefaultRule: Rule{},
					Quota:       2,
					Subpaths: []*Path{
						{
							Path: "/foo",
							Rule: Rule{
								Global: true,
							},
						},
						{
							Path: "foo/bar",
							Rule: Rule{
								Granular: &GranularProxy{
									Verbs: []string{
										"POST",
										"GET",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
