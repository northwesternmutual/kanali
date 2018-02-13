package builder

import (
	"strconv"
	"strings"

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
	// http://ipaddr:port
	one := strings.Split(url, "://")
	two := strings.Split(one[1], ":")
	port, _ := strconv.Atoi(two[1])

	b.curr.Spec.Target.Backend.Endpoint = &v2.Endpoint{
		Scheme: one[0],
		Host:   two[0],
		Port:   int64(port),
	}
	return b
}

func (b *ApiProxyBuilder) NewOrDie() *v2.ApiProxy {
	return &b.curr
}
