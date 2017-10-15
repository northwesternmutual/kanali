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

package crds

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

const (
	// KanaliGroupName represents the CRD group
	KanaliGroupName = "kanali.io"
  Version = "v2alpha1"
)

// CreateCRDs will create all the CRDs that Kanali requires
func CreateCRDs(i apiextensionsv1beta1client.ApiextensionsV1beta1Interface) error {
	return doCreateCRDs(i, apiProxyCRD, apiKeyBindingCRD, apiKeyCRD)
}

func doCreateCRDs(i apiextensionsv1beta1client.ApiextensionsV1beta1Interface, crds ...*apiextensionsv1beta1.CustomResourceDefinition) error {
	for _, crd := range crds {
		_, err := i.CustomResourceDefinitions().Create(crd)
		if err != nil && !errors.IsAlreadyExists(err) {
			return fmt.Errorf("Failed to create CRD: %v", err)
		}
	}

	return nil
}
