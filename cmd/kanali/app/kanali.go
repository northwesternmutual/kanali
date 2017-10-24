package app

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/northwesternmutual/kanali/cmd/kanali/app/options"
	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/client/informers/externalversions"
	apikey "github.com/northwesternmutual/kanali/pkg/controller/apikey"
	apikeybinding "github.com/northwesternmutual/kanali/pkg/controller/apikeybinding"
	apiproxy "github.com/northwesternmutual/kanali/pkg/controller/apiproxy"
	mocktarget "github.com/northwesternmutual/kanali/pkg/controller/mocktarget"
	logging "github.com/northwesternmutual/kanali/pkg/logging"
	traffic "github.com/northwesternmutual/kanali/pkg/traffic"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func Run(ctx context.Context) error {

	logger := logging.WithContext(ctx)

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	crdClientset, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return err
	}

	if err := createCRDs(crdClientset.ApiextensionsV1beta1()); err != nil {
		logger.Fatal(err.Error())
		return err
	}

	kanaliClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return err
	}

	k8sClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	decryptionKey, err := loadDecryptionKey(viper.GetString(options.FlagPluginsAPIKeyDecriptionKeyFile.GetLong()))
	if err != nil {
		return err
	}

	kanaliFactory := externalversions.NewSharedInformerFactory(kanaliClientset, 5*time.Minute)
	k8sFactory := informers.NewSharedInformerFactory(k8sClientset, 5*time.Minute)

	go apikey.NewApiKeyController(kanaliFactory.Kanali().V2().ApiKeies(), decryptionKey).Run(ctx.Done())
	go apikeybinding.NewApiKeyBindingController(kanaliFactory.Kanali().V2().ApiKeyBindings()).Run(ctx.Done())
	go apiproxy.NewApiProxyController(kanaliFactory.Kanali().V2().ApiProxies()).Run(ctx.Done())
	go mocktarget.NewMockTargetController(kanaliFactory.Kanali().V2().MockTargets()).Run(ctx.Done())
	go k8sFactory.Core().V1().Services().Informer().Run(ctx.Done())
	go k8sFactory.Core().V1().Secrets().Informer().Run(ctx.Done())

	// TODO: handle case that ctx.Done() stop channel sends an item through

	trafficCtlr, err := traffic.NewTrafficController()
	if err != nil {
		return err
	}
	defer func() {
		if err := trafficCtlr.Client.Close(); err != nil {
			logger.Warn(err.Error())
		}
	}()

	go trafficCtlr.MonitorTraffic(ctx)

	tracer, closer, err := newJaegerTracer()
	if err != nil {
		logger.Warn(err.Error())
	} else {
		opentracing.SetGlobalTracer(tracer)
		defer func() {
			if err := closer.Close(); err != nil {
				logger.Warn(err.Error())
			}
		}()
	}

	influxCtlr, err := newInfluxdbController()
	if err != nil {
		logger.Warn(err.Error())
	} else {
		defer func() {
			if err := influxCtlr.client.Close(); err != nil {
				logger.Warn(err.Error())
			}
		}()
	}

	// will always returns a non-nil error
	return startHTTP(ctx, getHTTPHandler(influxCtlr, k8sFactory.Core()))

}

func loadDecryptionKey(location string) (*rsa.PrivateKey, error) {
	// read in private key
	keyBytes, err := ioutil.ReadFile(location)
	if err != nil {
		return nil, err
	}
	// create a pem block from the private key provided
	block, _ := pem.Decode(keyBytes)
	// parse the pem block into a private key
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
