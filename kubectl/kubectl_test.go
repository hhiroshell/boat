package kubectl

import (
	"crypto/md5"
	"io"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubectl", func() {
	BeforeEach(func() {
		originalKubeconfig := os.Getenv("KUBECONFIG")
		DeferCleanup(func() {
			err := os.Setenv("KUBECONFIG", originalKubeconfig)
			Expect(err).NotTo(HaveOccurred())
		})

		err := os.Setenv("KUBECONFIG", "./kube-boat-test-kubeconfig")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(os.Getenv("KUBECONFIG"))
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("func SetContext()", func() {
		When("", func() {
			It("should return \"$HOME/.kube-boat/daemon.socket\"", func() {
				expect := `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://localhost:50000/
  name: kube-boat
contexts:
- context:
    cluster: kube-boat
    namespace: default
    user: kube-boat
  name: kube-boat
current-context: kube-boat
kind: Config
preferences: {}
users:
- name: kube-boat
  user:
    client-certificate-data: Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK
    client-key-data: Y2xpZW50LWtleS1kYXRhCg==
`

				err := SetContext(
					"https://localhost:50000/",
					"Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK",
					"Y2xpZW50LWtleS1kYXRhCg==",
					true,
				)
				Expect(err).NotTo(HaveOccurred())

				out, err := exec.Command("kubectl", "config", "view", "--minify", "--raw").Output()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(out)).Should(Equal(expect))
			})
		})

		When("", func() {
			It("should return \"$HOME/.kube-boat/daemon.socket\"", func() {
				expect := `apiVersion: v1
clusters:
- cluster:
    insecure-skip-tls-verify: true
    server: https://localhost:50000/
  name: kube-boat
contexts:
- context:
    cluster: kube-boat
    namespace: default
    user: kube-boat
  name: kube-boat
current-context: ""
kind: Config
preferences: {}
users:
- name: kube-boat
  user:
    client-certificate-data: Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK
    client-key-data: Y2xpZW50LWtleS1kYXRhCg==
`

				err := SetContext(
					"https://localhost:50000/",
					"Y2xpZW50LWNlcnRpZmljYXRlLWRhdGEK",
					"Y2xpZW50LWtleS1kYXRhCg==",
					false,
				)
				Expect(err).NotTo(HaveOccurred())

				rawConfig, err := os.Open(os.Getenv("KUBECONFIG"))
				Expect(err).NotTo(HaveOccurred())
				defer rawConfig.Close()

				hConfig := md5.New()
				_, err = io.Copy(hConfig, rawConfig)
				Expect(err).NotTo(HaveOccurred())

				hExpect := md5.New()
				_, err = io.WriteString(hExpect, expect)
				Expect(err).NotTo(HaveOccurred())

				Expect(hConfig.Sum(nil)).Should(Equal(hExpect.Sum(nil)))
			})
		})
	})

})
