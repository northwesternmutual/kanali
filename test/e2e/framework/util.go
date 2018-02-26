package framework

import (
	"time"

	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	Poll                    = 2 * time.Second
	NamespaceCleanupTimeout = 15 * time.Minute
)

func LoadConfig() (*restclient.Config, error) {
	return clientcmd.BuildConfigFromFlags("", TestContext.KubeConfig)
}
