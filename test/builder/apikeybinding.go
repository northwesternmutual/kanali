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

type ApiKeyBindingBuilder struct {
	curr v2.ApiKeyBinding
}

type KeyAccessBuilder struct {
	curr v2.Key
}

type RuleBuilder struct {
	curr v2.Rule
}

type PathBuilder struct {
	curr v2.Path
}

func NewApiKeyBinding(name, namespace string) *ApiKeyBindingBuilder {
	return &ApiKeyBindingBuilder{
		curr: v2.ApiKeyBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: v2.ApiKeyBindingSpec{
				Keys: []v2.Key{},
			},
		},
	}
}

func (b *ApiKeyBindingBuilder) WithKeys(keys ...v2.Key) *ApiKeyBindingBuilder {
	b.curr.Spec.Keys = append(b.curr.Spec.Keys, keys...)
	return b
}

func (b *ApiKeyBindingBuilder) NewOrDie() *v2.ApiKeyBinding {
	return &b.curr
}

func NewKeyAccess(keyName string) *KeyAccessBuilder {
	return &KeyAccessBuilder{
		curr: v2.Key{
			Name: keyName,
		},
	}
}

func NewPathBuilder(path string, rule v2.Rule) *PathBuilder {
	return &PathBuilder{
		curr: v2.Path{
			Path: path,
			Rule: rule,
		},
	}
}

func (b *PathBuilder) NewOrDie() v2.Path {
	return b.curr
}

func (b *KeyAccessBuilder) WithDefaultRule(rule v2.Rule) *KeyAccessBuilder {
	b.curr.DefaultRule = rule
	return b
}

func (b *KeyAccessBuilder) WithSubpaths(paths ...v2.Path) *KeyAccessBuilder {
	b.curr.Subpaths = paths
	return b
}

func (b *KeyAccessBuilder) NewOrDie() v2.Key {
	return b.curr
}

func NewRule() *RuleBuilder {
	return &RuleBuilder{}
}

func (b *RuleBuilder) WithGlobal() *RuleBuilder {
	b.curr.Global = true
	return b
}

func (b *RuleBuilder) WithGranular(verbs ...string) *RuleBuilder {
	b.curr.Granular.Verbs = verbs
	return b
}

func (b *RuleBuilder) NewOrDie() v2.Rule {
	return b.curr
}
