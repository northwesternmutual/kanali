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
package v2

import (
	v2 "github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	scheme "github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ApiKeysGetter has a method to return a ApiKeyInterface.
// A group's client should implement this interface.
type ApiKeysGetter interface {
	ApiKeys() ApiKeyInterface
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

// apiKeys implements ApiKeyInterface
type apiKeys struct {
	client rest.Interface
}

// newApiKeys returns a ApiKeys
func newApiKeys(c *KanaliV2Client) *apiKeys {
	return &apiKeys{
		client: c.RESTClient(),
	}
}

// Get takes name of the apiKey, and returns the corresponding apiKey object, and an error if there is any.
func (c *apiKeys) Get(name string, options v1.GetOptions) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Get().
		Resource("apikeys").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApiKeys that match those selectors.
func (c *apiKeys) List(opts v1.ListOptions) (result *v2.ApiKeyList, err error) {
	result = &v2.ApiKeyList{}
	err = c.client.Get().
		Resource("apikeys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apiKeys.
func (c *apiKeys) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("apikeys").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a apiKey and creates it.  Returns the server's representation of the apiKey, and an error, if there is any.
func (c *apiKeys) Create(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Post().
		Resource("apikeys").
		Body(apiKey).
		Do().
		Into(result)
	return
}

// Update takes the representation of a apiKey and updates it. Returns the server's representation of the apiKey, and an error, if there is any.
func (c *apiKeys) Update(apiKey *v2.ApiKey) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Put().
		Resource("apikeys").
		Name(apiKey.Name).
		Body(apiKey).
		Do().
		Into(result)
	return
}

// Delete takes name of the apiKey and deletes it. Returns an error if one occurs.
func (c *apiKeys) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("apikeys").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apiKeys) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("apikeys").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched apiKey.
func (c *apiKeys) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiKey, err error) {
	result = &v2.ApiKey{}
	err = c.client.Patch(pt).
		Resource("apikeys").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
