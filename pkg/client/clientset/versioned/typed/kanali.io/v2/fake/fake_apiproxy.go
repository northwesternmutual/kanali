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

// FakeApiProxies implements ApiProxyInterface
type FakeApiProxies struct {
	Fake *FakeKanaliV2
	ns   string
}

var apiproxiesResource = schema.GroupVersionResource{Group: "kanali.io", Version: "v2", Resource: "apiproxies"}

var apiproxiesKind = schema.GroupVersionKind{Group: "kanali.io", Version: "v2", Kind: "ApiProxy"}

// Get takes name of the apiProxy, and returns the corresponding apiProxy object, and an error if there is any.
func (c *FakeApiProxies) Get(name string, options v1.GetOptions) (result *v2.ApiProxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(apiproxiesResource, c.ns, name), &v2.ApiProxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiProxy), err
}

// List takes label and field selectors, and returns the list of ApiProxies that match those selectors.
func (c *FakeApiProxies) List(opts v1.ListOptions) (result *v2.ApiProxyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(apiproxiesResource, apiproxiesKind, c.ns, opts), &v2.ApiProxyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2.ApiProxyList{}
	for _, item := range obj.(*v2.ApiProxyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested apiProxies.
func (c *FakeApiProxies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(apiproxiesResource, c.ns, opts))

}

// Create takes the representation of a apiProxy and creates it.  Returns the server's representation of the apiProxy, and an error, if there is any.
func (c *FakeApiProxies) Create(apiProxy *v2.ApiProxy) (result *v2.ApiProxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(apiproxiesResource, c.ns, apiProxy), &v2.ApiProxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiProxy), err
}

// Update takes the representation of a apiProxy and updates it. Returns the server's representation of the apiProxy, and an error, if there is any.
func (c *FakeApiProxies) Update(apiProxy *v2.ApiProxy) (result *v2.ApiProxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(apiproxiesResource, c.ns, apiProxy), &v2.ApiProxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiProxy), err
}

// Delete takes name of the apiProxy and deletes it. Returns an error if one occurs.
func (c *FakeApiProxies) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(apiproxiesResource, c.ns, name), &v2.ApiProxy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApiProxies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(apiproxiesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v2.ApiProxyList{})
	return err
}

// Patch applies the patch and returns the patched apiProxy.
func (c *FakeApiProxies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiProxy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(apiproxiesResource, c.ns, name, data, subresources...), &v2.ApiProxy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v2.ApiProxy), err
}
