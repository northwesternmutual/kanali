package context

import (
	"flag"
	"fmt"
)

type Provider string

var (
	Minikube Provider = "MINIKUBE"
	AWS      Provider = "AWS"
)

type TestContextType struct {
	KubeConfig string
	KanaliConfig
	ProjectConfig
	CloudConfig
}

type CloudConfig struct {
	Provider
}

type ProjectConfig struct {
	CommitSHA string
}

type KanaliConfig struct {
	Scheme   string
	Host     string
	Port     int
  ImageConfig
}

type ImageConfig struct {
  Name, Tag string
}

var TestContext TestContextType

func RegisterCommonFlags() {
	flag.Var(&TestContext.CloudConfig.Provider, "e2e.cloud.provider", "")
	flag.StringVar(&TestContext.KubeConfig, "e2e.kube.config", "", "")
	flag.StringVar(&TestContext.KanaliConfig.Host, "e2e.kanali.host", "", "")
	flag.StringVar(&TestContext.ProjectConfig.CommitSHA, "e2e.project.commit_sha", "", "")
  flag.StringVar(&TestContext.KanaliConfig.ImageConfig.Name, "e2e.kanali.image.name", "", "")
  flag.StringVar(&TestContext.KanaliConfig.ImageConfig.Tag, "e2e.kanali.image.tag", "", "")
}

func (cfg KanaliConfig) GetEndpoint() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Scheme, cfg.Host, cfg.Port)
}

// String implements the flag.Value interface
func (p *Provider) String() string {
	return string(*p)
}

// Set implements the flag.Value interface
func (p *Provider) Set(data string) error {
	switch data {
	case "minikube":
		*p = Minikube
	case "aws":
		*p = AWS
	default:
		return fmt.Errorf("%s is not of type context.Provider", data)
	}
	return nil
}
