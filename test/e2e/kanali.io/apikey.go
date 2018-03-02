package kanaliio

import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

  "github.com/northwesternmutual/kanali/test/e2e/framework"
  "github.com/northwesternmutual/kanali/pkg/errors"
  "github.com/northwesternmutual/kanali/test/builder"
  "github.com/northwesternmutual/kanali/test/e2e/context"
  testutils "github.com/northwesternmutual/kanali/test/utils"
)

var _ = Describe("ApiKey", func() {
  f := framework.NewDefaultFramework("apikey")

	It("should be forbidden if no binding attached", func() {
    By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("apikey", f.Namespace.GetName()).WithSourcePath("/apikey").WithTargetPath("/").WithTargetBackendEndpoint("http://foo.bar.com").WithPlugin("apiKey", "", nil).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

    By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

    By("verifying the result")
    ok, err := testutils.RepresentJSONifiedObject(errors.ErrorForbidden).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())
	})
})
