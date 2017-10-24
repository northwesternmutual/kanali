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

// ApiKeyBindingsGetter has a method to return a ApiKeyBindingInterface.
// A group's client should implement this interface.
type ApiKeyBindingsGetter interface {
	ApiKeyBindings(namespace string) ApiKeyBindingInterface
}

// ApiKeyBindingInterface has methods to work with ApiKeyBinding resources.
type ApiKeyBindingInterface interface {
	Create(*v2.ApiKeyBinding) (*v2.ApiKeyBinding, error)
	Update(*v2.ApiKeyBinding) (*v2.ApiKeyBinding, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.ApiKeyBinding, error)
	List(opts v1.ListOptions) (*v2.ApiKeyBindingList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKeyBinding, err error)
	ApiKeyBindingExpansion
}

// apiKeyBindings implements ApiKeyBindingInterface
type apiKeyBindings struct {
	client rest.Interface
	ns     string
}

// newApiKeyBindings returns a ApiKeyBindings
func newApiKeyBindings(c *KanaliV2Client, namespace string) *apiKeyBindings {
	return &apiKeyBindings{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the apiKeyBinding, and returns the corresponding apiKeyBinding object, and an error if there is any.
func (c *apiKeyBindings) Get(name string, options v1.GetOptions) (result *v2.ApiKeyBinding, err error) {
	result = &v2.ApiKeyBinding{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apikeybindings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApiKeyBindings that match those selectors.
func (c *apiKeyBindings) List(opts v1.ListOptions) (result *v2.ApiKeyBindingList, err error) {
	result = &v2.ApiKeyBindingList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apikeybindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apiKeyBindings.
func (c *apiKeyBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apikeybindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a apiKeyBinding and creates it.  Returns the server's representation of the apiKeyBinding, and an error, if there is any.
func (c *apiKeyBindings) Create(apiKeyBinding *v2.ApiKeyBinding) (result *v2.ApiKeyBinding, err error) {
	result = &v2.ApiKeyBinding{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apikeybindings").
		Body(apiKeyBinding).
		Do().
		Into(result)
	return
}

// Update takes the representation of a apiKeyBinding and updates it. Returns the server's representation of the apiKeyBinding, and an error, if there is any.
func (c *apiKeyBindings) Update(apiKeyBinding *v2.ApiKeyBinding) (result *v2.ApiKeyBinding, err error) {
	result = &v2.ApiKeyBinding{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apikeybindings").
		Name(apiKeyBinding.Name).
		Body(apiKeyBinding).
		Do().
		Into(result)
	return
}

// Delete takes name of the apiKeyBinding and deletes it. Returns an error if one occurs.
func (c *apiKeyBindings) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apikeybindings").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apiKeyBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apikeybindings").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched apiKeyBinding.
func (c *apiKeyBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKeyBinding, err error) {
	result = &v2.ApiKeyBinding{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apikeybindings").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
