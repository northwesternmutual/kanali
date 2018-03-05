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

func (b *ApiKeyBuilder) WithRevision(status v2.RevisionStatus, encryptedKey []byte) *ApiKeyBuilder {
	b.curr.Spec.Revisions = append(b.curr.Spec.Revisions, v2.Revision{
		Data:   string(encryptedKey),
		Status: status,
	})

	return b
}

func (b *ApiKeyBuilder) NewOrDie() *v2.ApiKey {
	return &b.curr
}
