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

package e2e

import (
  "io"
  "bytes"
  "testing"
  "net/http"
  "io/ioutil"
  "encoding/json"

  "k8s.io/api/core/v1"

  "github.com/northwesternmutual/kanali/pkg/errors"
)

func RunE2ETests(t *testing.T) {

  tests := []struct{
    name string
    secret *v1.Secret
    request *http.Request
    expectedStatusCode int
    expectedBody []byte
    exepectedHeaders http.Header
    setup: func() error
    teardown: func() error
  }{
   {
     name: "simple",
     secret: nil,
     request: newRequestOrDie(t, "GET", "http://127.0.0.1:8080/foo", nil),
     expectedStatusCode: 404,
     expectedBody: marshalOrDie(t, errors.ErrorProxyNotFound),
     setup: none,
     teardown: none,
   },
   {
     name: "endpoint",
     secret: nil,
     request: newRequestOrDie(t, "GET", "http://127.0.0.1:8080/foo", nil),
     expectedStatusCode: 200,
     expectedBody: nil,
     setup: none,
     teardown: none,
   },
  }

  client := &http.Client{}
  for _, test := range tests {
    setup()
    defer teardown()
    resp, err := client.Do(test.request)
    if err != nil {
      t.Fail()
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Fail()
    }
    if !bytes.Equal(body, test.expectedBody) {
      t.Fail()
    }
      t.Fail()
    }
  }
}

func none() error {
  return nil
}

func marshalOrDie(t *testing.T, d interface{}) []byte {
  data, err := json.Marshal(d)
  if err != nil {
    t.Fail()
  }
  return data
}

func newRequestOrDie(t *testing.T, method, url string, body io.Reader) *http.Request {
  req, err := http.NewRequest(method, url, body)
  if err != nil {
    t.Fail()
  }
  return req
}