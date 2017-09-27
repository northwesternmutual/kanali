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
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type tpr struct {
	Name        string
	Version     string
	Description string
}

// CreateTPRs will create the two tprs that kanali will use
func (c *Controller) CreateTPRs() error {
	if err := c.doCreateTPRs(&tpr{
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

func (c *Controller) doCreateTPRs(tprs ...*tpr) error {
	for _, tpr := range tprs {
		if _, err := c.ClientSet.ExtensionsV1beta1Client.ThirdPartyResources().Create(&v1beta1.ThirdPartyResource{
			ObjectMeta: metav1.ObjectMeta{
				Name: tpr.Name,
			},
			Versions: []v1beta1.APIVersion{
				{
					Name: tpr.Version,
				},
			},
			Description: tpr.Description,
		}); err != nil {
			return err
		}
	}

	return nil
}

func isKubernetesResourceAlreadyExistError(err error) bool {
	se, ok := err.(*errors.StatusError)
	if !ok {
		return false
	} else if se.Status().Code == http.StatusConflict && se.Status().Reason == unversioned.StatusReasonAlreadyExists {
		return true
	} else {
		return false
	}
}
