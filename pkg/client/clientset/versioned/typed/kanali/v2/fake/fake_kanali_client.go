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
package fake

import (
	v2 "github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/typed/kanali/v2"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeKanaliV2 struct {
	*testing.Fake
}

func (c *FakeKanaliV2) ApiKeies() v2.ApiKeyInterface {
	return &FakeApiKeies{c}
}

func (c *FakeKanaliV2) ApiKeyBindings(namespace string) v2.ApiKeyBindingInterface {
	return &FakeApiKeyBindings{c, namespace}
}

func (c *FakeKanaliV2) ApiProxies(namespace string) v2.ApiProxyInterface {
	return &FakeApiProxies{c, namespace}
}

func (c *FakeKanaliV2) MockTargets(namespace string) v2.MockTargetInterface {
	return &FakeMockTargets{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeKanaliV2) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
