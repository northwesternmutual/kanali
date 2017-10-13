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

import (
	"bytes"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"k8s.io/kubernetes/pkg/api"
)

// ComputeTargetPath calcuates the target or destination path based on the incoming path,
// desired target path prefix and the assicated proxy
func ComputeTargetPath(proxyPath, proxyTarget, requestPath string) string {

	proxyPath = NormalizeURLPath(proxyPath)
	proxyTarget = NormalizeURLPath(proxyTarget)
	requestPath = NormalizeURLPath(requestPath)

	var buffer bytes.Buffer

	if len(strings.SplitAfter(requestPath, proxyPath)) == 0 {
		buffer.WriteString("/")
	} else if proxyTarget != "/" {
		buffer.WriteString(proxyTarget)
	}

	buffer.WriteString(strings.SplitAfter(requestPath, proxyPath)[1])

	if len(buffer.Bytes()) == 0 {
		return "/"
	}

	return buffer.String()
}

// GetAbsPath returns the absolute path given any path
// the returned path is in a form that Kanali prefers
func GetAbsPath(path string) (string, error) {

	p, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if p[len(p)-1] == '/' {
		if len(p) < 2 {
			return "", nil
		}
		return p[:len(p)-2], nil
	}

	return p, nil

}

// CompareObjectMeta will loosly determine whether two ObjectMeta objects are equal.
// It does this by comparing the name and namespace
func CompareObjectMeta(c1, c2 api.ObjectMeta) bool {
	return c1.Namespace == c2.Namespace && c1.Name == c2.Name
}

// ComputeURLPath will correct a URL path that might be valid but not ideally formatted
func ComputeURLPath(u *url.URL) string {
	return NormalizeURLPath(u.EscapedPath())
}

// NormalizeURLPath will normalize any string treating it like a URL path.
func NormalizeURLPath(path string) string {
	if len(path) < 1 {
		return "/"
	}

	path = regexp.MustCompile(`/{2,}`).ReplaceAllString(path, "/")

	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	if len(path) < 1 {
		return "/"
	}

	if path[0] != '/' {
		path = "/" + path
	}

	return path
}
