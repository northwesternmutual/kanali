package kanaliio

import (
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
