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

package kanaliio

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/frankgreco/tester/pkg/apis"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/context"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
	"github.com/northwesternmutual/kanali/test/e2e/deploy/tester"
	"github.com/northwesternmutual/kanali/test/e2e/framework"
	testutils "github.com/northwesternmutual/kanali/test/utils"
)

var _ = Describe("ApiProxy", func() {
	f := framework.NewDefaultFramework("apiproxy")
	requestDetails := func(r *http.Request, p *v2.ApiProxy) apis.RequestDetails {
		return apis.RequestDetails{
			Method: r.Method,
			Path:   p.Spec.Target.Path,
			Query:  map[string]string{},
		}
	}

	It("should not match any apiproxy", func() {
		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath("/foo").NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(errors.ErrorProxyNotFound).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to a backend Kubernetes service", func() {
		By("deploying an upstream application")
		_, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("endpoint", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/").WithTargetBackendStaticService(tester.Name, tester.InsecurePort).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should properly handle encoded urls", func() {
		By("deploying an upstream application")
		_, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("endpoint", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetBackendStaticService(tester.Name, tester.InsecurePort).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(
			fmt.Sprintf("%s/%%47%%6f%%2f", apiproxy.Spec.Source.Path),
		).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(apis.RequestDetails{
			Method: req.Method,
			Path:   fmt.Sprintf("/%%47%%6f%%2f"),
			Query:  map[string]string{},
		}).Match(resp)

		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should honor vhost routing", func() {
		By("deploying an upstream application")
		_, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxyOne, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("foo", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/foo").WithSourceHost("foo.bar.com").WithTargetBackendStaticService(tester.Name, tester.InsecurePort).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating another apiproxy")
		apiproxyTwo, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("bar", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/bar").WithSourceHost("bar.foo.com").WithTargetBackendStaticService(tester.Name, tester.InsecurePort).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxyOne.Spec.Source.Path).NewOrDie()
		req.Host = "foo.bar.com"
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxyOne)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		By("preforming another http request")
		req = builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxyTwo.Spec.Source.Path).NewOrDie()
		req.Host = "bar.foo.com"
		resp, err = f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err = testutils.RepresentJSONifiedObject(requestDetails(req, apiproxyTwo)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary http endpoint", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("endpoint", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/").WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary https endpoint using tls", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypePresent),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("endpoint", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/").WithSecret(tester.Name).WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary https endpoint using mutual tls", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeMutual),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("endpoint", f.Namespace.GetName()).WithSourcePath("/endpoint").WithTargetPath("/").WithSecret(tester.Name).WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

})
