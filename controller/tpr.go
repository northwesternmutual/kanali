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
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/northwesternmutual/kanali/utils"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// tpr is an internal struct that stores
// information about a Kubernetes ThirdPartyResource
type tpr struct {
	Name        string
	Version     string
	Description string
}

// CreateTPRs will create the two tprs that kanali uses
// ApiProxy an ApiKey
func (c *Controller) CreateTPRs() error {

	logrus.Debug("creating TPRs")

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
		if !utils.IsKubernetesResourceAlreadyExistError(err) {
			return fmt.Errorf("Fail to create TPR: %v", err)
		}
	}

	return nil

}

// doCreateTPRs is a helper function that takes a
// list of tprs and adds each of them to our cluster
func (c *Controller) doCreateTPRs(tprs ...*tpr) error {

	for _, tpr := range tprs {
		if _, err := c.ClientSet.Extensions().ThirdPartyResources().Create(&extensions.ThirdPartyResource{
			ObjectMeta: api.ObjectMeta{
				Name: tpr.Name,
			},
			Versions: []extensions.APIVersion{
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
