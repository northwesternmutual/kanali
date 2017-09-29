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
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	server := &http.Server{Addr: "127.0.0.1:40123", Handler: Logger(Handler{InfluxController: nil, H: IncomingRequest})}
	listener, _ := net.Listen("tcp4", "127.0.0.1:40123")
	go server.Serve(listener)
	defer server.Close()

	writer := new(bytes.Buffer)
	logrus.SetOutput(writer)
	resp, err := http.Get("http://127.0.0.1:40123/")
	logrus.SetOutput(os.Stdout)
	assert.Nil(t, err)
	assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, string(body), fmt.Sprintf("%s\n", `{"code":404,"msg":"proxy not found"}`))
	assert.Equal(t, resp.StatusCode, 404)

	logOutput := writer.String()
	assert.True(t, strings.Contains(logOutput, `msg="proxy not found"`))
	assert.True(t, strings.Contains(logOutput, `msg="request details"`))
	assert.True(t, strings.Contains(logOutput, `level=error`))
	assert.True(t, strings.Contains(logOutput, `client ip=127.0.0.1`))
	assert.True(t, strings.Contains(logOutput, `method=GET`))
	assert.True(t, strings.Contains(logOutput, `uri="/"`))
}
