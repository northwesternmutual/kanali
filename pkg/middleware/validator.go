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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/pkg/log"
	"github.com/northwesternmutual/kanali/pkg/validate"
)

func Validator(i versioned.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.WithContext(r.Context())

		data, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			logger.Error(fmt.Sprintf("error reading request body: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		review := new(v1beta1.AdmissionReview)

		if err := json.Unmarshal(data, review); err != nil {
			logger.Error(fmt.Sprintf("error unmarshaling request body into v1beta1.AdmissionReview: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		review.Response = new(v1beta1.AdmissionResponse)

		if err := validate.New(r.Context(), i).IsValidResource(review.Request.Kind, review.Request.Object.Raw); err != nil {
			review.Response.Result = &metav1.Status{
				Status:  metav1.StatusFailure,
				Message: err.Error(),
			}
			logger.Info(review.Response.Result.String())
		} else {
			logger.Debug(fmt.Sprintf("%s is valid", review.Request.Kind.String()))
			review.Response.Allowed = true
		}

		data, err = json.Marshal(review)
		if err != nil {
			logger.Error(fmt.Sprintf("error marshaling v1beta1.AdmissionReview: %s", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if n, err := w.Write(data); err != nil {
			logger.Error(fmt.Sprintf("error writing request to response - wrote %d bytes: %s", n, err))
		}

	}
}
