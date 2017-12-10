package crds

import (
	"errors"
	"fmt"
	"time"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func Create(i apiextensionsv1beta1client.ApiextensionsV1beta1Interface, crds ...*apiextensionsv1beta1.CustomResourceDefinition) error {
	for _, crd := range crds {
		_, err := i.CustomResourceDefinitions().Create(crd)
		if err != nil && !k8sErrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create CRD %s: %v", crd.ObjectMeta.Name, err)
		}

		err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
			crd, err = i.CustomResourceDefinitions().Get(crd.ObjectMeta.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}
			for _, cond := range crd.Status.Conditions {
				switch cond.Type {
				case apiextensionsv1beta1.Established:
					if cond.Status == apiextensionsv1beta1.ConditionTrue {
						return true, err
					}
				case apiextensionsv1beta1.NamesAccepted:
					if cond.Status == apiextensionsv1beta1.ConditionFalse {
						return false, errors.New(cond.Reason)
					}
				}
			}
			return false, err
		})
		if err != nil {
			deleteErr := i.CustomResourceDefinitions().Delete(crd.ObjectMeta.Name, nil)
			if deleteErr != nil {
				return deleteErr
			}
			return err
		}
	}

	return nil
}
