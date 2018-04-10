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

		if err := validate.New(r.Context(), i).IsValidResource(review.Request.Kind, review.Request.Object.Raw); err != nil {
			review.Response.Result = &metav1.Status{
				Status:  metav1.StatusSuccess,
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
