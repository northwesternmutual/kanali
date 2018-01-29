package tracer

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"

	"github.com/northwesternmutual/kanali/pkg/tags"
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
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestMethod], "GET")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestURLPath], "/")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestURLHost], "foo.bar.com")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestBody], "test data")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestHeaders], `{"Bar":["foo"],"Foo":["bar","car"]}`)
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPRequestURLQuery], `{"foo":["bar"]}`)

	testSpanTwo := mockTracer.StartSpan("test span two")
	HydrateSpanFromRequest(nil, testSpanTwo)
	testSpanTwo.Finish()
	assert.Nil(t, mockTracer.FinishedSpans()[1].Tags()[tags.HTTPRequest])
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
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPResponseBody], "test data")
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPResponseHeaders], `{"Bar":["foo"],"Foo":["bar","car"]}`)
	assert.Equal(t, mockTracer.FinishedSpans()[0].Tags()[tags.HTTPResponseStatusCode], 200)

	testSpanTwo := mockTracer.StartSpan("test span two")
	HydrateSpanFromResponse(nil, testSpanTwo)
	testSpanTwo.Finish()
	assert.Nil(t, mockTracer.FinishedSpans()[1].Tags()[tags.HTTPResponse])
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

func BenchmarkOmitHeaderValues(b *testing.B) {
	for n := 0; n < b.N; n++ {
		omitHeaderValues(http.Header{
			"One":   []string{"two"},
			"Three": []string{"four"},
		}, "omitted", "one")
	}
}

func BenchmarkHydrateSpanFromRequest(b *testing.B) {
  mockTracer := mocktracer.New()
  testReqOne, _ := http.NewRequest("GET", "https://foo.bar.com/?foo=bar", bytes.NewReader([]byte("test data")))
  testReqOne.Header.Add("foo", "bar")
  testSpanOne := mockTracer.StartSpan("test span one")

	for n := 0; n < b.N; n++ {
  	HydrateSpanFromRequest(testReqOne, testSpanOne)
	}
}
