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

// ApiKeiesGetter has a method to return a ApiKeyInterface.
// A group's client should implement this interface.
type ApiKeiesGetter interface {
	ApiKeies() ApiKeyInterface
}

// ApiKeyInterface has methods to work with ApiKey resources.
type ApiKeyInterface interface {
	Create(*v2.ApiKey) (*v2.ApiKey, error)
	Update(*v2.ApiKey) (*v2.ApiKey, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.ApiKey, error)
	List(opts v1.ListOptions) (*v2.ApiKeyList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKey, err error)
	ApiKeyExpansion
}

// apiKeies implements ApiKeyInterface
type apiKeies struct {
	client rest.Interface
}

// newApiKeies returns a ApiKeies
func newApiKeies(c *KanaliV2Client) *apiKeies {
	return &apiKeies{
		client: c.RESTClient(),
	}
}

// Get takes name of the apiKey, and returns the corresponding apiKey object, and an error if there is any.
func (c *apiKeies) Get(name string, options v1.GetOptions) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Get().
		Resource("apikeies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApiKeies that match those selectors.
func (c *apiKeies) List(opts v1.ListOptions) (result *v2.ApiKeyList, err error) {
	result = &v2.ApiKeyList{}
	err = c.client.Get().
		Resource("apikeies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apiKeies.
func (c *apiKeies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("apikeies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a apiKey and creates it.  Returns the server's representation of the apiKey, and an error, if there is any.
func (c *apiKeies) Create(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Post().
		Resource("apikeies").
		Body(apiKey).
		Do().
		Into(result)
	return
}

// Update takes the representation of a apiKey and updates it. Returns the server's representation of the apiKey, and an error, if there is any.
func (c *apiKeies) Update(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Put().
		Resource("apikeies").
		Name(apiKey.Name).
		Body(apiKey).
		Do().
		Into(result)
	return
}

// Delete takes name of the apiKey and deletes it. Returns an error if one occurs.
func (c *apiKeies) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("apikeies").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apiKeies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("apikeies").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched apiKey.
func (c *apiKeies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Patch(pt).
		Resource("apikeies").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
