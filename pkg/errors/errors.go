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

package errors

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Details string `json:"details"`
}

// Error satifies the error interface
func (e Error) Error() string {
	return e.Message
}

// String implements the stringer interface
func (e Error) String() string {
	return e.Message
}

func New(code, status int, msg string) Error {
	return Error{
		Status:  status,
		Code:    code,
		Message: msg,
		Details: fmt.Sprintf("Visit https://kanali.io/docs/v2/errorcodes/#%02d for more details.", code),
	}
}

var (
	ErrorProxyNotFound              = New(0, http.StatusNotFound, "No ApiProxy resource was not found that matches the request.")
	ErrorUnknown                    = New(1, http.StatusInternalServerError, "An unknown error occurred.")
	ErrorMockTargetNotFound         = New(2, http.StatusNotFound, "No MockTarget resource was not found that matches the request.")
	ErrorCouldNotLoadPlugin         = New(3, http.StatusInternalServerError, "Could not open or load plugin.")
	ErrorCouldNotLookupPluginSymbol = New(4, http.StatusInternalServerError, "Could not lookup plugin symbol.")
	ErrorPluginIncorrectInterface   = New(5, http.StatusInternalServerError, "Plugin does not implement the correct interface.")
	ErrorKubernetesSecretError      = New(6, http.StatusInternalServerError, "Could not retrieve Kubernetes secret.")
	ErrorCreateKeyPair              = New(7, http.StatusInternalServerError, "Could not create x509 key pair.")
	ErrorBadGateway                 = New(8, http.StatusBadGateway, "Could not get a valid or any response from the upstream server.")
	ErrorKubernetesServiceError     = New(9, http.StatusInternalServerError, "Could not retrieve Kubernetes services.")
	ErrorNoMatchingServices         = New(10, http.StatusInternalServerError, "Could not retrieve any matching Kubernetes services.")
	ErrorPluginRuntimeError         = New(11, http.StatusInternalServerError, "Plugin threw a runtime error.")
	ErrorForbidden                  = New(12, http.StatusForbidden, "My lips are sealed.")
	ErrorApiKeyUnauthorized         = New(13, http.StatusUnauthorized, "Api key is not authorized.")
	ErrorTooManyRequests            = New(14, http.StatusTooManyRequests, "The Api key you are using has exceeded its rate limit.")
)
