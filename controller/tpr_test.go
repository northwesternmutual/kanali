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
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
	e "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
)

func TestCreateTPRs(t *testing.T) {
	ctlr := Controller{
		RestClient: nil,
		ClientSet:  fake.NewSimpleClientset(),
		MasterHost: "foo.bar.com",
	}

	err := ctlr.CreateTPRs()
	assert.Nil(t, err)
	err = ctlr.CreateTPRs()
	assert.Nil(t, err)
	resource, _ := ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-proxy.kanali.io")
	assert.Equal(t, resource, &extensions.ThirdPartyResource{
		ObjectMeta: api.ObjectMeta{
			Name: "api-proxy.kanali.io",
		},
		Versions: []extensions.APIVersion{
			{
				Name: "v1",
			},
		},
		Description: "api proxy TPR",
	})
	resource, _ = ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-key.kanali.io")
	assert.Equal(t, resource, &extensions.ThirdPartyResource{
		ObjectMeta: api.ObjectMeta{
			Name: "api-key.kanali.io",
		},
		Versions: []extensions.APIVersion{
			{
				Name: "v1",
			},
		},
		Description: "api key TPR",
	})
	resource, _ = ctlr.ClientSet.Extensions().ThirdPartyResources().Get("api-key-binding.kanali.io")
	assert.Equal(t, resource, &extensions.ThirdPartyResource{
		ObjectMeta: api.ObjectMeta{
			Name: "api-key-binding.kanali.io",
		},
		Versions: []extensions.APIVersion{
			{
				Name: "v1",
			},
		},
		Description: "api key binding TPR",
	})
}

func TestIsKubernetesResourceAlreadyExistError(t *testing.T) {
	assert.False(t, isKubernetesResourceAlreadyExistError(errors.New("test error")))
	se := e.StatusError{
		ErrStatus: unversioned.Status{
			Code:   http.StatusConflict,
			Reason: unversioned.StatusReasonAlreadyExists,
		},
	}
	assert.True(t, isKubernetesResourceAlreadyExistError(&se))
	se.ErrStatus.Code = http.StatusNotFound
	assert.False(t, isKubernetesResourceAlreadyExistError(&se))
}
