// Copyright (c) 2018 Northwestern Mutual.
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

package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
)

type MockTargetBuilder struct {
	curr v2.MockTarget
}

type RouteBuilder struct {
	curr v2.Route
}

func NewMockTarget(name, namespace string) *MockTargetBuilder {
	return &MockTargetBuilder{
		curr: v2.MockTarget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v2.MockTargetSpec{
				Routes: []v2.Route{},
			},
		},
	}
}

func (b *MockTargetBuilder) WithRoute(r v2.Route) *MockTargetBuilder {
	b.curr.Spec.Routes = append(b.curr.Spec.Routes, r)
	return b
}

func NewRoute(path string) *RouteBuilder {
	return &RouteBuilder{
		curr: v2.Route{
			Path: path,
		},
	}
}

func (b *RouteBuilder) WithStatusCode(code int) *RouteBuilder {
	b.curr.StatusCode = code
	return b
}

func (b *RouteBuilder) WithMethods(methods ...string) *RouteBuilder {
	b.curr.Methods = methods
	return b
}

func (b *RouteBuilder) WithBody(data []byte) *RouteBuilder {
	b.curr.Body = data
	return b
}

func (b *RouteBuilder) WithHeaders(headers map[string]string) *RouteBuilder {
	b.curr.Headers = headers
	return b
}

func (b *RouteBuilder) NewOrDie() *v2.Route {
	return &b.curr
}

func (b *MockTargetBuilder) NewOrDie() *v2.MockTarget {
	return &b.curr
}
