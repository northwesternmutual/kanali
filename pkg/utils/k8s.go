package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesRestConfig(location string) (*rest.Config, error) {
	if len(location) > 0 {
		// user has specified a path to their own kubeconfig file so we'll use that
		return clientcmd.BuildConfigFromFlags("", location)
	}
	// use the in cluster config as the user has not specified their own
	return rest.InClusterConfig()
}
