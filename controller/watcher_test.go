// // Copyright (c) 2017 Northwestern Mutual.
// //
// // Permission is hereby granted, free of charge, to any person obtaining a copy
// // of this software and associated documentation files (the "Software"), to deal
// // in the Software without restriction, including without limitation the rights
// // to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// // copies of the Software, and to permit persons to whom the Software is
// // furnished to do so, subject to the following conditions:
// //
// // The above copyright notice and this permission notice shall be included in
// // all copies or substantial portions of the Software.
// //
// // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// // IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// // FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// // AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// // LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// // OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// // THE SOFTWARE.
//
package controller

//
// import (
// 	"crypto/x509"
// 	"encoding/pem"
// 	"net/http"
// 	"testing"
// 	"time"
//
// 	"github.com/northwesternmutual/kanali/spec"
// 	"github.com/stretchr/testify/assert"
// 	"k8s.io/kubernetes/pkg/api"
// 	"k8s.io/kubernetes/pkg/api/unversioned"
// )
//
// func TestMonitor(t *testing.T) {
//
// 	assert := assert.New(t)
//
// 	testEventCh := make(chan *event)
// 	monitor(testEventCh)
// 	testAPIProxy := getTestAPIProxy()
// 	testAPIKeyBinding := getTestAPIKeyBinding()
// 	testAPIKey := getTestAPIKey()
// 	testService := getTestService()
// 	testSecret := getTestSecret()
//
// 	testEventCh <- &event{added, testAPIProxy}
// 	time.Sleep(100 * time.Millisecond)
// 	proxy, _ := spec.ProxyStore.Get(testAPIProxy.Spec.Path)
// 	assert.Equal(testAPIProxy, proxy)
//
// 	testAPIProxy.Spec.Target = "/foo"
// 	testEventCh <- &event{modified, testAPIProxy}
// 	time.Sleep(100 * time.Millisecond)
// 	proxy, _ = spec.ProxyStore.Get(testAPIProxy.Spec.Path)
// 	assert.Equal(testAPIProxy, proxy)
//
// 	testEventCh <- &event{deleted, testAPIProxy}
// 	time.Sleep(100 * time.Millisecond)
// 	assert.Nil(spec.ProxyStore.Get(testAPIProxy.Spec.Path))
//
// 	testEventCh <- &event{added, testAPIKeyBinding}
// 	time.Sleep(100 * time.Millisecond)
// 	binding, _ := spec.BindingStore.Get(testAPIKeyBinding.Spec.APIProxyName, testAPIKeyBinding.ObjectMeta.Namespace)
// 	assert.Equal(testAPIKeyBinding, binding)
//
// 	testAPIKeyBinding.Spec.Keys[0].DefaultRule.Global = false
// 	testEventCh <- &event{modified, testAPIKeyBinding}
// 	time.Sleep(100 * time.Millisecond)
// 	binding, _ = spec.BindingStore.Get(testAPIKeyBinding.Spec.APIProxyName, testAPIKeyBinding.ObjectMeta.Namespace)
// 	assert.Equal(testAPIKeyBinding, binding)
//
// 	testEventCh <- &event{deleted, testAPIKeyBinding}
// 	time.Sleep(100 * time.Millisecond)
// 	assert.Nil(spec.BindingStore.Get(testAPIKeyBinding.Spec.APIProxyName, testAPIKeyBinding.ObjectMeta.Namespace))
//
// 	// setup a private key
// 	rawPrivateKeyData := []byte(`-----BEGIN RSA PRIVATE KEY-----
// MIIEowIBAAKCAQEAqYdnX0jeOX0zZuTG0zDJ+t1qzA63MMxYllwcNdSIuDCvT6Rn
// wSg0nx+PSWYQQqXCN7q8CSsBgp6QNneCxL3A/1JzV7w/fMyWLIuuSOn7Gi8Iz+7E
// Mb9dbPzejHbx44TDzjIm++xwpyR56e6Zqi8h+XFfNTxQ1IWsiUQJsEvNOus9km92
// gTZ9hJNX8GgfCvuP0BjDsXGjQVhUu7tLO4eccXvZjnLLYrOM9qtkkEf8eD/1dh5+
// nvmgCl2QI9YHq+OvlCLtAc2m1txPxtvCur51RjoPUXkmCgljZdqadSKcURu/DXME
// cLF13smQl6Jq4gGzQY919PC1EjkhGKZA/EFtUwIDAQABAoIBADN6r5RKr1irwTkk
// jY/CCAOKywxuB4jk9J2sGNDr2hx8hC/eD7ei+t+7GKrEOHnUlfaQWNs72PiOJ+Ky
// Rd5ydLHTmrzwqCLAiXW7cNApZRvdXoKt0Zv9rWQUIYxr7iYVwdPSfO4RLWBD/lVg
// I/9+0oVJvQyQZUcz1GHWbE7Bpe+W0vkDeFXxlCP39UWmzfCChhzCFXTTgvl2EHdx
// QMPnn3dhjf+uBtZZXUpjo6lFNIrWlqnSgE7krJUfzD5TVgG1q8AfF5BdYmAtEoIa
// nWQrn51++seJQcCh0g4bRWbjk79Qp2uanoeVZpZdRQSaguUnJbOLXkSJAiqAczoD
// MWqXSWkCgYEA2SpJMMdN30yjTZeHPLb72wum7l+ZZ6Wrv8zt6AWPA3T7i6dBRSKi
// D6ycCLk6V6SWOEgm2MruYxwvd3lY3XHjJwmfpWqQjYrKp1u6w+B8NKhcdAx2lcn3
// Uv7rv1A/Et7Q1abzL9e/vteP0o5sYDvfDGxMFt9jgu43S/rpHOWTe10CgYEAx9ha
// iHHcPfpysZR5kX0eOT8pXpOG6UfrEw38yZqIyCGz4fyZKQ+bttg0WqHMpWe1mw9P
// pYBq2PtzUuyJVGft9xTP7Ov1oa741cYFpBbYtW4e81CWPBj4cZNBzq1Y0dg4Hw0s
// aGYBQ9L87koM0elfjNo2/HJfTVc8OWU60EWLsG8CgYBg2BSvphG6JQkmTw7GKqwC
// MR4Oa5+TszP2YsMtl10Bo6eRzdKzrBAtgUJMOZ4k+4bqLnL0dvr8Q9N/KiRRDLrJ
// 6+a/89fm5yAcpjGRrIh3SyV/sxcnEVw0LO6g8H5QQgFLZhpJGaOuzZ6bvVvjRo/f
// kGQWRySvfOA4B/rxIgg1GQKBgCM/VqBwLJ9F2ArYHCT8A2Onbz1+GbJ1e9GtiuNn
// /S4HO7nlGoJyfU1fjsRZe0XFJ/PEXJDdOHsyxmFe1M3tUrxckFvCNl2hBcR2m7IY
// UXqWhKD3mrfY06D8jwPL8Tl5wFRBt45mR1zWDsRcjSxM1Ax8xGv8JDD47OdWomvv
// iDbDAoGBAK/WpzoN59vsaw+PcxS93WgBAnGe9q+mPxNPLj0xaakiNFn7FCdiGV8z
// 4Wn0sua48v7QcJbKcL0ZbXky61EZV3HqAyzhWZ6jSSQoM37S/nPpq5EPTfs4Dvnl
// BA7dWeyLFnN7ePVSkL1SES58MpMMiunrUo/Ci6OjiOyN7ynnU6tE
// -----END RSA PRIVATE KEY-----`)
//
// 	block, _ := pem.Decode(rawPrivateKeyData)
// 	// parse the pem block into a private key
// 	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		assert.Fail(err.Error())
// 	}
//
// 	spec.APIKeyDecryptionKey = privateKey
//
// 	testEventCh <- &event{added, testAPIKey}
// 	time.Sleep(100 * time.Millisecond)
//
// 	result, _ := spec.KeyStore.Get("CByPKc5XL8U77Ag8c6RbVgKJkJ9B2J5b")
//
// 	assert.Equal(spec.APIKey{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "abc123",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeySpec{
// 			APIKeyData: "CByPKc5XL8U77Ag8c6RbVgKJkJ9B2J5b",
// 		},
// 	}, result)
//
// 	testAPIKey.Spec.APIKeyData = "3e4d0a287d22be64eb43cca350a27204c51a9da07a37d9d5ec6e5689fc4d659ea7931d814935c24b5a4fb3a724cf40e76c597933b47c1c9b86e04a7a895ec305c4c0b1c4e2066e2fc319ec077adfb12dc267cb1443c0cd2bc4195b3655e418ac45137f475ff6fa5f6b1394b7a63b5f31dca7c3fc6846ca3794c2b7da66e517c4540c5dd51299b6dc65a65e947545601d146b6af7b45aa51f869f13635eb3a8b2a9d7bbce4f7895ade3509333a0b707e49427bfd2da168258095ba585910c99da70c7106e57e9e52ae6556d6626253856afb0a8fcb8b18a3f9b83adeab19fcfe8baf2898d143824aac5450d8b5bb9f772f44614ea5958ff5774c328432a2ec38d"
// 	testEventCh <- &event{modified, testAPIKey}
// 	time.Sleep(100 * time.Millisecond)
// 	result, _ = spec.KeyStore.Get("WS6JlCUtOYBB3cxokdcJAnkIciZaOwyS")
// 	assert.Equal(spec.APIKey{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "abc123",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeySpec{
// 			APIKeyData: "WS6JlCUtOYBB3cxokdcJAnkIciZaOwyS",
// 		},
// 	}, result)
//
// 	testEventCh <- &event{deleted, testAPIKey}
// 	time.Sleep(100 * time.Millisecond)
// 	assert.Nil(spec.KeyStore.Get("WS6JlCUtOYBB3cxokdcJAnkIciZaOwyS"))
//
// 	testEventCh <- &event{added, testService}
// 	time.Sleep(100 * time.Millisecond)
// 	untypedResult, _ := spec.ServiceStore.Get(spec.CreateService(testService), http.Header{})
// 	resultSvc := untypedResult.(spec.Service)
// 	assert.Equal("1.2.3.4", resultSvc.ClusterIP)
//
// 	testService.Spec.ClusterIP = "2.3.4.5"
// 	testEventCh <- &event{modified, testService}
// 	time.Sleep(100 * time.Millisecond)
// 	untypedResult, _ = spec.ServiceStore.Get(spec.CreateService(testService), http.Header{})
// 	resultSvc = untypedResult.(spec.Service)
// 	assert.Equal("2.3.4.5", resultSvc.ClusterIP)
//
// 	testEventCh <- &event{deleted, testService}
// 	time.Sleep(100 * time.Millisecond)
// 	untypedResult, _ = spec.ServiceStore.Get(spec.CreateService(testService), http.Header{})
// 	assert.Nil(untypedResult)
//
// 	testEventCh <- &event{added, testSecret}
// 	time.Sleep(100 * time.Millisecond)
// 	result, _ = spec.SecretStore.Get(testSecret.ObjectMeta.Name, testSecret.ObjectMeta.Namespace)
// 	assert.Equal(testSecret, result)
//
// 	testSecret.Data["tls.ca"] = []byte("YWJjMTIz")
// 	testEventCh <- &event{modified, testSecret}
// 	time.Sleep(100 * time.Millisecond)
// 	result, _ = spec.SecretStore.Get(testSecret.ObjectMeta.Name, testSecret.ObjectMeta.Namespace)
// 	assert.Equal(testSecret, result)
//
// 	testEventCh <- &event{deleted, testSecret}
// 	time.Sleep(100 * time.Millisecond)
// 	assert.Nil(spec.SecretStore.Get(testSecret.ObjectMeta.Name, testSecret.ObjectMeta.Namespace))
//
// }
<<<<<<< HEAD
//
// func getTestSecret() api.Secret {
// 	return api.Secret{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "secret-two",
// 			Namespace: "foo",
// 		},
// 		Type: "kubernetes.io/tls",
// 		Data: map[string][]byte{
// 			"tls.key": []byte("YWJjMTIz"),
// 			"tls.crt": []byte("ZGVmNDU2"),
// 		},
// 	}
// }
//
// func getTestService() api.Service {
// 	return api.Service{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "foo",
// 			Namespace: "bar",
// 			Labels: map[string]string{
// 				"one":   "two",
// 				"three": "four",
// 			},
// 		},
// 		Spec: api.ServiceSpec{
// 			ClusterIP: "1.2.3.4",
// 		},
// 	}
// }
//
// func getTestAPIKey() spec.APIKey {
// 	return spec.APIKey{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "abc123",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeySpec{
// 			APIKeyData: "5e0329991f70e018f9503aa9c8f03ebe861df9be92b91dd5188baecd975572ab4d9973292768991200f433ef8830dfd914b134886ccb13b1cb31925028b44b10a03845718ee96db46cc738a294ff4ebd756d84ce3124e5f25d2818af7e325d16805977f658a6331e7b82db77a4366285f817df4394c45e6fb09ce9764e6813bf5ba449e2377f496bcdda07d2e27561b1c215582a1560f7b76ed5d5db29bd6d57d93e8068bb776bd7fd83a2e99319a8ff1ca27250b683a96af655566147302d75db31af3d6c0b27e9e1f1a471eea7f6cc31343b0caad14ad3ba8be8eabe3add5d9dd25594290325a1d4bdaa45c8502fd8edad015c93646aae513ac1e2cf22377b",
// 		},
// 	}
// }
//
// func getTestAPIKeyBinding() spec.APIKeyBinding {
// 	return spec.APIKeyBinding{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "abc123",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIKeyBindingSpec{
// 			APIProxyName: "api-proxy-one",
// 			Keys: []spec.Key{
// 				{
// 					Name:        "franks-api-key",
// 					DefaultRule: spec.Rule{},
// 					Subpaths: []*spec.Path{
// 						{
// 							Path: "/foo",
// 							Rule: spec.Rule{
// 								Global: true,
// 							},
// 						},
// 						{
// 							Path: "foo/bar",
// 							Rule: spec.Rule{
// 								Granular: &spec.GranularProxy{
// 									Verbs: []string{
// 										"POST",
// 										"GET",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }
//
// func getTestAPIProxy() spec.APIProxy {
// 	return spec.APIProxy{
// 		TypeMeta: unversioned.TypeMeta{},
// 		ObjectMeta: api.ObjectMeta{
// 			Name:      "exampleAPIProxyOne",
// 			Namespace: "foo",
// 		},
// 		Spec: spec.APIProxySpec{
// 			Path: "api/v1/accounts",
// 			Hosts: []spec.Host{
// 				{
// 					Name: "foo.bar.com",
// 					SSL: spec.SSL{
// 						SecretName: "mySecretTwo",
// 					},
// 				},
// 				{
// 					Name: "bar.foo.com",
// 				},
// 			},
// 			Service: spec.Service{
// 				Name:      "my-service",
// 				Namespace: "foo",
// 				Port:      8080,
// 			},
// 			Plugins: []spec.Plugin{
// 				{
// 					Name:    "apikey",
// 					Version: "1.0.0",
// 				},
// 				{
// 					Name: "jwt",
// 				},
// 			},
// 			SSL: spec.SSL{
// 				SecretName: "mySecret",
// 			},
// 		},
// 	}
// }
=======

func getTestSecret() api.Secret {
	return api.Secret{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "secret-two",
			Namespace: "foo",
		},
		Type: "kubernetes.io/tls",
		Data: map[string][]byte{
			"tls.key": []byte("YWJjMTIz"),
			"tls.crt": []byte("ZGVmNDU2"),
		},
	}
}

func getTestService() api.Service {
	return api.Service{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "foo",
			Namespace: "bar",
			Labels: map[string]string{
				"one":   "two",
				"three": "four",
			},
		},
		Spec: api.ServiceSpec{
			ClusterIP: "1.2.3.4",
		},
	}
}

func getTestAPIKey() spec.APIKey {
	return spec.APIKey{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeySpec{
			APIKeyData: "5e0329991f70e018f9503aa9c8f03ebe861df9be92b91dd5188baecd975572ab4d9973292768991200f433ef8830dfd914b134886ccb13b1cb31925028b44b10a03845718ee96db46cc738a294ff4ebd756d84ce3124e5f25d2818af7e325d16805977f658a6331e7b82db77a4366285f817df4394c45e6fb09ce9764e6813bf5ba449e2377f496bcdda07d2e27561b1c215582a1560f7b76ed5d5db29bd6d57d93e8068bb776bd7fd83a2e99319a8ff1ca27250b683a96af655566147302d75db31af3d6c0b27e9e1f1a471eea7f6cc31343b0caad14ad3ba8be8eabe3add5d9dd25594290325a1d4bdaa45c8502fd8edad015c93646aae513ac1e2cf22377b",
		},
	}
}

func getTestAPIKeyBinding() spec.APIKeyBinding {
	return spec.APIKeyBinding{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "abc123",
			Namespace: "foo",
		},
		Spec: spec.APIKeyBindingSpec{
			APIProxyName: "api-proxy-one",
			Keys: []spec.Key{
				{
					Name:        "franks-api-key",
					DefaultRule: spec.Rule{},
					Subpaths: []*spec.Path{
						{
							Path: "/foo",
							Rule: spec.Rule{
								Global: true,
							},
						},
						{
							Path: "foo/bar",
							Rule: spec.Rule{
								Granular: &spec.GranularProxy{
									Verbs: []string{
										"POST",
										"GET",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getTestAPIProxy() spec.APIProxy {
	return spec.APIProxy{
		TypeMeta: unversioned.TypeMeta{},
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path: "api/v1/accounts",
			Hosts: []spec.Host{
				{
					Name: "foo.bar.com",
					SSL: spec.SSL{
						SecretName: "mySecretTwo",
					},
				},
			},
			Service: spec.Service{
				Name:      "my-service",
				Namespace: "foo",
				Port:      8080,
			},
			Plugins: []spec.Plugin{
				{
					Name:    "apikey",
					Version: "1.0.0",
				},
				{
					Name: "jwt",
				},
			},
			SSL: spec.SSL{
				SecretName: "mySecret",
			},
		},
	}
}
>>>>>>> master
