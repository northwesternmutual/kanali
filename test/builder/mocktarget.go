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
