package validate

import (
	"encoding/json"
	"io"

	"k8s.io/api/admission/v1beta1"
)

func UnmarshalAdmissionReview(closer io.ReadCloser) (*v1beta1.AdmissionReview, error) {

	var ar *v1beta1.AdmissionReview

	if err := json.NewDecoder(closer).Decode(ar); err != nil {
		return nil, err
	}

	return ar, nil

}
