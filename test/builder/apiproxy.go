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

type ApiProxyBuilder struct {
	curr v2.ApiProxy
}

func NewApiProxy(name, namespace string) *ApiProxyBuilder {
	return &ApiProxyBuilder{
		curr: v2.ApiProxy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}

func (b *ApiProxyBuilder) WithSourcePath(path string) *ApiProxyBuilder {
	b.curr.Spec.Source.Path = path
	return b
}

func (b *ApiProxyBuilder) WithSourceHost(virtualHost string) *ApiProxyBuilder {
	b.curr.Spec.Source.VirtualHost = virtualHost
	return b
}

func (b *ApiProxyBuilder) WithTargetPath(path string) *ApiProxyBuilder {
	b.curr.Spec.Target.Path = path
	return b
}

func (b *ApiProxyBuilder) WithTargetBackendEndpoint(url string) *ApiProxyBuilder {
	b.curr.Spec.Target.Backend.Endpoint = &url
	return b
}

func (b *ApiProxyBuilder) WithPlugin(name, version string, config map[string]string) *ApiProxyBuilder {
  b.curr.Spec.Plugins = append(b.curr.Spec.Plugins, v2.Plugin{
    Name: name,
    Version: version,
    Config: config,
  })
  return b
}

func (b *ApiProxyBuilder) WithTargetBackendStaticService(name string, port int) *ApiProxyBuilder {
	b.curr.Spec.Target.Backend.Service = &v2.Service{
		Name: name,
		Port: int64(port),
	}
	return b
}

func (b *ApiProxyBuilder) WithSecret(name string) *ApiProxyBuilder {
	b.curr.Spec.Target.SSL = &v2.SSL{
		SecretName: name,
	}
	return b
}

func (b *ApiProxyBuilder) NewOrDie() *v2.ApiProxy {
	return &b.curr
}
