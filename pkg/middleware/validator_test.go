// Copyright (c) 2018 Northwestern Mutual.
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

package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned/fake"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/test/builder"
)

type errorReader struct{}

func (r *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("foo bar car")
}

func TestValidator(t *testing.T) {
	core, _ := observer.New(zap.NewAtomicLevelAt(zapcore.DebugLevel))
	defer log.SetLogger(zap.New(core)).Restore()

	apiproxy := builder.NewApiProxy("foo", "bar").WithSourcePath("/foo").NewOrDie()
	data, err := json.Marshal(apiproxy)
	assert.Nil(t, err)

	review := &v1beta1.AdmissionReview{
		Request: &v1beta1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{
				Group:   "kanali.io",
				Version: "v2",
				Kind:    "ApiProxy",
			},
			Object: runtime.RawExtension{
				Raw: data,
			},
		},
	}
	data, err = json.Marshal(review)
	assert.Nil(t, err)

	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(nil))
	rec := httptest.NewRecorder()
	clientset := fake.NewSimpleClientset()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Validator(clientset).ServeHTTP(w, r)
		assert.Equal(t, 500, rec.Code)
	}).ServeHTTP(rec, req)

	req, _ = http.NewRequest("GET", "/", new(errorReader))
	rec = httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Validator(clientset).ServeHTTP(w, r)
		assert.Equal(t, 500, rec.Code)
	}).ServeHTTP(rec, req)

	req, _ = http.NewRequest("GET", "/", bytes.NewBuffer(data))
	rec = httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Validator(clientset).ServeHTTP(w, r)
		assert.Equal(t, 200, rec.Code)
	}).ServeHTTP(rec, req)

	_, err = clientset.KanaliV2().ApiProxies("foo").Create(
		builder.NewApiProxy("bar", "foo").WithSourcePath("/foo").NewOrDie(),
	)
	assert.Nil(t, err)

	req, _ = http.NewRequest("GET", "/", bytes.NewBuffer(data))
	rec = httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Validator(clientset).ServeHTTP(w, r)
		assert.Equal(t, 200, rec.Code)
	}).ServeHTTP(rec, req)
}
