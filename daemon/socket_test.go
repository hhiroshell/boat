package daemon

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Socket", func() {
	var (
		sock *Socket
		home string
	)

	BeforeEach(func() {
		var err error
		sock, err = NewSocket()
		Expect(err).NotTo(HaveOccurred())

		home, err = os.UserHomeDir()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("func Path()", func() {
		It("should return \"$HOME/.kube-boat/daemon.socket\"", func() {
			path := sock.Path()
			Expect(path).Should(Equal(filepath.Join(home, ".kube-boat", "daemon.socket")))
		})
	})

	Describe("func Close()", Ordered, func() {
		//When("the kube-boat socket is not exist", func() {
		//	It("ends up with error", func() {
		//		err := sock.Close()
		//		Expect(err).To(HaveOccurred())
		//		Expect(err.Error()).Should(HavePrefix("failed remove unix socket file: "))
		//	})
		//})

		When("the kube-boat socket exists", func() {
			It("removes the socket", func() {
				file, err := os.Create(filepath.Join(home, ".kube-boat", "daemon.socket"))
				Expect(err).NotTo(HaveOccurred())

				err = sock.Close()
				Expect(err).NotTo(HaveOccurred())

				_, err = os.Stat(file.Name())
				Expect(errors.Is(err, fs.ErrNotExist)).Should(BeTrue())
			})
		})
	})
})
