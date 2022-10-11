package daemon

import (
	"net/http"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	BeforeEach(func() {
		httpmock.Reset()
	})

	Describe("func Readyz()", func() {
		When("the kube-boat daemon returns 200 OK", func() {
			It("ends with no error", func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					"http://localhost"+readyz,
					httpmock.NewStringResponder(http.StatusOK, "ok"),
				)

				client := &Client{client: &http.Client{}}
				err := client.Readyz()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		When("the kube-boat daemon is not running", func() {
			It("ends up with error", func() {
				// do not register responder to emulate that the kube-boat daemon is not running

				client := &Client{client: &http.Client{}}
				err := client.Readyz()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).Should(HavePrefix("kube-boat daemon or local Kubernetes API Server is not ready:"))
			})
		})

		When("the kube-boat daemon returns status code that indicates tha daemon or the API Server is unhealthy", func() {
			It("ends up with error", func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					"http://localhost"+readyz,
					httpmock.NewStringResponder(http.StatusInternalServerError, "ng"),
				)

				client := &Client{client: &http.Client{}}
				err := client.Readyz()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).Should(HavePrefix("kube-boat daemon or local Kubernetes API Server is not ready:"))
			})
		})
	})

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
					"http://localhost"+kubeconfig,
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

		When("the kube-boat daemon is not running", func() {
			It("ends up with error", func() {
				// do not register responder to emulate that the kube-boat daemon is not running

				client := &Client{client: &http.Client{}}
				kubeconfig, err := client.Kubeconfig()

				Expect(err).To(HaveOccurred())
				Expect(kubeconfig).To(BeNil())
			})
		})
	})

	Describe("func WebhookConfig()", func() {
		When("the kube-boat daemon returns correct webhook config response", func() {
			var response = `{
    "local-serving-host": "127.0.0.1",
    "local-serving-port": 53971,
    "local-serving-cert-dir": "/var/folders/vc/k8xt_qss50b2dcfpmrxwkbrc0000gn/T/envtest-serving-certs-2136998892"
}`
			It("should return a correctly unmarshalled result", func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					"http://localhost"+webhookConfig,
					httpmock.NewStringResponder(http.StatusOK, response),
				)

				client := &Client{client: &http.Client{}}
				webhookConfig, err := client.WebhookConfig()

				Expect(err).NotTo(HaveOccurred())
				Expect(webhookConfig.LocalServingHost).To(Equal("127.0.0.1"))
				Expect(webhookConfig.LocalServingPort).To(Equal(53971))
				Expect(webhookConfig.LocalServingCertDir).To(Equal("/var/folders/vc/k8xt_qss50b2dcfpmrxwkbrc0000gn/T/envtest-serving-certs-2136998892"))
			})
		})

		When("the kube-boat daemon is not running", func() {
			It("ends up with error", func() {
				// do not register responder to emulate that the kube-boat daemon is not running

				client := &Client{client: &http.Client{}}
				kubeconfig, err := client.WebhookConfig()

				Expect(err).To(HaveOccurred())
				Expect(kubeconfig).To(BeNil())
			})
		})
	})

	Describe("func StopDaemon()", func() {
		When("the kube-boat daemon returns accepting response", func() {
			var response = `{
    "message": "shutting down the API server..."
}`
			It("should return a correctly unmarshalled result", func() {
				httpmock.RegisterResponder(
					http.MethodDelete,
					"http://localhost"+base,
					httpmock.NewStringResponder(http.StatusAccepted, response),
				)

				client := &Client{client: &http.Client{}}
				msg, err := client.StopDaemon()

				Expect(err).NotTo(HaveOccurred())
				Expect(msg).To(Equal("shutting down the API server..."))
			})
		})

		When("the kube-boat daemon is not running", func() {
			It("ends up with error", func() {
				// do not register responder to emulate that the kube-boat daemon is not running

				client := &Client{client: &http.Client{}}
				msg, err := client.StopDaemon()

				Expect(err).To(HaveOccurred())
				Expect(msg).To(BeEmpty())
			})
		})
	})

})
