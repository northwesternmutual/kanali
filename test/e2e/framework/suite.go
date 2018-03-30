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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/test/e2e/context"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
	"github.com/northwesternmutual/kanali/test/e2e/deploy/kanali"
)

type Suite struct {
	BaseName  string
	ClientSet kubernetes.Interface
	Namespace *v1.Namespace
}

func NewDefaultSuite(name string) *Suite {
	s := &Suite{
		BaseName:  name,
		ClientSet: nil,
	}

	BeforeSuite(s.BeforeSuite)
	AfterSuite(s.AfterSuite)

	return s
}

func (s *Suite) BeforeSuite() {
	if s.ClientSet == nil {
		By("creating a kubernetes clientset")
		config, err := LoadConfig()
		Expect(err).NotTo(HaveOccurred())
		s.ClientSet, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	}

	ns, err := deploy.CreateNamespace(s.ClientSet, "e2e-"+context.TestContext.CommitSHA)
	Expect(err).NotTo(HaveOccurred())
	s.Namespace = ns
	err = kanali.Deploy(s.ClientSet, ns.GetName(),
		deploy.WithServer(deploy.TLSTypeNone),
	)
	Expect(err).NotTo(HaveOccurred())
}

func (s *Suite) AfterSuite() {
	err := deploy.DestroyNamespace(s.ClientSet, s.Namespace.GetName())
	Expect(err).NotTo(HaveOccurred())
	err = kanali.Destroy(s.ClientSet)
	Expect(err).NotTo(HaveOccurred())
}
