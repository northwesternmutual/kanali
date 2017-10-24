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
	v2 "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeApiKeyBindings implements ApiKeyBindingInterface
type FakeApiKeyBindings struct {
	Fake *FakeKanaliV2
	ns   string
}

var apikeybindingsResource = schema.GroupVersionResource{Group: "kanali.io", Version: "v2", Resource: "apikeybindings"}

var apikeybindingsKind = schema.GroupVersionKind{Group: "kanali.io", Version: "v2", Kind: "ApiKeyBinding"}

// Get takes name of the apiKeyBinding, and returns the corresponding apiKeyBinding object, and an error if there is any.
func (c *FakeApiKeyBindings) Get(name string, options v1.GetOptions) (result *v2.ApiKeyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(apikeybindingsResource, c.ns, name), &v2.ApiKeyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKeyBinding), err
}

// List takes label and field selectors, and returns the list of ApiKeyBindings that match those selectors.
func (c *FakeApiKeyBindings) List(opts v1.ListOptions) (result *v2.ApiKeyBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(apikeybindingsResource, apikeybindingsKind, c.ns, opts), &v2.ApiKeyBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.ApiKeyBindingList{}
	for _, item := range obj.(*v2.ApiKeyBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested apiKeyBindings.
func (c *FakeApiKeyBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(apikeybindingsResource, c.ns, opts))

}

// Create takes the representation of a apiKeyBinding and creates it.  Returns the server's representation of the apiKeyBinding, and an error, if there is any.
func (c *FakeApiKeyBindings) Create(apiKeyBinding *v2.ApiKeyBinding) (result *v2.ApiKeyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(apikeybindingsResource, c.ns, apiKeyBinding), &v2.ApiKeyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKeyBinding), err
}

// Update takes the representation of a apiKeyBinding and updates it. Returns the server's representation of the apiKeyBinding, and an error, if there is any.
func (c *FakeApiKeyBindings) Update(apiKeyBinding *v2.ApiKeyBinding) (result *v2.ApiKeyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(apikeybindingsResource, c.ns, apiKeyBinding), &v2.ApiKeyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKeyBinding), err
}

// Delete takes name of the apiKeyBinding and deletes it. Returns an error if one occurs.
func (c *FakeApiKeyBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(apikeybindingsResource, c.ns, name), &v2.ApiKeyBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApiKeyBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(apikeybindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v2.ApiKeyBindingList{})
	return err
}

// Patch applies the patch and returns the patched apiKeyBinding.
func (c *FakeApiKeyBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKeyBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(apikeybindingsResource, c.ns, name, data, subresources...), &v2.ApiKeyBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKeyBinding), err
}
