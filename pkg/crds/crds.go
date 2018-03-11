package crds

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// TODO: do this concurrently
func EnsureCRDs(
	i apiextensionsv1beta1client.ApiextensionsV1beta1Interface,
	crds []*apiextensionsv1beta1.CustomResourceDefinition,
	stopCh <-chan struct{},
) error {
	for _, crd := range crds {
		if _, err := i.CustomResourceDefinitions().Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}

		if err := WaitForEstablished(i, crd, nil); err != nil {
			return err
		}
	}
	return nil
}
