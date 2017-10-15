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

package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	//"net/url"
	"testing"
	"time"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/spec"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

func TestServeHTTP(t *testing.T) {

	randomHTTPPort := random(40000, 49999)
	viper.SetDefault(config.FlagProxyEnableMockResponses.GetLong(), true)

	handler := Handler{InfluxController: nil, H: IncomingRequest}
	server := &http.Server{Addr: fmt.Sprintf("127.0.0.1:%d", randomHTTPPort), Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})}
	listener, _ := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", randomHTTPPort))
	go server.Serve(listener)
	defer server.Close()

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d", randomHTTPPort))
	assert.Nil(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), fmt.Sprintf("%s\n", `{"code":404,"msg":"proxy not found"}`))
	assert.Equal(t, resp.StatusCode, 404)

	testProxyOne := spec.APIProxy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Mock: &spec.Mock{
				ConfigMapName: "testConfigMapOne",
			},
			Service: spec.Service{
				Name:      "my-service",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	mockOne, _ := json.Marshal([]spec.Route{
		{
			Route:  "/foo",
			Code:   200,
			Method: "GET",
			Body:   "{\"foo\": \"bar\"}",
		},
    {
			Route:  "/https%3A%2F%2Fgoogle.com",
			Code:   200,
			Method: "GET",
			Body:   "{\"car\": \"bar\"}",
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
	spec.MockResponseStore.Set(testConfigMapOne)
	spec.ProxyStore.Set(testProxyOne)
	resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/accounts/foo", randomHTTPPort))
	assert.Nil(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, resp.StatusCode, 200)

  resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/accounts", randomHTTPPort) + "/https%3A%2F%2Fgoogle.com")
	assert.Nil(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, resp.StatusCode, 200)

  resp, err = http.Get(fmt.Sprintf("http://127.0.0.1:%d////api///v1////accounts", randomHTTPPort) + "/https%3A%2F%2Fgoogle.com")
	assert.Nil(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, resp.StatusCode, 200)
}

// func TestNormalize(t *testing.T) {
// 	r1 := &http.Request{
// 		URL: &url.URL{
// 			Path: "///foo//bar/car",
// 		},
// 	}
// 	r2 := &http.Request{
// 		URL: &url.URL{
// 			Path: "foo//bar/car/",
// 		},
// 	}
// 	r3 := &http.Request{
// 		URL: &url.URL{
// 			Path: "",
// 		},
// 	}
// 	r4 := &http.Request{
// 		URL: &url.URL{
// 			Path: "////",
// 		},
// 	}
// 	// normalize(r1)
// 	// normalize(r2)
// 	// normalize(r3)
// 	// normalize(r4)
//
// 	assert.Equal(t, "/foo/bar/car", r1.URL.Path)
// 	assert.Equal(t, "/foo/bar/car", r2.URL.Path)
// 	assert.Equal(t, "/", r3.URL.Path)
// 	assert.Equal(t, "/", r4.URL.Path)
//
// }

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
