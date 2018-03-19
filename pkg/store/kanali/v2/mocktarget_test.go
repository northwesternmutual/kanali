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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/test/builder"
)

func TestMockTargetStore(t *testing.T) {
	_, ok := MockTargetStore().(MockTargetStoreInterface)
	assert.True(t, ok)
}

func TestMockTargetSet(t *testing.T) {
	defer MockTargetStore().Clear()

	routeOne := builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie()
	routeTwo := builder.NewRoute("foo/bar").WithStatusCode(200).WithMethods("GET").NewOrDie()

	MockTargetStore().Set(builder.NewMockTarget("foo", "bar").WithRoute(routeOne).NewOrDie())

	assert.NotNil(t, MockTargetStore().Set(builder.NewMockTarget("car", "bar").NewOrDie()))

	assert.Equal(t, 1, len(mockTargetStore.mockRespTree))
	assert.Equal(t, 1, len(mockTargetStore.mockRespTree["bar"]))
	assert.Equal(t, 1, len(mockTargetStore.mockRespTree["bar"]["foo"].children))
	assert.Equal(t, routeOne, mockTargetStore.mockRespTree["bar"]["foo"].children["foo"].value)
	assert.Nil(t, mockTargetStore.mockRespTree["bar"]["foo"].value)

	MockTargetStore().Set(builder.NewMockTarget("car", "bar").WithRoute(routeOne).WithRoute(routeTwo).NewOrDie())
	assert.Equal(t, 2, len(mockTargetStore.mockRespTree["bar"]))
	assert.Equal(t, 1, len(mockTargetStore.mockRespTree["bar"]["car"].children))
	assert.Equal(t, 1, len(mockTargetStore.mockRespTree["bar"]["car"].children["foo"].children))
	assert.Equal(t, routeOne, mockTargetStore.mockRespTree["bar"]["car"].children["foo"].value)
	assert.Equal(t, routeTwo, mockTargetStore.mockRespTree["bar"]["car"].children["foo"].children["bar"].value)
}

func TestMockTargetUpdate(t *testing.T) {
	defer MockTargetStore().Clear()

	routeOld := builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie()
	routeNew := builder.NewRoute("/bar").WithStatusCode(200).WithMethods("GET").NewOrDie()
	mocktargetOld := builder.NewMockTarget("foo", "bar").WithRoute(routeOld).NewOrDie()
	mocktargetNew := builder.NewMockTarget("foo", "bar").WithRoute(routeNew).NewOrDie()

	assert.Nil(t, MockTargetStore().Set(mocktargetOld))
	assert.Nil(t, MockTargetStore().Update(mocktargetOld, mocktargetNew))

	assert.Nil(t, mockTargetStore.mockRespTree["bar"]["foo"].children["foo"])
	assert.Equal(t, routeNew, mockTargetStore.mockRespTree["bar"]["foo"].children["bar"].value)
	assert.Equal(t, routeNew, MockTargetStore().Get("bar", "foo", "/bar", "GET"))
	assert.Nil(t, MockTargetStore().Get("bar", "foo", "/foo", "GET"))
}

func TestMockTargetGet(t *testing.T) {
	defer MockTargetStore().Clear()

	route := builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie()
	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(route).NewOrDie()

	assert.Nil(t, MockTargetStore().Get("", "", "", ""))

	MockTargetStore().Set(mocktarget)
	assert.Equal(t, route, MockTargetStore().Get("bar", "foo", "/foo", "GET"))
	assert.Equal(t, route, MockTargetStore().Get("bar", "foo", "foo/bar", "GET"))
	assert.Nil(t, MockTargetStore().Get("bar", "foo", "/foo", "POST"))
	assert.Nil(t, MockTargetStore().Get("bar", "foo", "/", "GET"))
}

func TestMockTargetDelete(t *testing.T) {
	defer MockTargetStore().Clear()

	route := builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie()
	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(route).NewOrDie()

	assert.Equal(t, 0, len(mockTargetStore.mockRespTree))
	MockTargetStore().Set(mocktarget)
	assert.Equal(t, 1, len(mockTargetStore.mockRespTree))
	assert.True(t, MockTargetStore().Delete(mocktarget))
	assert.Equal(t, 0, len(mockTargetStore.mockRespTree))
}

func TestMockTargetClear(t *testing.T) {
	defer MockTargetStore().Clear()
	assert.Equal(t, 0, len(mockTargetStore.mockRespTree))

	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie(),
	).NewOrDie()
	MockTargetStore().Set(mocktarget)

	assert.Equal(t, 1, len(mockTargetStore.mockRespTree))
	MockTargetStore().Clear()
	assert.Equal(t, 0, len(mockTargetStore.mockRespTree))
}

func TestMockTargetIsEmpty(t *testing.T) {
	defer MockTargetStore().Clear()
	assert.True(t, MockTargetStore().IsEmpty())

	mocktarget := builder.NewMockTarget("foo", "bar").WithRoute(
		builder.NewRoute("/foo").WithStatusCode(200).WithMethods("GET").NewOrDie(),
	).NewOrDie()

	MockTargetStore().Set(mocktarget)
	assert.False(t, MockTargetStore().IsEmpty())
	MockTargetStore().Delete(mocktarget)
	assert.True(t, MockTargetStore().IsEmpty())
}
