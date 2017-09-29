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

package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateTPRs(t *testing.T) {
	ctlr := Controller{
		RESTClient:                    nil,
		ClientSet:                     nil,
		APIExtensionsV1beta1Interface: fake.NewSimpleClientset().ApiextensionsV1beta1(),
	}

	err := ctlr.CreateCRDs()
	assert.Nil(t, err)

	proxyCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apiproxies.kanali.io", v1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, proxyCRD.ObjectMeta.Name, "apiproxies.kanali.io")

	bindingCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apikeybindings.kanali.io", v1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, bindingCRD.ObjectMeta.Name, "apikeybindings.kanali.io")

	keyCRD, err := ctlr.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Get("apikeys.kanali.io", v1.GetOptions{})
	assert.Nil(t, err)
	assert.Equal(t, keyCRD.ObjectMeta.Name, "apikeys.kanali.io")
}
