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

package v2

// import (
// 	"crypto/x509"
// 	"encoding/pem"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )
//
// func TestGetKeyStore(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	message := "store is not empty"
//
// 	store.Clear()
// 	assert.Equal(0, len(store.keyMap), message)
//
// 	v := KeyFactory{}
// 	var i interface{} = &v
// 	_, ok := i.(Store)
// 	assert.True(ok, "KeyFactory does not implement the Store interface")
// }
//
// func TestAPIKeySet(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
// 	message := "key received is not expected"
//
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	store.Set(keyList.Items[1])
// 	store.Set(keyList.Items[2])
// 	err := store.Set(APIProxy{})
// 	assert.Equal("grrr - you're only allowed add api keys to the api key store.... duh", err.Error(), "error expected")
// 	assert.Equal(keyList.Items[0], store.keyMap["iamencrypted1"], message)
// 	assert.Equal(keyList.Items[1], store.keyMap["iamencrypted2"], message)
// 	assert.Equal(keyList.Items[2], store.keyMap["iamencrypted3"], message)
// }
//
// func TestAPIKeyUpdate(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
// 	message := "key received is not expected"
//
// 	store.Clear()
// 	store.Update(keyList.Items[0], keyList.Items[0])
// 	store.Update(keyList.Items[1], keyList.Items[1])
// 	store.Update(keyList.Items[2], keyList.Items[2])
// 	err := store.Update(APIProxy{}, APIProxy{})
// 	assert.Equal("grrr - you're only allowed add api keys to the api key store.... duh", err.Error(), "error expected")
// 	assert.Equal(keyList.Items[0], store.keyMap["iamencrypted1"], message)
// 	assert.Equal(keyList.Items[1], store.keyMap["iamencrypted2"], message)
// 	assert.Equal(keyList.Items[2], store.keyMap["iamencrypted3"], message)
// }
//
// func TestAPIKeyClear(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
//
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	assert.Equal(1, len(store.keyMap), "store should have one key in it")
// 	store.Clear()
// 	assert.Equal(0, len(store.keyMap), "empty store should have no keys")
// }
//
// func TestAPIKeyIsEmpty(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
//
// 	store.Clear()
// 	assert.True(store.IsEmpty())
// 	store.Set(keyList.Items[0])
// 	assert.False(store.IsEmpty())
// 	store.Clear()
// 	assert.True(store.IsEmpty())
//
// 	store.Set(keyList.Items[0])
// 	assert.False(store.IsEmpty())
// 	store.Delete(keyList.Items[0])
// 	assert.True(store.IsEmpty())
// }
//
// func TestAPIKeyGet(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
//
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	store.Set(keyList.Items[1])
// 	store.Set(keyList.Items[2])
// 	_, err := store.Get("foo", "bar")
// 	assert.Equal("should only pass the name of the api key", err.Error(), "error expected")
// 	_, err = store.Get(5)
// 	assert.Equal("when retrieving a key, use the keys name", err.Error(), "error expected")
// 	key, _ := store.Get("")
// 	assert.Nil(key, "no key name was passed")
// 	key, _ = store.Get("jkl012")
// 	assert.Nil(key, "key should not exist")
// 	key, _ = store.Get("iamencrypted1")
// 	assert.Equal(keyList.Items[0], key, "keys should be equal")
// 	key, _ = store.Get("iamencrypted2")
// 	assert.Equal(keyList.Items[1], key, "keys should be equal")
// 	key, _ = store.Get("iamencrypted3")
// 	assert.Equal(keyList.Items[2], key, "keys should be equal")
// }
//
// func TestAPIKeyDelete(t *testing.T) {
// 	assert := assert.New(t)
// 	store := KeyStore
// 	keyList := getTestAPIKeyList()
//
// 	store.Clear()
// 	store.Set(keyList.Items[0])
// 	store.Set(keyList.Items[1])
// 	result, _ := store.Delete(nil)
// 	assert.Nil(result, "should return nil")
// 	result, _ = store.Delete(keyList.Items[2])
// 	assert.Nil(result, "should return nil")
// 	result, _ = store.Delete(keyList.Items[1])
// 	assert.Equal(keyList.Items[1], result, "deleted key should be returned")
// 	key, _ := store.Get("iamencrypted2")
// 	assert.Nil(key, "deleted key should no longer be in the store")
// 	assert.Equal(1, len(store.keyMap), "store should have a length of 1")
// 	store.Set(keyList.Items[1])
// 	key, _ = store.Get("iamencrypted2")
// 	assert.Equal(keyList.Items[1], key, "key should have been added back")
// 	assert.Equal(2, len(store.keyMap), "store should have a length of 2")
// 	store.Set(keyList.Items[2])
// 	result, _ = store.Delete(keyList.Items[2])
// 	assert.Equal(keyList.Items[2], result, "deleted key should be returned")
// }
//
// func TestDecrypt(t *testing.T) {
//
// 	assert := assert.New(t)
//
// 	// setup a private key
// 	rawPrivateKeyData := []byte(`-----BEGIN RSA PRIVATE KEY-----
// MIIEpAIBAAKCAQEAqgey9x7FDnREl2eHX80OoF5ZH1zM37Na7O+F6VY7jcQDRZEd
// rIWCuASeypCfcxpIFGauLsjZ1uWaDj8i78nU5kpN01YDwlidMu0m5wBJbQBwApWp
// 53BqPOjKyMrLKBBE3roct1Lya4gUnKtvsKsq1ZukmCL2mZR63LqUnGwPUV+DhR8e
// aCB+c2cRmUjgmuO+EcfBWgETOEWw//BGEnnKxtNlNj/T0qevvQWBvmLgcY2GV4t1
// qd82fTQ9dCRP8pGgEAYaCospaF0SIhwdkD+IMNBxswvNypFHe+ee10H40bab/HtL
// fgD+0DdAcyk30Xd5UowOL0G5I9Om4P9xDR0OJQIDAQABAoIBAQCNjRtQ4Czte54e
// 7fGlr/EdUW6gzYFCOu7XkhDJ0SCDRUvz/nvVxNCuoioQOZaFHLwlP9aC3HN+lGdM
// tQNA3aaAkht4dYNrqJo2a3iXl4WJWXPmsvJf2xaW3rbzsuCu2dx8EDUX6dYn74I5
// a//v9JAUhR7iCTmDYjlmyW4vS0b1VlMY1fTRn2pV1MobOXAxtL4AeKckCZlqvvPs
// vI7JtcZiUH5X6w2S069hQjcnSOBaeNm2h96sElywXLK7ocn8WLfiK87KhYb+iwV3
// y4nLc4FWn1x84TeYUDr7KPKeNmJN9qUvHiyjq3cZeBh/kS3RKyJnmqIfxGOEjweA
// 81pMbfWBAoGBANy6Z2UqU5d8P1ENwruol3NVBr1hsjuLlPqyEr1RbRygXxhKo51Z
// 1YlvuLElrLI+EBXG5sZfoNqCNzrRjE2I/x6yD4YXTzQAbIIyuvepm02GTZkwCbeV
// ogtnUgZs/+NeK5SPxL1SHnEC2CKYzniqz5X4wU89Sy1IU0fVs3cWmWr7AoGBAMUz
// U8/gw8MhgDBUQN+GxRYWIpM5sYAxmxjlSTXYACSAvDe10EmzX0bvhqgpml4VCsh6
// kHcxFtRintqU+/sQX1xOB8PkbKyQ2F6sfARPdN96wt5illXmOJDcHFzF9JsPuokf
// 0oCmSqFMUwVEnOYHKHjPEgUXBMt6RQG6ZFKM7CFfAoGAEGN+vdIwtrSw9tGwvcI8
// zIFpqLISS72HaD3Eim37Tf7SULlqAIrlAflbmzuS4QEerJ9CD0SWauL0OtfINRGQ
// 77SE1o3wgYnuhJRDsZPURiNMq60Xa5welf4r5vI85aLeVX9/OhdYz0dIyPjTi74S
// /62inMGOTWRrWpARbSwuY/UCgYAp849i69f5jQjuPx5F7y/CURct5toiAWuTUkPO
// 1pLBMTnZSC51X1lnh9FHuQID/cn7eEd/c8v9hrVwnr00XflLc7YnMFMGKxI8ockI
// Sbb9S+pk0OhWgdGKz6ig++PbkO8H+OMZR0cdJVWMC1qtfNxZRwP4XTtB5lMD2Rk6
// jow1VQKBgQDUj8sLDAVS3C64U5BCitlnQ4hdScpBq6P09B624A3EjseYwCWK1h7L
// xpyQYWffQdGb6IoawFPWy7nT42d4IER+JmRJRKm2a7XD3LsNud5tBCB57ksN27v7
// I0teIby/h2MN7yPNjFgIq64JneHpKwEH04+qy/dWQafxrpBlf9tfwg==
// -----END RSA PRIVATE KEY-----`)
//
// 	block, _ := pem.Decode(rawPrivateKeyData)
// 	// parse the pem block into a private key
// 	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		assert.Fail(err.Error())
// 	}
//
// 	APIKeyDecryptionKey = privateKey
//
// 	apiKey := &APIKey{
// 		TypeMeta: metav1.TypeMeta{},
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "my-api-key",
// 			Namespace: "default",
// 		},
// 		Spec: APIKeySpec{
// 			APIKeyData: "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba",
// 		},
// 	}
//
// 	if err := apiKey.Decrypt(); err != nil {
// 		assert.Fail(err.Error())
// 	}
//
// 	assert.Equal("i3CZlcRnDhJZeZfkDw9BgeEtZuFQKiw9", apiKey.Spec.APIKeyData)
//
// 	apiKey.Spec.APIKeyData = ":=-0"
//
// 	err = apiKey.Decrypt()
// 	assert.NotNil(err)
//
// 	apiKey.Spec.APIKeyData = "9210f613f32a54eca4601d199b81dda5a4f93c0540ee6a8b9634c2d4976b13399a03276820cd85a35b625a96ffdeffa2e094f1349e1ed7510afd7f0f904595f0f1bd8707170a46e6d366395456568323e4de71973977d872ab9aa733b35fbdeec279fc1f4bc147e242f414652bae8d46b7c53af76a1c37254096e4e0aa89dfdf86d599692ab74849bfedd7ecc6b4409b01d1e4d989cdd9ca6db7c1a90cd86086da7508f85186d938ab2922e862832eb07281e5934d417addaba0ddc43f57f3613ab0aff4f353fdadc1116f9dca10338562a842904eb7b3ab77b6f919ac244a8b8fa4d2634ac2f9bec60ee4631894e6b823dd200dc0c793f5d1dfc08b749b2bba"
//
// 	rawPrivateKeyData = []byte(`-----BEGIN RSA PRIVATE KEY-----
// MIIEowIBAAKCAQEAwJg6donVve47GNVk8IT+uCMwwdGFJOKcvTssDDl5kOO0nNE6
// sQbKt+UDXvFg7hXQTnBmGlQp7ndQOiSJnmsLMzD35NEQCSuigu2hUTSgoiCslYXd
// glGSpIdIt4wu4ynZhXAoRnQvsVBm3me9q+OojyfNzlUM+BnRyUBYi63024moWSFT
// Oo0l4zPA6hHAznm+u48X7/0RKZAqCwMnhHotyt8OOuKH2uTDY7PZMECcZHPBYQcY
// 1u+9UJagdw6S4ok14FycOfr++ImNtonezT+P2XKJntCT03P+qSeLqMNbC8WxRb5Y
// tpHFaWNgeHJ2LFpYsAET4czoqodzq7i/JnLqdQIDAQABAoIBACq5Gi71gI6zbDSk
// EdIxDng2hjhYjBqUkoR/vdOwZEE54fTvjP98LDgC+25ySThHaoAb+upYRTz7Brb/
// J/AWetOOR09g+RevsEsu/5wN+pw8d0xr96lTAx6wS/c7h6Mow3IQYT7Pw0yoM0me
// 5bpdtCSdhdABmfDneVwVhA7oNesFCeLN9J6UKCpV0Q6OS08KmM44feNCvBIOBe1x
// FmxYuZOYLyB1kjl4la9DGmGSbgnuJIztaYYIF3QWL9hWr7qlk5PFJgjfWTkMY7np
// xlbZOmhi7j4F5eFe7ce3e3ThPAjhS3e8A9XtH0qX62WBBJ0HEP+sWOCblr+Ce1jB
// G0ONEMECgYEA41Q9NXBiQyaMj2gS5k8bb70KYuMYUj6m6Ghy/r7ughJ49OPW5qsZ
// rrJQDdKXVSM3damSOEcK4UzCdiiZ9lqLFX9GnB2X6GZ3qaZiONaz8TkklE4Bomcp
// q46FGjDZ4Un82lKrCiszVH9tYfK5kqf5mb+pBxV6otPnkAVs1AhI+20CgYEA2OKF
// Y9Q56SNQZEAdLbq/aE01N9rEbWi+XHIDyFHDfiR8u0Ljbw2VEjr5DgWRDh/zht02
// AUoZPNyIIXUdliIefOiXmzX3Wr7oxxvJHlzTg0S/zqEjZcvbnTEC0Hr2vJlLqyfF
// TnCWGW/xu6dHGN8j9YOIHOlDLPSAxtHIafYvfikCgYAgsU+wLL5k69lREmWNa5kF
// g8lHmWM5CLuWNtc63ZiNajys80tHlxm1bv1mT3/mBY+HQ2C86TKeXhylQ4eIz1Kr
// P8eW0gazrHDPHBeKFw5/xXkGPwWrJmzuuCtaLeZXqP1NJJBmgPn7z+uaJ38aoT+b
// Jd4YH7WoTxNBOhaBk8NzLQKBgDqmz4HWs657kVL7D6j9NqQDNSActkAMMmaXorQ6
// gE6NmXmethijKKwcbJvBH6AdPeM70U359udBkELUBCNEqNAIqm8b4f+VXyfxOKtQ
// WRbHscM0SnlN7t8MkQFGL5GDNzZa8/2pbr4ESu//ZbexKG1Nm7KN1k/C80xREQyu
// uds5AoGBAJuJeZhTz7eDGVTL/w+D+LcJ8XrPKFsSe4cexbctQ/CsAn2apObxyMA/
// vV+J+U6XO9amQH/QNUNg/JJObJhQjeGvunZk3FkMXzAh+lUC264d28UYiMTwP/Md
// 355Aa3eUhvyldATlYxSBiZz03j8vDh+hBd4Qaxp9hXYW/HI8r7lP
// -----END RSA PRIVATE KEY-----`)
//
// 	block, _ = pem.Decode(rawPrivateKeyData)
// 	// parse the pem block into a private key
// 	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		assert.Fail(err.Error())
// 	}
// 	APIKeyDecryptionKey = privateKey
//
// 	err = apiKey.Decrypt()
// 	assert.NotNil(err)
//
// }
//
// func getTestAPIKeyList() *APIKeyList {
//
// 	return &APIKeyList{
// 		ListMeta: metav1.ListMeta{},
// 		Items: []APIKey{
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "abc123",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeySpec{
// 					APIKeyData: "iamencrypted1",
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "def456",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeySpec{
// 					APIKeyData: "iamencrypted2",
// 				},
// 			},
// 			{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Name:      "ghi789",
// 					Namespace: "foo",
// 				},
// 				Spec: APIKeySpec{
// 					APIKeyData: "iamencrypted3",
// 				},
// 			},
// 		},
// 	}
//
// }
