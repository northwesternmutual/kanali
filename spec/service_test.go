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
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func TestGetServiceStore(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore

	store.Clear()
	assert.Equal(0, len(store.serviceMap), "store should be empty")

	v := ServiceFactory{}
	var i interface{} = &v
	_, ok := i.(Store)
	assert.True(ok, "ServiceFactory does not implement the Store interface")
}

func TestIsEmpty(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Clear()
	assert.True(store.IsEmpty())
	store.Set(serviceList[0])
	assert.False(store.IsEmpty())
	store.Clear()
	assert.True(store.IsEmpty())

	store.Set(serviceList[0])
	assert.False(store.IsEmpty())
	store.Delete(serviceList[0])
	assert.True(store.IsEmpty())
}

func TestCreateService(t *testing.T) {
	assert := assert.New(t)
	message := "service received is not expected"

	svc := CreateService(api.Service{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
			Labels: map[string]string{
				"one":   "two",
				"three": "four",
			},
		},
		Spec: api.ServiceSpec{},
	})

	assert.Equal("foo", svc.Name, message)
	assert.Equal("bar", svc.Namespace, message)
	assert.Equal(2, len(svc.Labels), message)
	assert.True((svc.Labels[0] == Label{
		Name:  "one",
		Value: "two",
	}) || (svc.Labels[1] == Label{
		Name:  "one",
		Value: "two",
	}), message)
	assert.True((svc.Labels[0] == Label{
		Name:  "three",
		Value: "four",
	}) || (svc.Labels[1] == Label{
		Name:  "three",
		Value: "four",
	}), message)

}

func TestServiceSet(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Clear()
	store.Set(serviceList[0])
	store.Set(serviceList[1])
	err := store.Set(APIProxy{})
	assert.Equal(err.Error(), "grrr - you're only allowed add services to the services store.... duh", "wrong error")
	assert.Equal(2, len(store.serviceMap), "store should have 2 namespaces represented")
	assert.Equal(serviceList[0], store.serviceMap["foo"][0], "service should be present")
	assert.Equal(serviceList[1], store.serviceMap["bar"][0], "service should be present")
	svcOne := serviceList[0]
	svcTwo := serviceList[1]
	svcOne.Labels[1].Name = "name-foo"
	svcTwo.Labels[1].Name = "name-foo"
	store.Set(serviceList[0])
	store.Set(serviceList[1])
	assert.Equal(svcOne, store.serviceMap["foo"][0], "service should be present")
	assert.Equal(svcTwo, store.serviceMap["bar"][0], "service should be present")
	store.Set(serviceList[2])
	assert.Equal(serviceList[2], store.serviceMap["bar"][1], "service should be present")
}

func TestServiceUpdate(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Clear()
	store.Update(serviceList[0])
	store.Update(serviceList[1])
	err := store.Update(APIProxy{})
	assert.Equal(err.Error(), "grrr - you're only allowed add services to the services store.... duh", "wrong error")
	assert.Equal(2, len(store.serviceMap), "store should have 2 namespaces represented")
	assert.Equal(serviceList[0], store.serviceMap["foo"][0], "service should be present")
	assert.Equal(serviceList[1], store.serviceMap["bar"][0], "service should be present")
	svcOne := serviceList[0]
	svcTwo := serviceList[1]
	svcOne.Labels[1].Name = "name-foo"
	svcTwo.Labels[1].Name = "name-foo"
	store.Update(serviceList[0])
	store.Update(serviceList[1])
	assert.Equal(svcOne, store.serviceMap["foo"][0], "service should be present")
	assert.Equal(svcTwo, store.serviceMap["bar"][0], "service should be present")
	store.Update(serviceList[2])
	assert.Equal(serviceList[2], store.serviceMap["bar"][1], "service should be present")
}

func TestServiceClear(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Set(serviceList[0])
	store.Set(serviceList[1])
	store.Set(serviceList[2])
	store.Clear()
	assert.Equal(0, len(store.serviceMap), "store should be empty")
}

func TestServiceGet(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Clear()
	store.Set(serviceList[0])
	store.Set(serviceList[1])
	_, err := store.Get("")
	assert.Equal(err.Error(), "getting a service requires 2 parameters", "wrong error")
	_, err = store.Get(APIProxy{}, "")
	assert.Equal(err.Error(), "first argument should be a service", "wrong error")
	_, err = store.Get(Service{}, 5)
	assert.Equal(err.Error(), "second argument should either be nil or http.Header", "wrong error")

	headerOne := http.Header{}
	headerOne.Add("x-nm-deploy", "production")
	headerTwo := http.Header{}
	headerTwo.Add("X-nM-dePloy", "pRoduCtion")

	viper.SetDefault("headers.x-nm-deploy", "production")

	result, _ := store.Get(Service{}, nil)
	assert.Nil(result, "nil service should return nil")
	result, _ = store.Get(Service{
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
		},
	}, headerOne)
	assert.Equal(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, headerOne)
	assert.Equal(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.Equal(serviceList[0], result, "service should exist")
	viper.SetDefault("headers.x-nm-deploy", "")
	result, _ = store.Get(Service{
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.Nil(result, "service should not exist")
	result, _ = store.Get(Service{
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, headerOne)
	assert.Equal(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "bar",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-three",
			},
		},
	}, headerOne)
	assert.NotEqual(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Name:      "my-cool-service-one",
		Namespace: "foo",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-three",
			},
		},
	}, headerOne)
	assert.Equal(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "bar",
		Labels: Labels{
			Label{
				Name:   "release",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-two",
				Value: "value-three",
			},
		},
	}, headerTwo)
	assert.NotEqual(serviceList[0], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "bar",
		Labels: Labels{
			Label{
				Name:   "name-prime-one",
				Header: "x-nm-deploy",
			},
			Label{
				Name:  "name-prime-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.NotEqual(serviceList[1], result, "service should exist")
	result, _ = store.Get(Service{
		Namespace: "bar",
		Labels: Labels{
			Label{
				Name:  "name-prime-one",
				Value: "value-one",
			},
			Label{
				Name:  "name-prime-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.Equal(serviceList[1], result, "service should exist")
}

func TestServiceDelete(t *testing.T) {
	assert := assert.New(t)
	store := ServiceStore
	serviceList := getTestServiceList()

	store.Clear()
	store.Set(serviceList[0])
	store.Set(serviceList[1])
	store.Set(serviceList[2])
	store.Set(serviceList[0])
	store.Set(serviceList[1])
	result, _ := store.Delete(nil)
	assert.Nil(result, "nil service should return nil")
	result, _ = store.Delete(Service{
		Name:      "my-cool-service-one",
		Namespace: "foo",
		Labels:    Labels{},
	})
	assert.Equal(serviceList[0], result, "delete should return deleted service")
	result, _ = store.Delete(Service{
		Name:      "my-cool-service-one",
		Namespace: "foo",
		Labels:    Labels{},
	})
	assert.Nil(result, "service should not exist")
	assert.Equal(0, len(store.serviceMap["foo"]), "length not equal")
	assert.Equal(1, len(store.serviceMap), "length not equal")
	assert.Equal(2, len(store.serviceMap["bar"]), "length not equal")
	result, _ = store.Delete(Service{
		Name:      "my-cool-service-two",
		Namespace: "bar",
		Labels:    Labels{},
	})
	assert.Equal(serviceList[1], result, "delete should return deleted service")
	assert.Equal(1, len(store.serviceMap), "store should be empty")
	result, _ = store.Delete(Service{
		Name:      "my-cool-service-two",
		Namespace: "bar",
		Labels:    Labels{},
	})
	assert.Nil(result, "deleted service should return nil")
}

func getTestServiceList() services {

	return services{
		Service{
			Name:      "my-cool-service-one",
			Namespace: "foo",
			Labels: Labels{
				Label{
					Name:  "release",
					Value: "production",
				},
				Label{
					Name:  "name-two",
					Value: "value-two",
				},
			},
		},
		Service{
			Name:      "my-cool-service-two",
			Namespace: "bar",
			Labels: Labels{
				Label{
					Name:  "name-prime-one",
					Value: "value-one",
				},
				Label{
					Name:  "name-prime-two",
					Value: "value-two",
				},
			},
		},
		Service{
			Name:      "my-cool-service-three",
			Namespace: "bar",
			Labels: Labels{
				Label{
					Name:  "name-prime-one",
					Value: "value-one",
				},
				Label{
					Name:  "name-prime-two",
					Value: "value-two",
				},
			},
		},
	}

}
