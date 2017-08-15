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

package steps

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/northwesternmutual/kanali/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestWriteResponseGetName(t *testing.T) {
	step := WriteResponseStep{}
	assert.Equal(t, step.GetName(), "Write Response", "step name is incorrect")
}

func TestWriteResponseDo(t *testing.T) {
	step := WriteResponseStep{}
	writer := httptest.NewRecorder()
	response := &httptest.ResponseRecorder{
		Code: 200,
		HeaderMap: http.Header{
			"One":   []string{"two"},
			"Three": []string{"four"},
		},
		Body: bytes.NewBuffer([]byte("this is my mock response body")),
	}
	err := step.Do(context.Background(), &metrics.Metrics{}, nil, writer, nil, response.Result(), opentracing.StartSpan("test span"))
	defer writer.Result().Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, writer.Result().StatusCode, 200)
	assert.Equal(t, writer.Result().Status, "OK")
	assert.Equal(t, len(writer.Result().Header), 2)
	assert.Equal(t, writer.Result().Header.Get("one"), "two")
	assert.Equal(t, writer.Result().Header.Get("three"), "four")
	bodyBytes, err := ioutil.ReadAll(writer.Result().Body)
	assert.Nil(t, err)
	assert.Equal(t, string(bodyBytes), "this is my mock response body")
}
