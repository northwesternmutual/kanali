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
