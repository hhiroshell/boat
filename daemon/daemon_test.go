package daemon

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	k8sClient "sigs.k8s.io/controller-runtime/pkg/client"
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

			o, err := exec.Command(
				"setup-envtest",
				"use",
				"-i",
				"-p=env",
				"1.24.x",
			).Output()
			Expect(err).NotTo(HaveOccurred())

			kubebuilderAssets := strings.TrimSuffix(strings.TrimPrefix(string(o), "export KUBEBUILDER_ASSETS='"), "'\n")
			err = os.Setenv("KUBEBUILDER_ASSETS", kubebuilderAssets)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("kube-boat daemon is started with no optional arguments", Ordered, func() {
			var (
				sock    *Socket
				testenv *envtest.Environment

				daemonClient *http.Client
			)

			BeforeAll(func() {
				sock = &Socket{path: "./kube-boat-test-daemon.sock"}
				testenv = &envtest.Environment{}

				daemonClient = &http.Client{
					Transport: &http.Transport{
						DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
							return net.Dial("unix", sock.Path())
						},
					},
				}
			})

			It("starts without any errors", func() {
				daemon := NewDaemon(sock, testenv)
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					err := daemon.Run(ctx, cancel)
					Expect(err).NotTo(HaveOccurred())
				}()
			})

			It("makes local Kubernetes API Server available", func() {
				Eventually(func() bool {
					return testenv.Config != nil
				}, 15, 1).Should(BeTrue())

				Eventually(func() error {
					k8s, err := k8sClient.New(testenv.Config, k8sClient.Options{Scheme: scheme.Scheme})
					if err != nil {
						return err
					}

					err = k8s.Get(context.Background(), k8sClient.ObjectKey{Name: "default"}, &corev1.Namespace{})
					return err
				}, 15, 1).Should(Succeed())
			})

			It("makes kube-boat's \"readyz\" endpoint available", func() {
				var res *http.Response
				Eventually(func() error {
					var err error
					res, err = daemonClient.Get("http://localhost" + readyz)
					return err
				}, 15, 1).Should(Succeed())
				defer res.Body.Close()
			})

			It("makes kube-boat's \"kubeconfig\" endpoint available", func() {
				var res *http.Response
				Eventually(func() error {
					var err error
					res, err = daemonClient.Get("http://localhost" + kubeconfig)
					return err
				}, 15, 1).Should(Succeed())
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				Expect(err).NotTo(HaveOccurred())

				kubeconfig := &Kubeconfig{}
				err = json.Unmarshal(body, kubeconfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(kubeconfig.Server).Should(Equal(testenv.Config.Host))
				Expect(kubeconfig.ClientCert).Should(Equal(base64.StdEncoding.EncodeToString(testenv.Config.CertData)))
				Expect(kubeconfig.ClientKey).Should(Equal(base64.StdEncoding.EncodeToString(testenv.Config.KeyData)))
			})

			It("makes kube-boat's shutdown endpoint available", func() {
				req, err := http.NewRequest(http.MethodDelete, "http://localhost"+base, nil)
				Expect(err).NotTo(HaveOccurred())

				res, err := daemonClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				Expect(err).NotTo(HaveOccurred())

				msg := &Message{}
				err = json.Unmarshal(body, msg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("kube-boat daemon and Kubernetes API Server are shutdown", func() {
				// TODO: check the API Server is not running
				Eventually(func() bool {
					_, err := os.Stat(sock.Path())
					return errors.Is(err, fs.ErrNotExist)
				}, 10, 1).Should(BeTrue())
			})

		})
	})

})
