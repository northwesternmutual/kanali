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
	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/crds"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/spf13/viper"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Controller is an exported struct
// that holds all of the information
// we will need throughout our program
// to talk to the kubernetes api server.
// this includes a rest client and well
// a clientset that gives us access to the libs
type Controller struct {
	RESTClient                    *rest.RESTClient
	ClientSet                     *kubernetes.Clientset
	APIExtensionsV1beta1Interface apiextensionsv1beta1.ApiextensionsV1beta1Interface
}

var (
	schemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	addToScheme        = schemeBuilder.AddToScheme
	schemeGroupVersion = schema.GroupVersion{Group: crds.KanaliGroupName, Version: "v1"}
)

// New creates a new kubernetes controller
// the controller is secure and uses the
// cluster's kubeconfig file to construct
// permissions
func New() (*Controller, error) {
	cfg, err := buildConfig(viper.GetString(config.FlagKubernetesKubeConfig.GetLong()))
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return nil, err
	}

	apiextensionsclientset, err := apiextensionsclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	config := *cfg
	config.APIPath = "/apis"
	config.GroupVersion = &schemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	restClient, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &Controller{restClient, clientset, apiextensionsclientset.ApiextensionsV1beta1()}, nil
}

func buildConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiProxy"), &spec.APIProxy{})
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiProxyList"), &spec.APIProxyList{})
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiKey"), &spec.APIKey{})
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiKeyList"), &spec.APIKeyList{})
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiKeyBinding"), &spec.APIKeyBinding{})
	scheme.AddKnownTypeWithName(schemeGroupVersion.WithKind("ApiKeyBindingList"), &spec.APIKeyBindingList{})
	metav1.AddToGroupVersion(scheme, schemeGroupVersion)
	return nil
}
