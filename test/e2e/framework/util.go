package framework

import (
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/northwesternmutual/kanali/test/e2e/context"
)

func LoadConfig() (*restclient.Config, error) {
	return clientcmd.BuildConfigFromFlags("", context.TestContext.KubeConfig)
}
