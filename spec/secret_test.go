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

package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSecretGetSecretStore(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore

	store.Clear()
	assert.Equal(0, len(store.secretMap), "store should be empty")

	v := SecretFactory{}
	var i interface{} = &v
	_, ok := i.(Store)
	assert.True(ok, "SecretFactory does not implement the Store interface")
}

func TestSecretSet(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	store.Set(secretList[0])
	store.Set(secretList[1])
	err := store.Set(APIProxy{})
	assert.Equal(err.Error(), "grrr - you're only allowed add secrets to the secrets store.... duh", "error not what expected")
	assert.Equal(1, len(store.secretMap), "there should be one namespace in the secret store")
	assert.Equal(2, len(store.secretMap["foo"]), "there should be two secrets in the secret store")
	assert.Equal(secretList[0], store.secretMap["foo"]["secret-one"], "secret should exist")
	assert.Equal(secretList[1], store.secretMap["foo"]["secret-two"], "secret should exist")
}

func TestSecretUpdate(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	store.Update(secretList[0], secretList[0])
	store.Update(secretList[1], secretList[1])
	err := store.Update(APIProxy{}, APIProxy{})
	assert.Equal(err.Error(), "grrr - you're only allowed add secrets to the secrets store.... duh", "error not what expected")
	assert.Equal(1, len(store.secretMap), "there should be one namespace in the secret store")
	assert.Equal(2, len(store.secretMap["foo"]), "there should be two secrets in the secret store")
	assert.Equal(secretList[0], store.secretMap["foo"]["secret-one"], "secret should exist")
	assert.Equal(secretList[1], store.secretMap["foo"]["secret-two"], "secret should exist")
}

func TestSecretClear(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Set(secretList[0])
	store.Clear()
	assert.Equal(0, len(store.secretMap), "store should be empty")
}

func TestSecretIsEmpty(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	assert.True(store.IsEmpty())
	store.Set(secretList[0])
	assert.False(store.IsEmpty())
	store.Clear()
	assert.True(store.IsEmpty())
}

func TestSecretGet(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	store.Set(secretList[0])
	store.Set(secretList[1])
	result, _ := store.Get("", "")
	assert.Nil(result, "empty string should not exist")
	result, _ = store.Get("jkl012", "")
	assert.Nil(result, "secrect should not exist")
	_, err := store.Get("jkl012")
	assert.Equal(err.Error(), "should should take 2 params, name and namespace", "error not what expected")
	_, err = store.Get(5, "")
	assert.Equal(err.Error(), "secret name must be of type string", "error not what expected")
	_, err = store.Get("", 5)
	assert.Equal(err.Error(), "secret namespace must be of type string", "error not what expected")
	result, _ = store.Get("secret-one", "foo")
	assert.Equal(secretList[0], result, "secret should exist")
	result, _ = store.Get("secret-two", "foo")
	assert.Equal(secretList[1], result, "secret should exist")
}

func TestSecretDelete(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	store.Set(secretList[0])
	store.Set(secretList[1])
	_, err := store.Delete(APIProxy{})
	assert.Equal(err.Error(), "there's no way this secret could've gotten in here", "wrong error")
	result, _ := store.Delete(nil)
	assert.Nil(result, "nil should be nil")
	result, _ = store.Delete(secretList[1])
	assert.Equal(secretList[1], result, "secret that is deleted should be returned")
	result, _ = store.Delete(secretList[1])
	assert.Nil(result, "deleted secret should return nil")
	assert.Equal(1, len(store.secretMap), "store should have one namespace represented")
	store.Set(secretList[1])
	assert.Equal(2, len(store.secretMap["foo"]), "store should have two secrets represented in this namespace")
	result, _ = store.Delete(secretList[3])
	assert.Nil(result, "secret namespace should not exist")
	store.Set(secretList[3])
	assert.Equal(2, len(store.secretMap), "should have 2 namespaces represented")
	store.Delete(secretList[3])
	assert.Equal(1, len(store.secretMap), "should have 1 namespace represented")
}

func TestX509KeyPair(t *testing.T) {
	assert := assert.New(t)
	store := SecretStore
	secretList := getTestSecretList()

	store.Clear()
	store.Set(secretList[0])
	store.Set(secretList[1])
	store.Set(secretList[2])

	result, _ := store.Get("secret-one", "foo")
	typed := result.(v1.Secret)
	_, err := X509KeyPair(typed)
	assert.Nil(err, "should be able to get key pair")
	result, _ = store.Get("secret-two", "foo")
	typed = result.(v1.Secret)
	_, err = X509KeyPair(typed)
	assert.Equal(err.Error(), "tls: failed to find any PEM data in certificate input", "wrong error")
	result, _ = store.Get("secret-three", "foo")
	typed = result.(v1.Secret)
	_, err = X509KeyPair(typed)
	assert.Equal(err.Error(), "tls: failed to find any PEM data in certificate input", "wrong error")
}

func getTestSecretList() []v1.Secret {

	return []v1.Secret{
		{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-one",
				Namespace: "foo",
			},
			Type: "kubernetes.io/tls",
			Data: map[string][]byte{
				"tls.key": []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAt4Mb8EJTLdmvEYmOMv2e1xgHDW8zW+fl5N+hEjvPO1hj7XE/
IiPmPpyGbBHiXnHWM3PRTHXw85GZNAv1m7CFYRap7qFjjFxx5PpHv2yrAk7VOaM5
vW4saBDW71PYvYV3Sne7poj360bfPjjAq9EPx2BACzYtKXSXVegsMF99TiRDU60g
U4PmyUyKWHqHjIXORxrDYwQAEKLrns+86+077Pek6EitMb0j5RJGoLoEozpQo783
qCmN3MwwUyjqZWbAymwXbeDNpLK8fyGeumm/rZ5uWAFY0bjiKlJqj5H9dL0GLFs4
Wj/b1E167Hw0fpL2FQcTdceyyDEeXTP7yHTJbQIDAQABAoIBAD9GSKrV46JKjY+p
c0tnoY4ercfXEMlesyjwZrRHxRN8TbBAPPmDkU8DX2IB5KCV9bp187O+ct84069b
dEtDWSn5y4wJP67U7Tx7g6OgT3KTRfgCrNUBGPSk4cdRjMkELPaTRDOOxJTuMeld
XimgAyXGrS5wdq/1kHzBegMP2b/IMaYXWqvbZzvVVhM62HhfCoy+LVO6wAnyzg+4
j0PgogLuMsZqV9sJKte5dDjYMUGU4eqT3qGw/c7u8Nc2wXDt9PEDa8I0dbDt8aJO
hwFz/cGGN6OBmucPr8FpwDwoBAzW7zZVKr92d5wtgl/x6wnHcfoAJYeh3A2wy5Dx
0hJLbKECgYEA3AydRU814uw/AAdgYf+CEugRRuWFl/AKCK4RUEgRjdWbeW2gPorc
UXgJGqLrW9dqtxdDpnogSuB7MNJ3Z0jRttf/aibrNWInmIabMK34gzUPfpge6wAr
pMCAnDcD1ev/lgHGAOgjVAEI4pFGCPpHt9n7ZUOXrLlOVWKKJn4cjYkCgYEA1X5b
wxlaYqCbRxWgSOIiaV0OvfveAMESySZYaE62AzjKE0A7t3hQ5ok4B2grgpvdhHXc
ePPDr1vslLbiYLWQSyJOzhcL+19VoKd8xipRF/l+R/ibEDkAOw5JjpCnQcvXajFm
Cr726uxqtfQjOmryE2b3bfxzNTP0tF0KHhyMJ8UCgYB3hY6DEAQ37pRForwSXqM7
O7vuo0qN/X6luk6DKbIHNSgOq6Hndqs2wRlI04c9SmOXcVZ4fUs/AHFJUngyrJXk
V6xO7zbAt0DEkxCGP2iFc/Onkl7wzBXolfsFXsiBkH8x7mKrKWvO/ATkc33z8GW2
Eft/xFgFJF3M+QoK4GMH6QKBgH7HiKp09c233j36Q7LERvcWwdhkJ1N/VC6LTNi0
VTXYlXeVH2E2W6RrPBB59cybzpIIg6J3727FQdqWOByB9WwItY+A9CaUWH8MtvXe
z0VTbYBDy6F4mAB53YiBsjFKYWO4bgZXnL2pIz1l+z2HoLWQ4cI4thmT/u7UiVuU
TE6ZAoGBAINgO/tARvu0Za015acMq0U6eKIlxzEaWAIzWkFz+yOAu8S3d/5pmFlh
SkhyyDd+dXF1C6WDQd4z9jVR4KbE+4NEOUpFf/+5lgtaBm+cJpwxfhXhhOmdlTu7
3Gh4e/dRFE8pNugniSz5zPh0bCX81TsijgOpjN4LKiUIdgZ2R0dw
-----END RSA PRIVATE KEY-----`),
				"tls.crt": []byte(`-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAMhA5naPEAzGMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNzA0MjYxNTU4MjRaFw0yNzA0MjQxNTU4
MjRaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBALeDG/BCUy3ZrxGJjjL9ntcYBw1vM1vn5eTfoRI7zztY
Y+1xPyIj5j6chmwR4l5x1jNz0Ux18PORmTQL9ZuwhWEWqe6hY4xcceT6R79sqwJO
1TmjOb1uLGgQ1u9T2L2Fd0p3u6aI9+tG3z44wKvRD8dgQAs2LSl0l1XoLDBffU4k
Q1OtIFOD5slMilh6h4yFzkcaw2MEABCi657PvOvtO+z3pOhIrTG9I+USRqC6BKM6
UKO/N6gpjdzMMFMo6mVmwMpsF23gzaSyvH8hnrppv62eblgBWNG44ipSao+R/XS9
BixbOFo/29RNeux8NH6S9hUHE3XHssgxHl0z+8h0yW0CAwEAAaNQME4wHQYDVR0O
BBYEFD0PEoXECieWqj/QfXp3y0EaZzVxMB8GA1UdIwQYMBaAFD0PEoXECieWqj/Q
fXp3y0EaZzVxMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAHv4531M
rDfBHCgxMW/Vm+Lg2EIZI3V0L75GEDUmA/FJ3qF0A+fMxJtFOqbcEkQq6twCyq8v
BYA2A8twakVL3RMiJYdsSphX4Rxr26arWKpVLXbHLtu95p37am7AJZZRjMCkBNwa
fFCKBRPmL3H6UKrCt3stFfKWwbK/9nI+G7KzbbPaK/vaCbC0BUr7WObBTtB3erbW
f/AjKZmfokLhsMl6vhbIWz9QV3Ssyhbb8u+TmNQV79696yuJugPsM+gf5vYqj6M0
ls6gptH4+JWQAtYh6HsTl6caNRikSgXxIIae7cNcEJyuBTL8A2CgLdg8GA2LqUnL
GbTccols89QEcA4=
-----END CERTIFICATE-----`),
			},
		},
		{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-two",
				Namespace: "foo",
			},
			Type: "kubernetes.io/tls",
			Data: map[string][]byte{
				"tls.key": []byte("YWJjMTIz"),
				"tls.crt": []byte("ZGVmNDU2"),
			},
		},
		{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-three",
				Namespace: "foo",
			},
			Type: "kubernetes.io/tls",
			Data: map[string][]byte{
				"tls.key": []byte("abc123"),
				"tls.crt": []byte("def456"),
			},
		},
		{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-four",
				Namespace: "bar",
			},
			Type: "kubernetes.io/tls",
			Data: map[string][]byte{
				"tls.key": []byte("abc123"),
				"tls.crt": []byte("def456"),
			},
		},
	}

}
