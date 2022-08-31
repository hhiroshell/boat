package daemon

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("func Kubeconfig()", func() {
		When("the kube-boat daemon returns correct kubeconfig response", func() {
			var response = `{
    "server": "https://localhost:50000/",
    "client-certificate-data": "Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK",
    "client-key-data": "Y2xpZW50LWtleS1kYXRhCg=="
}`
			It("should return a correctly unmarshalled result", func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					"http://localhost/kube-boat",
					httpmock.NewStringResponder(http.StatusOK, response),
				)

				client := &Client{client: &http.Client{}}
				kubeconfig, err := client.Kubeconfig()

				Expect(err).NotTo(HaveOccurred())
				Expect(kubeconfig.Server).To(Equal("https://localhost:50000/"))
				Expect(kubeconfig.ClientCert).To(Equal("Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK"))
				Expect(kubeconfig.ClientKey).To(Equal("Y2xpZW50LWtleS1kYXRhCg=="))
			})
		})
	})

	Describe("func StopDaemon()", func() {
		When("the kube-boat daemon returns accepting response", func() {
			var response = `{
    "message": "shutting down the server..."
}`
			It("should return a correctly unmarshalled result", func() {
				httpmock.RegisterResponder(
					http.MethodDelete,
					"http://localhost/kube-boat",
					httpmock.NewStringResponder(http.StatusAccepted, response),
				)

				client := &Client{client: &http.Client{}}
				msg, err := client.StopDaemon()

				Expect(err).NotTo(HaveOccurred())
				Expect(msg).To(Equal("shutting down the server..."))
			})
		})
	})

})
