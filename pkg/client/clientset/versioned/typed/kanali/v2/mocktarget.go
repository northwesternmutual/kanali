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
package v2

import (
	v2 "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	scheme "github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MockTargetsGetter has a method to return a MockTargetInterface.
// A group's client should implement this interface.
type MockTargetsGetter interface {
	MockTargets(namespace string) MockTargetInterface
}

// MockTargetInterface has methods to work with MockTarget resources.
type MockTargetInterface interface {
	Create(*v2.MockTarget) (*v2.MockTarget, error)
	Update(*v2.MockTarget) (*v2.MockTarget, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.MockTarget, error)
	List(opts v1.ListOptions) (*v2.MockTargetList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.MockTarget, err error)
	MockTargetExpansion
}

// mockTargets implements MockTargetInterface
type mockTargets struct {
	client rest.Interface
	ns     string
}

// newMockTargets returns a MockTargets
func newMockTargets(c *KanaliV2Client, namespace string) *mockTargets {
	return &mockTargets{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the mockTarget, and returns the corresponding mockTarget object, and an error if there is any.
func (c *mockTargets) Get(name string, options v1.GetOptions) (result *v2.MockTarget, err error) {
	result = &v2.MockTarget{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mocktargets").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MockTargets that match those selectors.
func (c *mockTargets) List(opts v1.ListOptions) (result *v2.MockTargetList, err error) {
	result = &v2.MockTargetList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("mocktargets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested mockTargets.
func (c *mockTargets) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("mocktargets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a mockTarget and creates it.  Returns the server's representation of the mockTarget, and an error, if there is any.
func (c *mockTargets) Create(mockTarget *v2.MockTarget) (result *v2.MockTarget, err error) {
	result = &v2.MockTarget{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("mocktargets").
		Body(mockTarget).
		Do().
		Into(result)
	return
}

// Update takes the representation of a mockTarget and updates it. Returns the server's representation of the mockTarget, and an error, if there is any.
func (c *mockTargets) Update(mockTarget *v2.MockTarget) (result *v2.MockTarget, err error) {
	result = &v2.MockTarget{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("mocktargets").
		Name(mockTarget.Name).
		Body(mockTarget).
		Do().
		Into(result)
	return
}

// Delete takes name of the mockTarget and deletes it. Returns an error if one occurs.
func (c *mockTargets) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mocktargets").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *mockTargets) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("mocktargets").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched mockTarget.
func (c *mockTargets) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.MockTarget, err error) {
	result = &v2.MockTarget{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("mocktargets").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
