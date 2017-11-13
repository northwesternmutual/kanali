// // Copyright (c) 2017 Northwestern Mutual.
// //
// // Permission is hereby granted, free of charge, to any person obtaining a copy
// // of this software and associated documentation files (the "Software"), to deal
// // in the Software without restriction, including without limitation the rights
// // to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// // copies of the Software, and to permit persons to whom the Software is
// // furnished to do so, subject to the following conditions:
// //
// // The above copyright notice and this permission notice shall be included in
// // all copies or substantial portions of the Software.
// //
// // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// // IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// // FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// // AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// // LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// // OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// // THE SOFTWARE.

package controller

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/northwesternmutual/kanali/spec"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

type testHandlerFuncs struct {
	mutex        sync.RWMutex
	addResult    string
	updateResult string
	deleteResult string
}

func (f *testHandlerFuncs) addFunc(obj interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	foo, _ := obj.(string)
	f.addResult = fmt.Sprintf("modified %s", foo)
}

func (f *testHandlerFuncs) updateFunc(obj interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	foo, _ := obj.(string)
	f.updateResult = fmt.Sprintf("modified %s", foo)
}

func (f *testHandlerFuncs) deleteFunc(obj interface{}) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	foo, _ := obj.(string)
	f.deleteResult = fmt.Sprintf("modified %s", foo)
}

func (f *testHandlerFuncs) getAddResult() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.addResult
}

func (f *testHandlerFuncs) getUpdateResult() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.updateResult
}

func (f *testHandlerFuncs) getDeleteResult() string {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.deleteResult
}

func TestMonitor(t *testing.T) {
	eventCh := make(chan *event)
	testFuncs := &testHandlerFuncs{}
	go monitor(eventCh, testFuncs)

	addEvent := &event{
		Type:   added,
		Object: "test add",
	}
	updateEvent := &event{
		Type:   modified,
		Object: "test update",
	}
	deleteEvent := &event{
		Type:   deleted,
		Object: "test delete",
	}

	eventCh <- addEvent
	eventCh <- updateEvent
	eventCh <- deleteEvent

	time.Sleep(500 * time.Millisecond)

	assert.Equal(t, "modified test add", testFuncs.getAddResult())
	assert.Equal(t, "modified test update", testFuncs.getUpdateResult())
	assert.Equal(t, "modified test delete", testFuncs.getDeleteResult())

}

func getTestSecret() api.Secret {
	return api.Secret{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "secret-two",
			Namespace: "foo",
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.key": []byte("YWJjMTIz"),
			"tls.crt": []byte("ZGVmNDU2"),
		},
	}
}

func getTestService() api.Service {
	return api.Service{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
			Labels: map[string]string{
				"one":   "two",
				"three": "four",
			},
		},
		Spec: api.ServiceSpec{
			ClusterIP: "1.2.3.4",
		},
	}
}

func getTestAPIKey() spec.APIKey {
	return spec.APIKey{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "5e0329991f70e018f9503aa9c8f03ebe861df9be92b91dd5188baecd975572ab4d9973292768991200f433ef8830dfd914b134886ccb13b1cb31925028b44b10a03845718ee96db46cc738a294ff4ebd756d84ce3124e5f25d2818af7e325d16805977f658a6331e7b82db77a4366285f817df4394c45e6fb09ce9764e6813bf5ba449e2377f496bcdda07d2e27561b1c215582a1560f7b76ed5d5db29bd6d57d93e8068bb776bd7fd83a2e99319a8ff1ca27250b683a96af655566147302d75db31af3d6c0b27e9e1f1a471eea7f6cc31343b0caad14ad3ba8be8eabe3add5d9dd25594290325a1d4bdaa45c8502fd8edad015c93646aae513ac1e2cf22377b",
		},
	}
}

func getTestAPIKeyBinding() spec.APIKeyBinding {
	return spec.APIKeyBinding{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeyBindingSpec{
			APIProxyName: "api-proxy-one",
			Keys: []spec.Key{
				{
					Name:        "franks-api-key",
					DefaultRule: spec.Rule{},
					Subpaths: []*spec.Path{
						{
							Path: "/foo",
							Rule: spec.Rule{
								Global: true,
							},
						},
						{
							Path: "foo/bar",
							Rule: spec.Rule{
								Granular: &spec.GranularProxy{
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

func getTestAPIProxy() spec.APIProxy {
	return spec.APIProxy{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "api/v1/accounts",
			Hosts: []spec.Host{
				{
					Name: "foo.bar.com",
					SSL: spec.SSL{
						SecretName: "mySecretTwo",
					},
				},
			},
			Service: spec.Service{
				Name:      "my-service",
				Namespace: "foo",
				Port:      8080,
			},
			Plugins: []spec.Plugin{
				{
					Name:    "apikey",
					Version: "1.0.0",
				},
				{
					Name: "jwt",
				},
			},
			SSL: spec.SSL{
				SecretName: "mySecret",
			},
		},
	}
}
