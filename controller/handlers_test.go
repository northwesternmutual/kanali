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
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"testing"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/logging"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func TestAPIProxyHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	testProxyOne := spec.APIProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name:      "my-service",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	apiProxyHandlerFuncs.AddFunc(&testProxyOne)
	result, _ := spec.ProxyStore.Get("/api/v1/accounts")
	assert.Equal(t, result, testProxyOne)
	assertLog(t, zapcore.DebugLevel, "added ApiProxy testProxyOne in foo namespace", obsvr)

	apiProxyHandlerFuncs.AddFunc(testProxyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed APIProxy from k8s apiserver", obsvr)

	apiProxyHandlerFuncs.UpdateFunc(&testProxyOne, &testProxyOne)
	assertLog(t, zapcore.DebugLevel, "updated ApiProxy testProxyOne in foo namespace", obsvr)

	apiProxyHandlerFuncs.UpdateFunc(testProxyOne, &testProxyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiProxy from k8s apiserver", obsvr)

	apiProxyHandlerFuncs.UpdateFunc(&testProxyOne, testProxyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiProxy from k8s apiserver", obsvr)

	testProxyTwo := testProxyOne
	testProxyTwo.ObjectMeta.Namespace = "bar"
	apiProxyHandlerFuncs.UpdateFunc(&testProxyOne, &testProxyTwo)
	assertLog(t, zapcore.ErrorLevel, "there exists an APIProxy as the targeted path - APIProxy can not be updated - consider using kanalictl to avoid this error in the future", obsvr)

	apiProxyHandlerFuncs.DeleteFunc(&testProxyOne)
	assertLog(t, zapcore.DebugLevel, "deleted ApiProxy testProxyOne in foo namespace", obsvr)
	result, _ = spec.ProxyStore.Get("/api/v1/accounts")
	assert.Nil(t, result)

	apiProxyHandlerFuncs.DeleteFunc(testProxyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiProxy from k8s apiserver", obsvr)
}

func TestAPIKeyHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)
	setDecryptionKey(t)

	testKeyOne := spec.APIKey{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testKeyOne",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba",
		},
	}

	testKeyTwo := testKeyOne
	testKeyTwo.Spec.APIKeyData = "badencryption"

	apiKeyHandlerFuncs.AddFunc(&testKeyOne)
	result, err := spec.KeyStore.Get("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "added ApiKey testKeyOne", obsvr)

	apiKeyHandlerFuncs.AddFunc(testKeyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKey from k8s apiserver", obsvr)

	apiKeyHandlerFuncs.AddFunc(&testKeyTwo)
	assertLog(t, zapcore.ErrorLevel, "encoding/hex: odd length hex string", obsvr)

	testKeyOne.Spec.APIKeyData = "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba"
	apiKeyHandlerFuncs.UpdateFunc(&testKeyOne, &testKeyOne)
	result, err = spec.KeyStore.Get("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "updated ApiKey testKeyOne", obsvr)

	apiKeyHandlerFuncs.UpdateFunc(testKeyOne, &testKeyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKey from k8s apiserver", obsvr)

	apiKeyHandlerFuncs.UpdateFunc(&testKeyOne, testKeyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKey from k8s apiserver", obsvr)

	apiKeyHandlerFuncs.UpdateFunc(&testKeyTwo, &testKeyTwo)
	assertLog(t, zapcore.ErrorLevel, "encoding/hex: odd length hex string", obsvr)

	testKeyOne.Spec.APIKeyData = "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba"
	apiKeyHandlerFuncs.DeleteFunc(&testKeyOne)
	result, err = spec.KeyStore.Get("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9")
	assert.Nil(t, err)
	assert.Nil(t, result)
	assertLog(t, zapcore.DebugLevel, "deleted ApiKey testKeyOne", obsvr)

	apiKeyHandlerFuncs.DeleteFunc(testKeyOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKey from k8s apiserver", obsvr)

	apiKeyHandlerFuncs.DeleteFunc(&testKeyTwo)
	assertLog(t, zapcore.ErrorLevel, "encoding/hex: odd length hex string", obsvr)
}

func TestAPIKeyBindingHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	testBindingOne := spec.APIKeyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testBindingOne",
			Namespace: "foo",
		},
		Spec: spec.APIKeyBindingSpec{
			APIProxyName: "testProxyOne",
			Keys: []spec.Key{
				{
					Name: "franks-api-key",
				},
			},
		},
	}

	apiKeyBindingHandlerFuncs.AddFunc(&testBindingOne)
	result, err := spec.BindingStore.Get("testProxyOne", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "added ApiKeyBinding testBindingOne in foo namespace", obsvr)

	apiKeyBindingHandlerFuncs.AddFunc(testBindingOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKeyBinding from k8s apiserver", obsvr)

	apiKeyBindingHandlerFuncs.UpdateFunc(&testBindingOne, &testBindingOne)
	assertLog(t, zapcore.DebugLevel, "updated ApiKeyBinding testBindingOne in foo namespace", obsvr)

	apiKeyBindingHandlerFuncs.UpdateFunc(&testBindingOne, testBindingOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKeyBinding from k8s apiserver", obsvr)

	apiKeyBindingHandlerFuncs.UpdateFunc(testBindingOne, &testBindingOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKeyBinding from k8s apiserver", obsvr)

	apiKeyBindingHandlerFuncs.DeleteFunc(&testBindingOne)
	assertLog(t, zapcore.DebugLevel, "deleted ApiKeyBinding testBindingOne in foo namespace", obsvr)

	apiKeyBindingHandlerFuncs.DeleteFunc(testBindingOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ApiKeyBinding from k8s apiserver", obsvr)
	result, err = spec.BindingStore.Get("testProxyOne", "foo")
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestSecretHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	testSecretOne := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testSecretOne",
			Namespace: "foo",
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.key": []byte("YWJjMTIz"),
			"tls.crt": []byte("ZGVmNDU2"),
		},
	}

	secretHandlerFuncs.AddFunc(&testSecretOne)
	result, err := spec.SecretStore.Get("testSecretOne", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "added Secret testSecretOne in foo namespace", obsvr)

	secretHandlerFuncs.AddFunc(testSecretOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Secret from k8s apiserver", obsvr)

	secretHandlerFuncs.UpdateFunc(&testSecretOne, &testSecretOne)
	assertLog(t, zapcore.DebugLevel, "updated Secret testSecretOne in foo namespace", obsvr)

	secretHandlerFuncs.UpdateFunc(&testSecretOne, testSecretOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Secret from k8s apiserver", obsvr)

	secretHandlerFuncs.UpdateFunc(testSecretOne, &testSecretOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Secret from k8s apiserver", obsvr)

	secretHandlerFuncs.DeleteFunc(&testSecretOne)
	assertLog(t, zapcore.DebugLevel, "deleted Secret testSecretOne in foo namespace", obsvr)
	result, err = spec.SecretStore.Get("testSecretOne", "foo")
	assert.Nil(t, err)
	assert.Nil(t, result)

	secretHandlerFuncs.DeleteFunc(testSecretOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Secret from k8s apiserver", obsvr)
}

func TestServiceHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	testServiceOne := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testServiceOne",
			Namespace: "foo",
			Labels: map[string]string{
				"release":  "production",
				"name-two": "value-two",
			},
		},
	}

	serviceHandlerFuncs.AddFunc(&testServiceOne)
	result, err := spec.ServiceStore.Get(spec.CreateService(testServiceOne), http.Header{})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "added Service testServiceOne in foo namespace", obsvr)

	serviceHandlerFuncs.AddFunc(testServiceOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Service from k8s apiserver", obsvr)

	serviceHandlerFuncs.UpdateFunc(&testServiceOne, &testServiceOne)
	assertLog(t, zapcore.DebugLevel, "updated Service testServiceOne in foo namespace", obsvr)

	serviceHandlerFuncs.UpdateFunc(&testServiceOne, testServiceOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Service from k8s apiserver", obsvr)

	serviceHandlerFuncs.UpdateFunc(testServiceOne, &testServiceOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Service from k8s apiserver", obsvr)

	serviceHandlerFuncs.DeleteFunc(&testServiceOne)
	assertLog(t, zapcore.DebugLevel, "deleted Service testServiceOne in foo namespace", obsvr)
	result, err = spec.ServiceStore.Get(spec.CreateService(testServiceOne), http.Header{})
	assert.Nil(t, err)
	assert.Nil(t, result)

	serviceHandlerFuncs.DeleteFunc(testServiceOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Service from k8s apiserver", obsvr)
}

func TestConfigMapHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	mockOne, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	testConfigMapOne := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testConfigMapOne",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": string(mockOne),
		},
	}

	configMapHandlerFuncs.AddFunc(&testConfigMapOne)
	result, err := spec.MockResponseStore.Get("foo", "testConfigMapOne", "/foo", "GET")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assertLog(t, zapcore.DebugLevel, "added ConfigMap testConfigMapOne in foo namespace", obsvr)

	configMapHandlerFuncs.AddFunc(testConfigMapOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ConfigMap from k8s apiserver", obsvr)

	configMapHandlerFuncs.UpdateFunc(&testConfigMapOne, &testConfigMapOne)
	assertLog(t, zapcore.DebugLevel, "updated ConfigMap testConfigMapOne in foo namespace", obsvr)

	configMapHandlerFuncs.UpdateFunc(testConfigMapOne, &testConfigMapOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ConfigMap from k8s apiserver", obsvr)

	configMapHandlerFuncs.UpdateFunc(&testConfigMapOne, testConfigMapOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ConfigMap from k8s apiserver", obsvr)

	configMapHandlerFuncs.DeleteFunc(&testConfigMapOne)
	result, err = spec.MockResponseStore.Get("foo", "testConfigMapOne", "/foo", "GET")
	assert.Nil(t, err)
	assert.Nil(t, result)
	assertLog(t, zapcore.DebugLevel, "deleted ConfigMap testConfigMapOne in foo namespace", obsvr)

	configMapHandlerFuncs.DeleteFunc(testConfigMapOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed ConfigMap from k8s apiserver", obsvr)
}

func TestEndpointsHandlerFuncs(t *testing.T) {
	viper.SetDefault(config.FlagProcessLogLevel.GetLong(), "debug")
	defer viper.Reset()
	core, obsvr := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	logging.Init(core)

	testEndpointsOne := v1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kanali",
			Namespace: "foo",
		},
	}

	assert.Equal(t, spec.KanaliEndpoints, &v1.Endpoints{})
	endpointsHandlerFuncs.AddFunc(&testEndpointsOne)
	assert.NotNil(t, spec.KanaliEndpoints)
	assertLog(t, zapcore.DebugLevel, "adding Endpoints kanali in foo namespace", obsvr)

	endpointsHandlerFuncs.AddFunc(testEndpointsOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Endpoints from k8s apiserver", obsvr)

	endpointsHandlerFuncs.UpdateFunc(&testEndpointsOne, &testEndpointsOne)
	assertLog(t, zapcore.DebugLevel, "updated Endpoints kanali in foo namespace", obsvr)

	endpointsHandlerFuncs.UpdateFunc(testEndpointsOne, &testEndpointsOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Endpoints from k8s apiserver", obsvr)

	endpointsHandlerFuncs.UpdateFunc(&testEndpointsOne, testEndpointsOne)
	assertLog(t, zapcore.ErrorLevel, "received malformed Endpoints from k8s apiserver", obsvr)
}

func assertLog(t *testing.T, l zapcore.LevelEnabler, msg string, obsvr *observer.ObservedLogs) {
	assert.Equal(t, l, obsvr.All()[obsvr.Len()-1].Entry.Level)
	assert.Equal(t, msg, obsvr.All()[obsvr.Len()-1].Entry.Message)
	assert.Equal(t, 0, len(obsvr.All()[obsvr.Len()-1].Context))
}

func setDecryptionKey(t *testing.T) {
	// setup a private key
	rawPrivateKeyData := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAqgey9x7FDnREl2eHX80OoF5ZH1zM37Na7O+F6VY7jcQDRZEd
rIWCuASeypCfcxpIFGauLsjZ1uWaDj8i78nU5kpN01YDwlidMu0m5wBJbQBwApWp
53BqPOjKyMrLKBBE3roct1Lya4gUnKtvsKsq1ZukmCL2mZR63LqUnGwPUV+DhR8e
aCB+c2cRmUjgmuO+EcfBWgETOEWw//BGEnnKxtNlNj/T0qevvQWBvmLgcY2GV4t1
qd82fTQ9dCRP8pGgEAYaCospaF0SIhwdkD+IMNBxswvNypFHe+ee10H40bab/HtL
fgD+0DdAcyk30Xd5UowOL0G5I9Om4P9xDR0OJQIDAQABAoIBAQCNjRtQ4Czte54e
7fGlr/EdUW6gzYFCOu7XkhDJ0SCDRUvz/nvVxNCuoioQOZaFHLwlP9aC3HN+lGdM
tQNA3aaAkht4dYNrqJo2a3iXl4WJWXPmsvJf2xaW3rbzsuCu2dx8EDUX6dYn74I5
a//v9JAUhR7iCTmDYjlmyW4vS0b1VlMY1fTRn2pV1MobOXAxtL4AeKckCZlqvvPs
vI7JtcZiUH5X6w2S069hQjcnSOBaeNm2h96sElywXLK7ocn8WLfiK87KhYb+iwV3
y4nLc4FWn1x84TeYUDr7KPKeNmJN9qUvHiyjq3cZeBh/kS3RKyJnmqIfxGOEjweA
81pMbfWBAoGBANy6Z2UqU5d8P1ENwruol3NVBr1hsjuLlPqyEr1RbRygXxhKo51Z
1YlvuLElrLI+EBXG5sZfoNqCNzrRjE2I/x6yD4YXTzQAbIIyuvepm02GTZkwCbeV
ogtnUgZs/+NeK5SPxL1SHnEC2CKYzniqz5X4wU89Sy1IU0fVs3cWmWr7AoGBAMUz
U8/gw8MhgDBUQN+GxRYWIpM5sYAxmxjlSTXYACSAvDe10EmzX0bvhqgpml4VCsh6
kHcxFtRintqU+/sQX1xOB8PkbKyQ2F6sfARPdN96wt5illXmOJDcHFzF9JsPuokf
0oCmSqFMUwVEnOYHKHjPEgUXBMt6RQG6ZFKM7CFfAoGAEGN+vdIwtrSw9tGwvcI8
zIFpqLISS72HaD3Eim37Tf7SULlqAIrlAflbmzuS4QEerJ9CD0SWauL0OtfINRGQ
77SE1o3wgYnuhJRDsZPURiNMq60Xa5welf4r5vI85aLeVX9/OhdYz0dIyPjTi74S
/62inMGOTWRrWpARbSwuY/UCgYAp849i69f5jQjuPx5F7y/CURct5toiAWuTUkPO
1pLBMTnZSC51X1lnh9FHuQID/cn7eEd/c8v9hrVwnr00XflLc7YnMFMGKxI8ockI
Sbb9S+pk0OhWgdGKz6ig++PbkO8H+OMZR0cdJVWMC1qtfNxZRwP4XTtB5lMD2Rk6
jow1VQKBgQDUj8sLDAVS3C64U5BCitlnQ4hdScpBq6P09B624A3EjseYwCWK1h7L
xpyQYWffQdGb6IoawFPWy7nT42d4IER+JmRJRKm2a7XD3LsNud5tBCB57ksN27v7
I0teIby/h2MN7yPNjFgIq64JneHpKwEH04+qy/dWQafxrpBlf9tfwg==
-----END RSA PRIVATE KEY-----`)

	block, _ := pem.Decode(rawPrivateKeyData)
	// parse the pem block into a private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	spec.APIKeyDecryptionKey = privateKey
}
