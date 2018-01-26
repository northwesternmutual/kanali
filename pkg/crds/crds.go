package crds

import (
	"fmt"
	"math/rand"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	utilwait "k8s.io/apimachinery/pkg/util/wait"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsv1beta1client "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
)

var (
	// retries is the maximum number of attempts that a CRD will be given to establish
	retries = 5
	// interval is the amount of time that will pass between retries
	interval = 500 * time.Millisecond
	// factor is the scalar by which the prevous interval will be increased
	factor = 1.0
)

// EnsureCRDs is a utility function that will ensure that a list of CustomResourceDefinitions
// are not only created, but ready for use by a Kubernetes cluster. If a CustomResourceDefinition
// cannot be established within a reasonable amount of retries, an error will be returned.
func EnsureCRDs(
	i apiextensionsv1beta1client.ApiextensionsV1beta1Interface,
	crds ...*apiextensionsv1beta1.CustomResourceDefinition,
) error {
	funcs := make([]func() error, len(crds))
	for index, crd := range crds {
		crd := crd.DeepCopy() // We want to process every crd, not just the last one.
		funcs[index] = func() error {
			return ensureCRD(i, crd)
		}
	}
	return utilerrors.AggregateGoroutines(funcs...)
}

func ensureCRD(
	i apiextensionsv1beta1client.ApiextensionsV1beta1Interface,
	crd *apiextensionsv1beta1.CustomResourceDefinition,
) error {
	// Attempt to create the crd (even if it's already present).
	// If it is present, we still need to check its status.
	if _, err := i.CustomResourceDefinitions().Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}

	err := utilwait.ExponentialBackoff(utilwait.Backoff{
		Factor:   factor, // Even though we are using a factor of 1, ExponentialBackoff is preferred over PollImmediate as it provides jitter.
		Steps:    retries,
		Jitter:   rand.Float64(),
		Duration: interval,
	}, func() (bool, error) {
		// Attempt to retrieve the CRD that was either already present or just created.
		crd, err := i.CustomResourceDefinitions().Get(crd.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		// Test if CRD is established. If it is not, attempt
		// to find a state in which it will never be established
		// and fail fast.
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				// This is the state we are looking for.
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, nil
				}
			case apiextensionsv1beta1.NamesAccepted:
				// If we have reached this state, the CRD will never become
				// established
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					return false, fmt.Errorf("due to the naming conflict %s, the CustomResourceDefinition %s will never become established", cond.Reason, crd.GetName())
				}
			}
		}
		return false, nil
	})

	if err == utilwait.ErrWaitTimeout {
		return fmt.Errorf("the CustomResourceDefinition %s was not established within a reasonable amount of time", crd.GetName())
	}
	return err
}
