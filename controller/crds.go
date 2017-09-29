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
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kanaliGroupName = "kanali.io"
)

var apiProxyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apiproxies.%s", kanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliGroupName,
		Version: "v1",
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apiproxies",
			Singular: "apiproxy",
			ShortNames: []string{
				"ap",
				"proxies",
			},
			Kind:     "ApiProxy",
			ListKind: "ApiProxyList",
		},
		Scope: apiextensionsv1beta1.NamespaceScoped,
	},
}

var apiKeyBindingCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeybindings.%s", kanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliGroupName,
		Version: "v1",
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apikeybindings",
			Singular: "apikeybinding",
			ShortNames: []string{
				"akb",
				"bindings",
			},
			Kind:     "ApiKeyBinding",
			ListKind: "ApiKeyBindingList",
		},
		Scope: apiextensionsv1beta1.NamespaceScoped,
	},
}

var apiKeyCRD = &apiextensionsv1beta1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("apikeys.%s", kanaliGroupName),
	},
	Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
		Group:   kanaliGroupName,
		Version: "v1",
		Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
			Plural:   "apikeys",
			Singular: "apikey",
			ShortNames: []string{
				"ak",
				"keys",
			},
			Kind:     "ApiKey",
			ListKind: "ApiKeyList",
		},
		Scope: apiextensionsv1beta1.ClusterScoped,
	},
}

// CreateCRDs will create all the CRDs that Kanali requires
func (c *Controller) CreateCRDs() error {
	return doCreateCRDs(c, apiProxyCRD, apiKeyBindingCRD, apiKeyCRD)
}

func doCreateCRDs(c *Controller, crds ...*apiextensionsv1beta1.CustomResourceDefinition) error {
	for _, crd := range crds {
		_, err := c.APIExtensionsV1beta1Interface.CustomResourceDefinitions().Create(crd)
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("Failed to create CRD: %v", err)
		}
	}

	return nil
}
