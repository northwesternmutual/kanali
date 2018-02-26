package builder

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/rsa"
)

type ApiKeyBuilder struct {
	curr v2.ApiKey
}

func NewApiKey(name string) *ApiKeyBuilder {
	return &ApiKeyBuilder{
		curr: v2.ApiKey{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: v2.ApiKeySpec{
				Revisions: []v2.Revision{},
			},
		},
	}
}

func (b *ApiKeyBuilder) WithRevision(status RevisionStatus, unencryptedKey string) *ApiKeyBuilder {
	encryptedKey, _ := rsa.Encrypt(unencryptedKey, nil, rsa.Base64Encode())

	b.curr.Spec.Revisions = append(b.curr.Spec.Revisions, v2.Revision{
		Data:   encryptedKey,
		Status: status,
	})

	return b
}

func (b *ApiKeyBuilder) NewOrDie() *v2.ApiKey {
	return &b.curr
}
