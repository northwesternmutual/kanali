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

	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// Controller is an exported struct
// that holds all of the information
// we will need throughout our program
// to talk to the kubernetes api server.
// this includes a rest client and well
// a clientset that gives us access to the libs
type Controller struct {
	RestClient *restclient.RESTClient
	ClientSet  internalclientset.Interface
	MasterHost string
}

// New creates a new kubernetes controller
// the controller is secure and uses the
// cluster's kubeconfig file to construct
// permissions
func New() (*Controller, error) {

	f := util.NewFactory(nil)

	// ClientSet gives you back an internal, generated clientset
	clientSet, err := f.ClientSet()
	if err != nil {
		return nil, err
	}

	// Returns a RESTClient for accessing Kubernetes resources or an error.
	restClient, err := f.RESTClient()
	if err != nil {
		return nil, err
	}

	controller := &Controller{
		RestClient: restClient,
		ClientSet:  clientSet,
	}

	// Returns a client.Config for accessing the Kubernetes server.
	k8sConfig, err := f.ClientConfig()
	if err != nil {
		return nil, err
	} else if k8sConfig == nil {
		return nil, errors.New("received nil k8sConfig, please check if k8s cluster is available")
	} else {
		controller.MasterHost = k8sConfig.Host
	}

	// return our newly created controller
	return controller, nil
}
