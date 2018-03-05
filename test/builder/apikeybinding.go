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

// type PathBuilder struct {
// 	curr v2.Path
// }

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

func (b *KeyAccessBuilder) WithDefaultRule(rule v2.Rule) *KeyAccessBuilder {
	b.curr.DefaultRule = rule
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
