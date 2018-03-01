package kanaliio

// import (
//   . "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
//
//   "github.com/northwesternmutual/kanali/pkg/rsa"
//   "github.com/northwesternmutual/kanali/test/e2e/deploy/kanali"
//   "github.com/northwesternmutual/kanali/test/e2e/framework"
// )
//
// var _ = Describe("ApiKey", func() {
//   f := framework.NewDefaultFramework("apikey")
//
// 	It("should do something", func() {
// 		By("doing this and that")
//     encryptedKey, err := rsa.Encrypt([]byte("abc123"), kanali.TestApiKeyDecryptionPublicKey,
//       rsa.Base64Encode(),
//       rsa.WithEncryptionLabel(rsa.EncryptionLabel),
//     )
//     Expect(err).NotTo(HaveOccurred())
// 	})
// })
