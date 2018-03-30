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

// ApiProxiesGetter has a method to return a ApiProxyInterface.
// A group's client should implement this interface.
type ApiProxiesGetter interface {
	ApiProxies(namespace string) ApiProxyInterface
}

// ApiProxyInterface has methods to work with ApiProxy resources.
type ApiProxyInterface interface {
	Create(*v2.ApiProxy) (*v2.ApiProxy, error)
	Update(*v2.ApiProxy) (*v2.ApiProxy, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v2.ApiProxy, error)
	List(opts v1.ListOptions) (*v2.ApiProxyList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiProxy, err error)
	ApiProxyExpansion
}

// apiProxies implements ApiProxyInterface
type apiProxies struct {
	client rest.Interface
	ns     string
}

// newApiProxies returns a ApiProxies
func newApiProxies(c *KanaliV2Client, namespace string) *apiProxies {
	return &apiProxies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the apiProxy, and returns the corresponding apiProxy object, and an error if there is any.
func (c *apiProxies) Get(name string, options v1.GetOptions) (result *v2.ApiProxy, err error) {
	result = &v2.ApiProxy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apiproxies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApiProxies that match those selectors.
func (c *apiProxies) List(opts v1.ListOptions) (result *v2.ApiProxyList, err error) {
	result = &v2.ApiProxyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apiproxies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apiProxies.
func (c *apiProxies) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apiproxies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a apiProxy and creates it.  Returns the server's representation of the apiProxy, and an error, if there is any.
func (c *apiProxies) Create(apiProxy *v2.ApiProxy) (result *v2.ApiProxy, err error) {
	result = &v2.ApiProxy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apiproxies").
		Body(apiProxy).
		Do().
		Into(result)
	return
}

// Update takes the representation of a apiProxy and updates it. Returns the server's representation of the apiProxy, and an error, if there is any.
func (c *apiProxies) Update(apiProxy *v2.ApiProxy) (result *v2.ApiProxy, err error) {
	result = &v2.ApiProxy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apiproxies").
		Name(apiProxy.Name).
		Body(apiProxy).
		Do().
		Into(result)
	return
}

// Delete takes name of the apiProxy and deletes it. Returns an error if one occurs.
func (c *apiProxies) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apiproxies").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apiProxies) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apiproxies").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched apiProxy.
func (c *apiProxies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v2.ApiProxy, err error) {
	result = &v2.ApiProxy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apiproxies").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
