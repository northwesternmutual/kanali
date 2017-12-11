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

package utils

// import (
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// )
//
// func TestComputeTargetPath(t *testing.T) {
// 	assert.Equal(t, "/", NormalizeURLPath(ComputeTargetPath("/foo/bar", "", "/foo/bar")))
// 	assert.Equal(t, "/", NormalizeURLPath(ComputeTargetPath("/foo/bar", "/", "/foo/bar")))
// 	assert.Equal(t, "/foo", NormalizeURLPath(ComputeTargetPath("/foo/bar", "/foo", "/foo/bar")))
// 	assert.Equal(t, "/foo/bar", NormalizeURLPath(ComputeTargetPath("/foo/bar", "/foo", "/foo/bar/bar")))
// 	assert.Equal(t, "/bar", NormalizeURLPath(ComputeTargetPath("/foo/bar", "", "/foo/bar/bar")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two", "/", "/api/v1/example-two/accounts/")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two", "", "/api/v1/example-two/accounts/")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two/", "/", "/api/v1/example-two/accounts/")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts/")))
// 	assert.Equal(t, "/accounts", NormalizeURLPath(ComputeTargetPath("/api/v1/example-two/", "", "/api/v1/example-two/accounts")))
// 	assert.Equal(t, "/", NormalizeURLPath(ComputeTargetPath("/", "", "/")))
// 	assert.Equal(t, "/", NormalizeURLPath(ComputeTargetPath("/", "/", "/")))
// }
//
// func TestAbsPath(t *testing.T) {
// 	p, _ := GetAbsPath("/")
// 	assert.Equal(t, "", p)
// 	p, _ = GetAbsPath("/foo/")
// 	assert.Equal(t, "/foo", p)
// 	p, _ = GetAbsPath("//")
// 	assert.Equal(t, "", p)
// }
//
// func TestNormalizeURLPath(t *testing.T) {
// 	assert.Equal(t, "/foo/bar", NormalizeURLPath("foo////bar"))
// 	assert.Equal(t, "/foo", NormalizeURLPath("foo"))
// 	assert.Equal(t, "/foo", NormalizeURLPath("foo////"))
// 	assert.Equal(t, "/foo/bar", NormalizeURLPath("///foo////bar//"))
// 	assert.Equal(t, "/", NormalizeURLPath(""))
// 	assert.Equal(t, "/", NormalizeURLPath("////"))
// 	assert.Equal(t, "/https%3A%2F%2Fgoogle.com", NormalizeURLPath("/////https%3A%2F%2Fgoogle.com"))
// }
