package builder

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type secretBuilder struct {
	curr v1.Secret
}

func NewSecretBuilder(name, namespace string) *secretBuilder {
	return &secretBuilder{
		curr: v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:        name,
				Namespace:   namespace,
				Annotations: map[string]string{},
			},
			Data: map[string][]byte{},
		},
	}
}

func (b *secretBuilder) WithAnnotation(key, val string) *secretBuilder {
	b.curr.ObjectMeta.Annotations[key] = val
	return b
}

func (b *secretBuilder) WithData(key string, val []byte) *secretBuilder {
	b.curr.Data[key] = val
	return b
}

func (b *secretBuilder) NewOrDie() *v1.Secret {
	return &b.curr
}
