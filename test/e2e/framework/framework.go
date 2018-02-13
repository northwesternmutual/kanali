package framework

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/northwesternmutual/kanali/pkg/client/clientset/versioned"
)

type Framework struct {
	BaseName        string
	ClientSet       clientset.Interface
	KanaliClientSet versioned.Interface
	HTTPClient      *http.Client
}

func NewDefaultFramework(name string) *Framework {
	f := &Framework{
		BaseName:  name,
		ClientSet: nil,
		HTTPClient: &http.Client{
			Timeout: 500 * time.Millisecond,
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
}

func (f *Framework) AfterEach() {}

func (f *Framework) GetKanaliNodePort() (string, error) {
	svc, err := f.ClientSet.Core().Services("default").Get("kanali", v1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, port := range svc.Spec.Ports {
		if port.NodePort != 0 {
			return strconv.Itoa(int(port.NodePort)), nil
		}
	}
	return "", errors.New("kanali is not exposed on a node port")
}
