package daemon

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.uber.org/multierr"
	"log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type Daemon struct {
	sock    *Socket
	testEnv *envtest.Environment
}

func NewDaemon(sock *Socket, testEnv *envtest.Environment) *Daemon {
	return &Daemon{
		sock:    sock,
		testEnv: testEnv,
	}
}

const (
	base          = "/kube-boat"
	readyz        = base + "/readyz"
	kubeconfig    = base + "/kubeconfig"
	webhookConfig = base + "/webhookconfig"
)

func (d *Daemon) Run(ctx context.Context, cancel context.CancelFunc) error {
	config, err := d.testEnv.Start()
	if err != nil {
		return err
	}

	cert := base64.StdEncoding.EncodeToString(config.CertData)
	key := base64.StdEncoding.EncodeToString(config.KeyData)

	engine := gin.Default()
	engine.GET(readyz, func(c *gin.Context) {
		client := http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
		res, err := client.Get(config.Host + "readyz")
		if err != nil {
			c.JSON(http.StatusInternalServerError, "ng")
			return
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			c.JSON(http.StatusOK, "ok")
		} else {
			c.JSON(http.StatusInternalServerError, "ng")
		}
	})
	engine.GET(kubeconfig, func(c *gin.Context) {
		c.JSON(http.StatusOK, Kubeconfig{
			Server:     config.Host,
			ClientCert: cert,
			ClientKey:  key,
		})
	})
	engine.GET(webhookConfig, func(c *gin.Context) {
		c.JSON(http.StatusOK, WebhookConfig{
			LocalServingHost:    d.testEnv.WebhookInstallOptions.LocalServingHost,
			LocalServingPort:    d.testEnv.WebhookInstallOptions.LocalServingPort,
			LocalServingCertDir: d.testEnv.WebhookInstallOptions.LocalServingCertDir,
		})
	})
	engine.DELETE(base, func(c *gin.Context) {
		c.JSON(http.StatusAccepted, Message{
			Message: "shutting down the API server...",
		})

		cancel()
	})

	go func() {
		if err := engine.RunUnix(d.sock.Path()); err != nil {
			log.Printf("failed to start the daemon: %s", err.Error())

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
