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
