package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
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
