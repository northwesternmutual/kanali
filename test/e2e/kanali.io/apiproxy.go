package kanaliio

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/frankgreco/tester/pkg/apis"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/framework"
	"github.com/northwesternmutual/kanali/test/e2e/tester"
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
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(framework.TestContext.KanaliEndpoint).WithPath("/foo").NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(errors.ErrorProxyNotFound).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to a backend Kubernetes service", func() {
		By("deploying an upstream application")
		_, _, err := tester.Deploy(f.BaseName, f.ClientSet,
			tester.WithServer(tester.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.BaseName).Create(
			builder.NewApiProxy("endpoint", f.BaseName).WithSourcePath("/endpoint").WithTargetPath("/").WithTargetBackendStaticService(tester.Name, tester.InsecurePort).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(framework.TestContext.KanaliEndpoint).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary http endpoint", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.ClientSet,
			tester.WithServer(tester.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.BaseName).Create(
			builder.NewApiProxy("endpoint", f.BaseName).WithSourcePath("/endpoint").WithTargetPath("/").WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(framework.TestContext.KanaliEndpoint).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary https endpoint using tls", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.ClientSet,
			tester.WithServer(tester.TLSTypePresent),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.BaseName).Create(
			builder.NewApiProxy("endpoint", f.BaseName).WithSourcePath("/endpoint").WithTargetPath("/").WithSecret(tester.Name).WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(framework.TestContext.KanaliEndpoint).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

	It("should proxy to arbitrary https endpoint using mutual tls", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.ClientSet,
			tester.WithServer(tester.TLSTypeMutual),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.BaseName).Create(
			builder.NewApiProxy("endpoint", f.BaseName).WithSourcePath("/endpoint").WithTargetPath("/").WithSecret(tester.Name).WithTargetBackendEndpoint(dns).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(framework.TestContext.KanaliEndpoint).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})

})
