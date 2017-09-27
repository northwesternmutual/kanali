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

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
  apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type crd struct {
	Name        string
  Group string
	Version     string
  Kind string
	Description string
}

// CreateTPRs will create the two tprs that kanali will use
func (c *Controller) CreateCRDs() error {
	if err := c.doCreateCRDs(&crd{
		Name:        "api-proxy.kanali.io",
		Version:     "v1",
		Description: "api proxy TPR",
	}, &tpr{
		Name:        "api-key.kanali.io",
		Version:     "v1",
		Description: "api key TPR",
	}, &tpr{
		Name:        "api-key-binding.kanali.io",
		Version:     "v1",
		Description: "api key binding TPR",
	}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return fmt.Errorf("Failed to create TPR: %v", err)
		}
	}

	return nil
}

func (c *Controller) doCreateCRDs(crds ...*crd) error {
	for _, crd := range crds {
		if _, err := c.APIExtensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(&apiextensionsv1beta1.CustomResourceDefinition{
  		ObjectMeta: metav1.ObjectMeta{
  			Name: crd.Name,
  		},
  		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
  			Group:   strings.Split(crd.Name, "."),
  			Version: crv1.SchemeGroupVersion.Version,
  			Scope:   apiextensionsv1beta1.NamespaceScoped,
  			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
  				Plural: crv1.ExampleResourcePlural,
  				Kind:   reflect.TypeOf(crv1.Example{}).Name(),
  			},
  		},
  	}); err != nil {
			return err
		}
	}

	return nil
}
