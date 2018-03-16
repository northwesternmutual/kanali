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

package app

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/chain"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/client/informers/externalversions"
	"github.com/northwesternmutual/kanali/pkg/controller"
	"github.com/northwesternmutual/kanali/pkg/crds"
	v2CRDs "github.com/northwesternmutual/kanali/pkg/crds/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/log"
	_ "github.com/northwesternmutual/kanali/pkg/metrics"
	"github.com/northwesternmutual/kanali/pkg/middleware"
	"github.com/northwesternmutual/kanali/pkg/run"
	"github.com/northwesternmutual/kanali/pkg/server"
	//storev1 "github.com/northwesternmutual/kanali/pkg/store/core/v1"
	"github.com/northwesternmutual/kanali/pkg/client/informers/externalversions/internalinterfaces"
	"github.com/northwesternmutual/kanali/pkg/tracer"
	"github.com/northwesternmutual/kanali/pkg/utils"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var (
	resyncPeriod = 5 * time.Minute
)

func Run(sigCtx context.Context) error {
	logger := log.WithContext(sigCtx)

	decryptionKey, err := utils.LoadDecryptionKey(viper.GetString(options.FlagPluginsAPIKeyDecriptionKeyFile.GetLong()))
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}

	crdClientset, kanaliClientset, k8sClientset, err := createClientsets()
	if err != nil {
		logger.Fatal(err.Error())
		return err
	}

	// we need to create a tempory shared informer specific to the kubernetes clientset
	// so that we can merge it into one shared index informer later.
	coreV1SharedInformer := informers.NewSharedInformerFactory(k8sClientset, resyncPeriod).Core().V1()

	sharedInformer := externalversions.NewSharedInformerFactory(kanaliClientset, resyncPeriod)
	sharedInformer.InformerFor(&v1.Service{}, internalinterfaces.NewInformerFunc(func(i versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return coreV1SharedInformer.Services().Informer()
	}))
	sharedInformer.InformerFor(&v1.Secret{}, internalinterfaces.NewInformerFunc(func(i versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
		return coreV1SharedInformer.Secrets().Informer()
	}))

	controller.InitEventHandlers(sharedInformer, decryptionKey)

	// TODO: this is messy
	//storev1.SetGlobalInterface(k8sFactory.Core().V1())

	if err := crds.EnsureCRDs(
		crdClientset.ApiextensionsV1beta1(),
		[]*apiextensionsv1beta1.CustomResourceDefinition{
			v2CRDs.ApiProxyCRD,
			v2CRDs.ApiKeyCRD,
			v2CRDs.ApiKeyBindingCRD,
			v2CRDs.MockTargetCRD,
		}, nil,
	); err != nil {
		return err
	}

	tracer, tracerErr := tracer.New()
	if tracerErr != nil {
		logger.Warn(tracerErr.Error())
	}

	gatewayServer := server.Prepare(&server.Options{
		Name:         "gateway",
		InsecureAddr: viper.GetString(options.FlagServerInsecureBindAddress.GetLong()),
		SecureAddr:   viper.GetString(options.FlagServerSecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(options.FlagServerInsecurePort.GetLong()),
		SecurePort:   viper.GetInt(options.FlagServerSecurePort.GetLong()),
		TLSKey:       viper.GetString(options.FlagServerTLSKeyFile.GetLong()),
		TLSCert:      viper.GetString(options.FlagServerTLSCertFile.GetLong()),
		TLSCa:        viper.GetString(options.FlagServerTLSCaFile.GetLong()),
		Handler: chain.New().Add(
			middleware.Correlation,
			middleware.Metrics,
		).Link(middleware.Gateway),
		Logger: logger.Sugar(),
	})

	profilingServer := server.Prepare(&server.Options{
		Name:         "profiling",
		InsecureAddr: viper.GetString(options.FlagProfilingInsecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(options.FlagProfilingInsecurePort.GetLong()),
		Handler:      server.ProfilingHandler(),
		Logger:       logger.Sugar(),
	})

	metricsServer := server.Prepare(&server.Options{
		Name:         "prometheus",
		InsecureAddr: viper.GetString(options.FlagPrometheusServerInsecureBindAddress.GetLong()),
		InsecurePort: viper.GetInt(options.FlagPrometheusServerInsecurePort.GetLong()),
		Handler:      promhttp.Handler(),
		Logger:       logger.Sugar(),
	})

	ctx, cancel := context.WithCancel(sigCtx)

	var g run.Group
	g.Add(ctx, run.Always, "shared index informer", run.SharedInformerWrapper(sharedInformer))
	g.Add(ctx, tracerErr == nil, "tracer", tracer)
	g.Add(ctx, run.Always, metricsServer.Name(), metricsServer)
	g.Add(ctx, run.Always, gatewayServer.Name(), gatewayServer)
	g.Add(ctx, viper.GetBool(options.FlagProfilingEnabled.GetLong()), profilingServer.Name(), profilingServer)
	g.Add(ctx, run.Always, "parent process", run.MonitorContext(cancel))
	return g.Run()
}

func createClientsets() (
	crdClientset *clientset.Clientset,
	kanaliClientset *versioned.Clientset,
	k8sClientset *kubernetes.Clientset,
	err error,
) {
	config, err := utils.GetRestConfig(viper.GetString(options.FlagKubernetesKubeConfig.GetLong()))
	if err != nil {
		return nil, nil, nil, err
	}

	crdClientset, err = clientset.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	kanaliClientset, err = versioned.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}

	k8sClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, err
	}
	return
}
