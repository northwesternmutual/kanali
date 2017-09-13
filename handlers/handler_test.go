package handlers

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	r1 := &http.Request{
		URL: &url.URL{
			Path: "///foo//bar/car",
		},
	}
	r2 := &http.Request{
		URL: &url.URL{
			Path: "foo//bar/car/",
		},
	}
	r3 := &http.Request{
		URL: &url.URL{
			Path: "",
		},
	}
	r4 := &http.Request{
		URL: &url.URL{
			Path: "////",
		},
	}
	normalize(r1)
	normalize(r2)
	normalize(r3)
	normalize(r4)

	assert.Equal(t, "/foo/bar/car", r1.URL.Path)
	assert.Equal(t, "/foo/bar/car", r2.URL.Path)
	assert.Equal(t, "/", r3.URL.Path)
	assert.Equal(t, "/", r4.URL.Path)

}
