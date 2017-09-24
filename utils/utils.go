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
	"net/http"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/purell"
	"github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/api"
)

// ComputeTargetPath calcuates the target or destination path based on the incoming path,
// desired target path prefix and the assicated proxy
func ComputeTargetPath(proxyPath, proxyTarget, requestPath string) string {

	target := ""

	// normalize paths
	if proxyPath[len(proxyPath)-1] == '/' {
		proxyPath = proxyPath[:len(proxyPath)-1]
	}
	if requestPath[len(requestPath)-1] == '/' {
		requestPath = requestPath[:len(requestPath)-1]
	}

	if proxyTarget == "/" {

		if len(strings.SplitAfter(requestPath, proxyPath)) == 0 {
			target = "/"
		} else {
			target = strings.SplitAfter(requestPath, proxyPath)[1]
		}

	} else {

		if len(strings.SplitAfter(requestPath, proxyPath)) == 0 {
			target = "/"
		} else {
			target = proxyTarget + strings.SplitAfter(requestPath, proxyPath)[1]
		}

	}

	if target == "" {
		return "/"
	}

	return target

}

// IsValidHTTPMethod validates whether a given string is a valid
// HTTP method or not
func IsValidHTTPMethod(m string) bool {
	m = strings.ToUpper(m)

	return m == http.MethodGet ||
		m == http.MethodHead ||
		m == http.MethodPost ||
		m == http.MethodPut ||
		m == http.MethodPatch ||
		m == http.MethodDelete ||
		m == http.MethodConnect ||
		m == http.MethodOptions ||
		m == http.MethodTrace
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

// NormalizePath will correct a URL path that might be valid but no ideally formatted
func NormalizePath(path string) string {
	result, err := purell.NormalizeURLString(path, purell.FlagRemoveDotSegments|purell.FlagRemoveDuplicateSlashes|purell.FlagRemoveTrailingSlash)
	if err != nil {
		logrus.Errorf("error normalizing url path - using original url path: %s", err.Error())
		return removeDupLeadingSlashes(path)
	}
	return removeDupLeadingSlashes(result)
}

func removeDupLeadingSlashes(path string) string {
	if len(path) < 1 {
		return "/"
	}
	var buffer bytes.Buffer
	var i int
	buffer.WriteString("/")
	for i = 0; i < len(path); i++ {
		if path[i] == '/' {
			continue
		}
		break
	}
	buffer.WriteString(path[i:])
	return buffer.String()
}
