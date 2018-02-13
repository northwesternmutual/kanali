package kanali

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/northwesternmutual/kanali/test/builder"
	"github.com/northwesternmutual/kanali/test/e2e/framework"
	testutils "github.com/northwesternmutual/kanali/test/utils"
)

var _ = Describe("ApiProxy", func() {
	f := framework.NewDefaultFramework("apiproxy")

	It("should proxy to arbitrary http endpoint", func() {
		By("creating a server for an arbitrary upstream service")
		server := builder.NewHTTPServer(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Hello, client")
			}),
		).NewOrDie()
		server.Start()
		defer server.Close()

		By("creating an apiproxy")
		apiproxy := builder.NewApiProxy("endpoint", "default").WithSourcePath("/endpoint").WithTargetBackendEndpoint(server.URL).NewOrDie()
		apiproxy, err := f.KanaliClientSet.KanaliV2().ApiProxies("default").Create(apiproxy)
		Expect(err).NotTo(HaveOccurred())

		By("waiting for that apiproxy to be applied")
		// TODO

		By("preforming an http request")
		port, err := f.GetKanaliNodePort()
		Expect(err).NotTo(HaveOccurred())
		req := builder.NewHTTPRequest("GET").WithHTTP().WithHostPort("localhost", port).WithPath("/").NewOrDie()
		resp, err := f.HTTPClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		ok, err := testutils.RepresentJSONifiedObject(struct{ foo string }{foo: "bar"}).Match(resp)
		Expect(err).NotTo(HaveOccurred())
		Expect(ok).To(BeTrue())

		By("deleting the apiproxy")
		err = f.KanaliClientSet.KanaliV2().ApiProxies("default").Delete(apiproxy.GetName(), &v1.DeleteOptions{})
		Expect(err).NotTo(HaveOccurred())
	})

})
