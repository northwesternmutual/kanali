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
	"crypto/x509"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/northwesternmutual/kanali/config"
	"github.com/northwesternmutual/kanali/metrics"
	"github.com/northwesternmutual/kanali/spec"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/api"
)

type mockHTTPClient struct{}

func (cli *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	responseRecorder := &httptest.ResponseRecorder{}
	mockTracer := mocktracer.New()

	_, err := mockTracer.Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(req.Header),
	)
	if err != nil {
		return nil, errors.New("error extracting headers")
	}

	if req.URL.Path == "/error" {
		return nil, errors.New("expected error")
	}

	responseRecorder.Code = http.StatusOK
	responseRecorder.Body = bytes.NewBuffer([]byte("response body"))
	return responseRecorder.Result(), nil
}

func TestProxyPassGetName(t *testing.T) {
	step := ProxyPassStep{}
	assert.Equal(t, step.GetName(), "Proxy Pass", "step name is incorrect")
}

func TestPreformTargetProxy(t *testing.T) {
	testMetrics := &metrics.Metrics{}
	mockTracer := mocktracer.New()
	testReqOne, _ := http.NewRequest("GET", "https://foo.bar.com/?foo=bar", bytes.NewReader([]byte("test data")))
	testReqTwo, _ := http.NewRequest("GET", "https://foo.bar.com/error", bytes.NewReader([]byte("test data")))

	testSpanOne := mockTracer.StartSpan("test span one")
	resp, err := preformTargetProxy(&mockHTTPClient{}, testReqOne, testMetrics, testSpanOne)
	testSpanOne.Finish()
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, 200)
	assert.Equal(t, (*testMetrics)[0].Name, "total_target_time")
	assert.False(t, (*testMetrics)[0].Index)

	testSpanTwo := mockTracer.StartSpan("test span two")
	_, err = preformTargetProxy(&mockHTTPClient{}, testReqTwo, testMetrics, testSpanTwo)
	testSpanTwo.Finish()
	assert.Equal(t, err.Error(), "expected error")
}

func TestCreateTargetRequest(t *testing.T) {
	originalReq, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts", nil)

	proxyOne := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	_, err := createTargetRequest(proxyOne, originalReq)
	assert.Equal(t, err.Error(), "no matching services")

	spec.ServiceStore.Set(spec.Service{
		Name:      "bar",
		Namespace: "foo",
		ClusterIP: "1.2.3.4",
		Port:      8080,
	})

	targetReq, _ := createTargetRequest(proxyOne, originalReq)
	assert.Equal(t, targetReq.URL, &url.URL{
		Scheme:     "http",
		Host:       "bar.foo.svc.cluster.local:8080",
		Path:       "/",
		RawPath:    "",
		ForceQuery: false,
	})
}

func TestCreateTargetClient(t *testing.T) {
	originalReq, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts", nil)

	proxyOne := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	duration, _ := time.ParseDuration("1m0s")
	viper.SetDefault(config.FlagProxyUpstreamTimeout.GetLong(), duration)
	cli, err := createTargetClient(proxyOne, originalReq)
	assert.Equal(t, cli.Timeout, viper.GetDuration(config.FlagProxyUpstreamTimeout.GetLong()))
	assert.Nil(t, err)
	assert.Nil(t, cli.Transport)
}

func TestConfigureTargetTLS(t *testing.T) {
	originalReq, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts", nil)

	proxyOne := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	transport, err := configureTargetTLS(proxyOne, originalReq)
	assert.Nil(t, err)
	assert.Nil(t, transport)

	proxyOne.Spec.SSL = spec.SSL{
		SecretName: "mysecretname",
	}

	testSecret := api.Secret{
		ObjectMeta: api.ObjectMeta{
			Name:      "mysecretname",
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
	}

	spec.SecretStore.Set(testSecret)
	cert, _ := spec.X509KeyPair(testSecret)

	transport, err = configureTargetTLS(proxyOne, originalReq)
	assert.Nil(t, err)
	assert.Equal(t, transport.TLSClientConfig.Certificates[0], *cert)
	assert.Equal(t, transport.TLSClientConfig.RootCAs, x509.NewCertPool())

	viper.SetDefault(config.FlagProxyTLSCommonNameValidation.GetLong(), false)
	defer viper.Reset()

	transport, err = configureTargetTLS(proxyOne, originalReq)
	assert.Nil(t, err)
	assert.True(t, transport.TLSClientConfig.InsecureSkipVerify)

}

func TestGetTargetURL(t *testing.T) {
	spec.ServiceStore.Set(spec.Service{
		Name:      "bar",
		Namespace: "foo",
		ClusterIP: "1.2.3.4",
		Port:      8080,
	})
	req, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts", nil)
	reqTwo, _ := http.NewRequest("GET", "http://foo.bar.com/api/v1/accounts"+"/https%3A%2F%2Fgoogle.com", nil)

	proxyOne := &spec.APIProxy{
		ObjectMeta: api.ObjectMeta{
			Name:      "exampleAPIProxyOne",
			Namespace: "foo",
		},
		Spec: spec.APIProxySpec{
			Path:   "/api/v1/accounts",
			Target: "/",
			Service: spec.Service{
				Name:      "bar",
				Namespace: "foo",
				Port:      8080,
			},
		},
	}

	urlOne, _ := getTargetURL(proxyOne, req)
	assert.Equal(t, *urlOne, url.URL{
		Scheme:     "http",
		Host:       "bar.foo.svc.cluster.local:8080",
		Path:       "/",
		RawPath:    "",
		ForceQuery: false,
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), true)

	urlOne, _ = getTargetURL(proxyOne, req)
	assert.Equal(t, *urlOne, url.URL{
		Scheme:     "http",
		Host:       "1.2.3.4:8080",
		Path:       "/",
		RawPath:    "",
		ForceQuery: false,
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), false)
	proxyOne.Spec.SSL = spec.SSL{
		SecretName: "mysecretname",
	}

	urlTwo, _ := getTargetURL(proxyOne, req)
	assert.Equal(t, *urlTwo, url.URL{
		Scheme:     "https",
		Host:       "bar.foo.svc.cluster.local:8080",
		Path:       "/",
		RawPath:    "",
		ForceQuery: false,
	})

	viper.SetDefault(config.FlagProxyEnableClusterIP.GetLong(), true)

	urlTwo, _ = getTargetURL(proxyOne, req)
	assert.Equal(t, *urlTwo, url.URL{
		Scheme:     "https",
		Host:       "1.2.3.4:8080",
		Path:       "/",
		RawPath:    "",
		ForceQuery: false,
	})

	urlThree, _ := getTargetURL(proxyOne, reqTwo)
	assert.Equal(t, *urlThree, url.URL{
		Scheme:     "https",
		Host:       "1.2.3.4:8080",
		Path:       "/https://google.com",
		RawPath:    "/https%3A%2F%2Fgoogle.com",
		ForceQuery: false,
	})
}
