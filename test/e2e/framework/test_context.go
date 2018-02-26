package framework

import (
	"flag"
)

type TestContextType struct {
	KubeConfig     string
	KanaliEndpoint string
}

var TestContext TestContextType

func RegisterCommonFlags() {
	flag.StringVar(&TestContext.KubeConfig, "kubeconfig", "", "")
	flag.StringVar(&TestContext.KanaliEndpoint, "kanali-endpoint", "", "")
}
