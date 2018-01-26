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

import "net/http"

type Error struct {
	Status  int `json:"status"`
	Message string `json:"message"`
	Code    int `json:"code"`
	Details string `json:"details"`
}

// Error satifies the error interface
func (e Error) Error() string {
	return e.Message
}

var (
	ErrorProxyNotFound              = Error{http.StatusNotFound, "No ApiProxy resource was not found that matches the request.", 0, "More details coming soon!"}
	ErrorUnknown                    = Error{http.StatusInternalServerError, "An unknown error occured.", 1, "More details coming soon!"}
	ErrorMockTargetNotFound         = Error{http.StatusNotFound, "No MockTarget resource was not found that matches the request.", 2, "More details coming soon!"}
	ErrorCouldNotLoadPlugin         = Error{http.StatusInternalServerError, "Could not open or load plugin.", 3, "More details coming soon!"}
	ErrorCouldNotLookupPluginSymbol = Error{http.StatusInternalServerError, "Could not lookup plugin symbol.", 4, "More details coming soon!"}
	ErrorPluginIncorrectInterface   = Error{http.StatusInternalServerError, "Plugin does not implement the correct interface.", 5, "More details coming soon!"}
	ErrorKubernetesSecretError      = Error{http.StatusInternalServerError, "Could not retreive Kubernetes TLS secret.", 6, "More details coming soon!"}
	ErrorCreateKeyPair              = Error{http.StatusInternalServerError, "Could not create x509 key pair.", 7, "More details coming soon!"}
	ErrorBadGateway = Error{http.StatusBadGateway, "Could not get a valid or any response from the upstream server.", 8, "More details coming soon!"}
	ErrorKubernetesServiceError     = Error{http.StatusInternalServerError, "Could not retreive Kubernetes services.", 9, "More details coming soon!"}
	ErrorNoMatchingServices         = Error{http.StatusInternalServerError, "Could not retreive Kubernetes services.", 9, "More details coming soon!"}
	ErrorPluginRuntimeError         = Error{http.StatusInternalServerError, "Plugin threw a runtime error.", 10, "More details coming soon!"}
  ErrorApiProxyBackendEndpointMalformed = Error{http.StatusInternalServerError, "The ApiProxy endpoint backend is malformed.", 11, "More details coming soon!"}
)
