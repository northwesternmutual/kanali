package tracer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
)

func TestHydrateSpanFromRequest(t *testing.T) {
	mockTracer := mocktracer.New()
	testReqOne, _ := http.NewRequest("GET", "https://foo.bar.com/?foo=bar", bytes.NewReader([]byte("test data")))
	testReqOne.Header.Add("foo", "bar")
	testReqOne.Header.Add("foo", "car")
	testReqOne.Header.Add("bar", "foo")

	testSpanOne := mockTracer.StartSpan("test span one")
	HydrateSpanFromRequest(testReqOne, testSpanOne)
	testSpanOne.Finish()
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestMethod], "GET")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestURLPath], "/")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestURLHost], "foo.bar.com")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestBody], "test data")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestHeaders], `{"Bar":["foo"],"Foo":["bar","car"]}`)
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPRequestURLQuery], `{"foo":["bar"]}`)

	testSpanTwo := mockTracer.StartSpan("test span two")
	HydrateSpanFromRequest(nil, testSpanTwo)
	testSpanTwo.Finish()
	assert.Nil(t, mockTracer.FinishedSpans()[1].Tags()[HTTPRequest])
}

func TestHydrateSpanFromResponse(t *testing.T) {
	mockTracer := mocktracer.New()
	responseRecorder := &httptest.ResponseRecorder{
		Code: 200,
		Body: bytes.NewBuffer([]byte("test data")),
		HeaderMap: http.Header{
			"Foo": []string{"bar", "car"},
			"Bar": []string{"foo"},
		},
	}
	mockResponseOne := responseRecorder.Result()

	testSpanOne := mockTracer.StartSpan("test span one")
	HydrateSpanFromResponse(mockResponseOne, testSpanOne)
	testSpanOne.Finish()
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPResponseBody], "test data")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPResponseHeaders], `{"Bar":["foo"],"Foo":["bar","car"]}`)
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[HTTPResponseStatusCode], 200)

	testSpanTwo := mockTracer.StartSpan("test span two")
	HydrateSpanFromResponse(nil, testSpanTwo)
	testSpanTwo.Finish()
	assert.Nil(t, mockTracer.FinishedSpans()[1].Tags()[HTTPResponse])
}

func TestDupReader(t *testing.T) {
	closer := ioutil.NopCloser(bytes.NewReader([]byte("test string")))
	closerOne, closerTwo, _ := dupReader(closer)
	data1, _ := ioutil.ReadAll(closerOne)
	assert.Equal(t, string(data1), "test string")
	data2, _ := ioutil.ReadAll(closerTwo)
	assert.Equal(t, string(data2), "test string")
}

func TestOmitHeaderValues(t *testing.T) {
	h := http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}
	copy := omitHeaderValues(h, "omitted", "one")
	assert.Equal(t, h, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")
	assert.Equal(t, copy, http.Header{
		"One":   []string{"omitted"},
		"Three": []string{"four"},
	}, "map should be equal")
	copy = omitHeaderValues(h, "omitted", "one", "foo", "bar")
	assert.Equal(t, copy, http.Header{
		"One":   []string{"omitted"},
		"Three": []string{"four"},
	}, "map should be equal")
	copy = omitHeaderValues(h, "omitted")
	assert.Equal(t, copy, http.Header{
		"One":   []string{"two"},
		"Three": []string{"four"},
	}, "original map should not change")
	copy = omitHeaderValues(nil, "omitted")
	assert.Equal(t, copy, http.Header{}, "map should be equal")
}
