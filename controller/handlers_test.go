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
	"testing"

	"github.com/northwesternmutual/kanali/spec"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

func TestAddFunc(t *testing.T) {
	clearAllStores()
	defer clearAllStores()
	handlers := k8sEventHandler{}
	setDecryptionKey(t)

	handlers.addFunc(spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "api/v1/accounts",
			Service: spec.Service{
				Name: "my-service",
				Port: 8080,
			},
		},
	})
	result, err := spec.ProxyStore.Get("/api/v1/accounts")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.addFunc(spec.APIKey{
		ObjectMeta: api.ObjectMeta{
			Name: "def456",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "encrypted",
		},
	})
	handlers.addFunc(spec.APIKey{
		ObjectMeta: api.ObjectMeta{
			Name: "abc123",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba",
		},
	})
	result, err = spec.KeyStore.Get("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.addFunc(spec.APIKeyBinding{
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeyBindingSpec{
			APIProxyName: "api-proxy-one",
			Keys: []spec.Key{
				{
					Name: "franks-api-key",
				},
			},
		},
	})
	result, err = spec.BindingStore.Get("api-proxy-one", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.addFunc(api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name:      "secret-one",
			Namespace: "foo",
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.key": []byte("YWJjMTIz"),
			"tls.crt": []byte("ZGVmNDU2"),
		},
	})
	result, err = spec.SecretStore.Get("secret-one", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.addFunc(api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
			Labels: map[string]string{
				"release":  "production",
				"name-two": "value-two",
			},
		},
	})
	result, err = spec.ServiceStore.Get(spec.Service{
		Namespace: "bar",
		Labels: spec.Labels{
			spec.Label{
				Name:  "release",
				Value: "production",
			},
			spec.Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	ep := api.Endpoints{
		ObjectMeta: api.ObjectMeta{
			Name:      "kanali",
			Namespace: "bar",
		},
	}
	handlers.addFunc(ep)
	assert.Equal(t, *(spec.KanaliEndpoints), ep)

	mockOne, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	handlers.addFunc(api.ConfigMap{
		ObjectMeta: api.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": string(mockOne),
		},
	})
	result, err = spec.MockResponseStore.Get("foo", "cm-one", "/foo", "GET")
	assert.Nil(t, err)
	assert.NotNil(t, result)

}

func TestUpdateFunc(t *testing.T) {
	clearAllStores()
	defer clearAllStores()
	handlers := k8sEventHandler{}
	setDecryptionKey(t)

	spec.ProxyStore.Set(spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "api/v1/accounts",
			Service: spec.Service{
				Name: "my-service",
				Port: 8080,
			},
		},
	})
	result, _ := spec.ProxyStore.Get("/api/v1/accounts")
	assert.NotNil(t, result)
	handlers.updateFunc(spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "/modified/path",
			Service: spec.Service{
				Name: "my-service",
				Port: 8080,
			},
		},
	})
	result, _ = spec.ProxyStore.Get("/api/v1/accounts")
	assert.Nil(t, result)
	result, _ = spec.ProxyStore.Get("/modified/path")
	assert.NotNil(t, result)

	handlers.updateFunc(spec.APIKey{
		ObjectMeta: api.ObjectMeta{
			Name: "def456",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "encrypted",
		},
	})
	handlers.addFunc(spec.APIKey{
		ObjectMeta: api.ObjectMeta{
			Name: "abc123",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba",
		},
	})
	result, err := spec.KeyStore.Get("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.updateFunc(spec.APIKeyBinding{
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeyBindingSpec{
			APIProxyName: "api-proxy-one",
			Keys: []spec.Key{
				{
					Name: "franks-api-key",
				},
			},
		},
	})
	result, err = spec.BindingStore.Get("api-proxy-one", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.updateFunc(api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name:      "secret-one",
			Namespace: "foo",
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.key": []byte("YWJjMTIz"),
			"tls.crt": []byte("ZGVmNDU2"),
		},
	})
	result, err = spec.SecretStore.Get("secret-one", "foo")
	assert.Nil(t, err)
	assert.NotNil(t, result)

	handlers.updateFunc(api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
			Labels: map[string]string{
				"release":  "production",
				"name-two": "value-two",
			},
		},
	})
	result, err = spec.ServiceStore.Get(spec.Service{
		Namespace: "bar",
		Labels: spec.Labels{
			spec.Label{
				Name:  "release",
				Value: "production",
			},
			spec.Label{
				Name:  "name-two",
				Value: "value-two",
			},
		},
	}, nil)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	ep := api.Endpoints{
		ObjectMeta: api.ObjectMeta{
			Name:      "kanali",
			Namespace: "bar",
		},
	}
	handlers.updateFunc(ep)
	assert.Equal(t, *(spec.KanaliEndpoints), ep)

	mockOne, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
	})

	handlers.updateFunc(api.ConfigMap{
		ObjectMeta: api.ObjectMeta{
			Name:      "cm-one",
			Namespace: "foo",
		},
		Data: map[string]string{
			"response": string(mockOne),
		},
	})
	result, err = spec.MockResponseStore.Get("foo", "cm-one", "/foo", "GET")
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func clearAllStores() {
	spec.ProxyStore.Clear()
	spec.KeyStore.Clear()
	spec.BindingStore.Clear()
	spec.SecretStore.Clear()
	spec.ServiceStore.Clear()
	spec.MockResponseStore.Clear()
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
