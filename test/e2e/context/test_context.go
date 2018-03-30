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
	Scheme string
	Host   string
	Port   int
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
