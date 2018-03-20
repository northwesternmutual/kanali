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

const (
	moreDetails = "Visit https://kanali.io/docs/errorcodes for more details."
)

var (
	ErrorProxyNotFound                    = Error{http.StatusNotFound, "No ApiProxy resource was not found that matches the request.", 0, moreDetails}
	ErrorUnknown                          = Error{http.StatusInternalServerError, "An unknown error occured.", 1, moreDetails}
	ErrorMockTargetNotFound               = Error{http.StatusNotFound, "No MockTarget resource was not found that matches the request.", 2, moreDetails}
	ErrorCouldNotLoadPlugin               = Error{http.StatusInternalServerError, "Could not open or load plugin.", 3, moreDetails}
	ErrorCouldNotLookupPluginSymbol       = Error{http.StatusInternalServerError, "Could not lookup plugin symbol.", 4, moreDetails}
	ErrorPluginIncorrectInterface         = Error{http.StatusInternalServerError, "Plugin does not implement the correct interface.", 5, moreDetails}
	ErrorKubernetesSecretError            = Error{http.StatusInternalServerError, "Could not retreive Kubernetes TLS secret.", 6, moreDetails}
	ErrorCreateKeyPair                    = Error{http.StatusInternalServerError, "Could not create x509 key pair.", 7, moreDetails}
	ErrorBadGateway                       = Error{http.StatusBadGateway, "Could not get a valid or any response from the upstream server.", 8, moreDetails}
	ErrorKubernetesServiceError           = Error{http.StatusInternalServerError, "Could not retreive Kubernetes services.", 9, moreDetails}
	ErrorNoMatchingServices               = Error{http.StatusInternalServerError, "Could not retreive Kubernetes services.", 9, moreDetails}
	ErrorPluginRuntimeError               = Error{http.StatusInternalServerError, "Plugin threw a runtime error.", 10, moreDetails}
	ErrorApiProxyBackendEndpointMalformed = Error{http.StatusInternalServerError, "The ApiProxy endpoint backend is malformed.", 11, moreDetails}
	ErrorForbidden                        = Error{http.StatusForbidden, "My lips are sealed.", 12, moreDetails}
	ErrorApiKeyUnauthorized               = Error{http.StatusUnauthorized, "Api key is not authorized.", 13, moreDetails}
	ErrorTooManyRequests                  = Error{http.StatusTooManyRequests, "The Api key you are using has exceeded its rate limit.", 14, moreDetails}
)
