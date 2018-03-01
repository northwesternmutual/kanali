package deploy

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/e2e/utils"
)

type config interface {
	SetServerType(TLSType)
}

func Pod(i kubernetes.Interface, po *v1.Pod) (*v1.Pod, error) {
	pod, err := i.CoreV1().Pods(po.GetNamespace()).Create(po)
	if err != nil {
		return nil, err
	}

	return pod, wait.Poll(utils.Poll, time.Second*30, func() (bool, error) {
		ep, err := i.CoreV1().Endpoints(po.GetNamespace()).Get(po.GetName(), metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if ep.Subsets == nil {
			return false, nil
		}
		for _, ss := range ep.Subsets {
			if ss.Addresses != nil && len(ss.Addresses) > 0 {
				return true, nil
			}
		}
		return false, nil
	})
}
