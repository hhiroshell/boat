package daemon

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type Daemon struct {
	ctx context.Context

	sock    *Socket
	testEnv *envtest.Environment
}

func NewDaemon(sock *Socket, testEnv *envtest.Environment) *Daemon {
	return &Daemon{
		sock:    sock,
		testEnv: testEnv,
	}
}

func (d *Daemon) Run(ctx context.Context, cancel context.CancelFunc) error {
	config, err := d.testEnv.Start()
	if err != nil {
		return err
	}

	engine := gin.Default()
	engine.GET("/kube-boat", func(c *gin.Context) {
		c.JSON(http.StatusOK, Kubeconfig{
			Server:     config.Host,
			ClientCert: base64.StdEncoding.EncodeToString(config.CertData),
			ClientKey:  base64.StdEncoding.EncodeToString(config.KeyData),
		})
	})
	engine.DELETE("/kube-boat", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "shutting down the server...",
		})

		cancel()
	})

	go func() {
		if err := engine.RunUnix(d.sock.Path()); err != nil {
			log.Println("failed to start the daemon: %w", err)

			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the daemon...")
	if err := d.clean(); err != nil {
		return err
	}

	return nil
}

func (d *Daemon) clean() error {
	var errors error

	if err := d.testEnv.Stop(); err != nil {
		multierr.Append(errors, err)
	}
	if err := d.sock.Close(); err != nil {
		multierr.Append(errors, err)
	}

	return errors
}
