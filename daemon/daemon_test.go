package daemon

import (
	"context"
	"errors"
	"io/fs"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var _ = Describe("Daemon", func() {
	Describe("func Run()", Ordered, func() {

		BeforeAll(func() {
			originalKubebuilderAssets := os.Getenv("KUBEBUILDER_ASSETS")
			DeferCleanup(func() {
				err := os.Setenv("KUBEBUILDER_ASSETS", originalKubebuilderAssets)
				Expect(err).NotTo(HaveOccurred())
			})

			err := os.Setenv("KUBEBUILDER_ASSETS", "/Users/hirhayak/Library/Application Support/io.kubebuilder.envtest/k8s/1.24.2-darwin-arm64")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterAll(func() {
		})

		When("kube-boat daemon is started with no optional arguments", Ordered, func() {
			var (
				sock    *Socket
				daemon  *Daemon
				testenv *envtest.Environment
			)

			BeforeAll(func() {
				sock = &Socket{path: "./kube-boat-test-daemon.sock"}
				testenv = &envtest.Environment{}
				daemon = NewDaemon(sock, testenv)
			})

			It("starts and stops without any problems", func() {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					err := daemon.Run(ctx, cancel)
					Expect(err).NotTo(HaveOccurred())
				}()

				By("checking Kubernetes API Server is available")

				Eventually(func() bool {
					return testenv.Config != nil
				}, 10, 1).Should(BeTrue())

				Eventually(func() error {
					k8s, err := client.New(testenv.Config, client.Options{Scheme: scheme.Scheme})
					if err != nil {
						return err
					}

					err = k8s.Get(context.Background(), client.ObjectKey{Name: "default"}, &corev1.Namespace{})
					return err
				}, 20, 1).Should(Succeed())

				By("checking kube-boat daemon endpoints are available")

				By("checking all resources are deleted")

				// TODO: shutting down with the daemon endpoint
				cancel()

				// TODO: check the API Server is not running
				Eventually(func() bool {
					_, err := os.Stat(sock.Path())
					return errors.Is(err, fs.ErrNotExist)
				}, 20, 1).Should(BeTrue())
			})
		})
	})

})
