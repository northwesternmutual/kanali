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

// FakeApiKeies implements ApiKeyInterface
type FakeApiKeies struct {
	Fake *FakeKanaliV2
}

var apikeiesResource = schema.GroupVersionResource{Group: "kanali.io", Version: "v2", Resource: "apikeies"}

var apikeiesKind = schema.GroupVersionKind{Group: "kanali.io", Version: "v2", Kind: "ApiKey"}

// Get takes name of the apiKey, and returns the corresponding apiKey object, and an error if there is any.
func (c *FakeApiKeies) Get(name string, options v1.GetOptions) (result *v2.ApiKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(apikeiesResource, name), &v2.ApiKey{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKey), err
}

// List takes label and field selectors, and returns the list of ApiKeies that match those selectors.
func (c *FakeApiKeies) List(opts v1.ListOptions) (result *v2.ApiKeyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(apikeiesResource, apikeiesKind, opts), &v2.ApiKeyList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.ApiKeyList{}
	for _, item := range obj.(*v2.ApiKeyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested apiKeies.
func (c *FakeApiKeies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(apikeiesResource, opts))
}

// Create takes the representation of a apiKey and creates it.  Returns the server's representation of the apiKey, and an error, if there is any.
func (c *FakeApiKeies) Create(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(apikeiesResource, apiKey), &v2.ApiKey{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKey), err
}

// Update takes the representation of a apiKey and updates it. Returns the server's representation of the apiKey, and an error, if there is any.
func (c *FakeApiKeies) Update(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(apikeiesResource, apiKey), &v2.ApiKey{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKey), err
}

// Delete takes name of the apiKey and deletes it. Returns an error if one occurs.
func (c *FakeApiKeies) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(apikeiesResource, name), &v2.ApiKey{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApiKeies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(apikeiesResource, listOptions)

	_, err := c.Fake.Invokes(action, &v2.ApiKeyList{})
	return err
}

// Patch applies the patch and returns the patched apiKey.
func (c *FakeApiKeies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKey, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(apikeiesResource, name, data, subresources...), &v2.ApiKey{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiKey), err
}
