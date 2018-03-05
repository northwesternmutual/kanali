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

package framework

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
)

type Framework struct {
	BaseName        string
	ClientSet       clientset.Interface
	KanaliClientSet versioned.Interface
	HTTPClient      *http.Client
	Namespace       *v1.Namespace
}

func NewDefaultFramework(name string) *Framework {
	f := &Framework{
		BaseName:  name,
		ClientSet: nil,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	BeforeEach(f.BeforeEach)
	AfterEach(f.AfterEach)

	return f
}

// BeforeEach gets a client and makes a namespace.
func (f *Framework) BeforeEach() {
	if f.ClientSet == nil {
		By("creating a kubernetes clientset")
		config, err := LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		f.ClientSet, err = clientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
		By("creating a kanali clientset")
		f.KanaliClientSet, err = versioned.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	}

	ns, err := deploy.CreateNamespace(f.ClientSet, fmt.Sprintf("e2e-%s-%d", f.BaseName, rand.Intn(1000)))
	Expect(err).NotTo(HaveOccurred())
	f.Namespace = ns
}

func (f *Framework) AfterEach() {
	err := deploy.DestroyNamespace(f.ClientSet, f.Namespace.GetName())
	Expect(err).NotTo(HaveOccurred())
}
