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
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/frankgreco/tester/pkg/apis"

	"github.com/northwesternmutual/kanali/pkg/apis/kanali.io/v2"
	"github.com/northwesternmutual/kanali/pkg/errors"
	"github.com/northwesternmutual/kanali/pkg/rsa"
	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/context"
	"github.com/northwesternmutual/kanali/test/e2e/deploy"
	"github.com/northwesternmutual/kanali/test/e2e/deploy/kanali"
	"github.com/northwesternmutual/kanali/test/e2e/deploy/tester"
	"github.com/northwesternmutual/kanali/test/e2e/framework"
	testutils "github.com/northwesternmutual/kanali/test/utils"
)

var _ = Describe("ApiKey", func() {
	f := framework.NewDefaultFramework("apikey")
	requestDetails := func(r *http.Request, p *v2.ApiProxy) apis.RequestDetails {
		return apis.RequestDetails{
			Method: r.Method,
			Path:   p.Spec.Target.Path,
			Query:  map[string]string{},
		}
	}

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

	It("should have access", func() {
		By("deploying an upstream application")
		dns, _, err := tester.Deploy(f.BaseName, f.Namespace.GetName(), f.ClientSet,
			deploy.WithServer(deploy.TLSTypeNone),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apiproxy")
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies(f.Namespace.GetName()).Create(
			builder.NewApiProxy("apikey", f.Namespace.GetName()).WithSourcePath("/apikey").WithTargetPath("/").WithTargetBackendEndpoint(dns).WithPlugin("apiKey", "", map[string]string{
				"apiKeyBindingName": "apikey",
			}).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apikey")
		encryptedKey, err := rsa.Encrypt([]byte("abc123"), kanali.TestApiKeyDecryptionPublicKey, rsa.Base64Encode(), rsa.WithEncryptionLabel(rsa.EncryptionLabel))
		Expect(err).NotTo(HaveOccurred())
		apikey, err := f.KanaliClientSet.KanaliV2().ApiKeys().Create(
			builder.NewApiKey("john").WithRevision(v2.RevisionStatusActive, encryptedKey).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("creating an apikeybinding")
		_, err = f.KanaliClientSet.KanaliV2().ApiKeyBindings(f.Namespace.GetName()).Create(
			builder.NewApiKeyBinding("apikey", f.Namespace.GetName()).WithKeys(
				builder.NewKeyAccess(apikey.GetName()).WithDefaultRule(
					builder.NewRule().WithGlobal().NewOrDie(),
				).NewOrDie(),
			).NewOrDie(),
		)
		Expect(err).NotTo(HaveOccurred())

		By("preforming an http request")
		req := builder.NewHTTPRequest().WithMethod("GET").WithHeader("apikey", "abc123").WithHost(context.TestContext.KanaliConfig.GetEndpoint()).WithPath(apiproxy.Spec.Source.Path).NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())

		By("verifying the result")
		ok, err := testutils.RepresentJSONifiedObject(requestDetails(req, apiproxy)).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		By("deleting an apikey")
		err = f.KanaliClientSet.KanaliV2().ApiKeys().Delete(apikey.GetName(), nil)
		Expect(err).NotTo(HaveOccurred())
	})
})
