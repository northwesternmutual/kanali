// Copyright (c) 2018 Northwestern Mutual.
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

package middleware

import (
	"net/http"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/log"
)

func Recover(next http.Handler) http.Handler {
	runtime.ReallyCrash = false
	runtime.PanicHandlers = nil

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithContext(r.Context())

		defer runtime.HandleCrash(func(err interface{}) {
			logger.Error("kanali gateway paniced")

			jsonErr, data := errors.ToJSON(nil)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(jsonErr.Status)
			if _, err := w.Write(data); err != nil {
				logger.Error(err.Error())
			}

		})

		next.ServeHTTP(w, r)
	})
}
