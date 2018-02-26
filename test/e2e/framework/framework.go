package framework

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
)

type Framework struct {
	BaseName        string
	ClientSet       clientset.Interface
	KanaliClientSet versioned.Interface
	HTTPClient      *http.Client
	KanaliEndpoint  string
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

	var got *v1.Namespace
	err := wait.PollImmediate(Poll, 30*time.Second, func() (bool, error) {
		var err error
		got, err = f.ClientSet.CoreV1().Namespaces().Create(&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: f.BaseName,
			},
		})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	Expect(err).NotTo(HaveOccurred())
}

func (f *Framework) AfterEach() {
	err := f.ClientSet.CoreV1().Namespaces().Delete(f.BaseName, &metav1.DeleteOptions{})
	Expect(err).NotTo(HaveOccurred())

	err = wait.PollImmediate(Poll, NamespaceCleanupTimeout, func() (bool, error) {
		if _, err := f.ClientSet.CoreV1().Namespaces().Get(f.BaseName, metav1.GetOptions{}); err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, nil
		}
		return false, nil
	})
	Expect(err).NotTo(HaveOccurred())
}
