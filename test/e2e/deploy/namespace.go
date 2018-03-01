package deploy

import (
	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/e2e/utils"
)

func CreateNamespace(i kubernetes.Interface, name string) (*v1.Namespace, error) {
	var got *v1.Namespace
	err := wait.PollImmediate(utils.Poll, utils.NamespaceCreationTimeout, func() (bool, error) {
		var err error
		got, err = i.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	return got, err
}

func DestroyNamespace(i kubernetes.Interface, name string) error {
	if err := i.CoreV1().Namespaces().Delete(name, nil); err != nil {
		return err
	}

	return wait.PollImmediate(utils.Poll, utils.NamespaceCleanupTimeout, func() (bool, error) {
		if _, err := i.CoreV1().Namespaces().Get(name, metav1.GetOptions{}); err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, nil
		}
		return false, nil
	})
}
