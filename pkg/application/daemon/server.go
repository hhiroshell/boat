package daemon

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/hhiroshell/kube-boat/pkg/common"
	"github.com/hhiroshell/kube-boat/pkg/infrastructure/socket"
)

type Daemon struct {
	ctx context.Context

	sock    *socket.Socket
	testEnv *envtest.Environment

	config *rest.Config
}

func NewDaemon(sock *socket.Socket, testEnv *envtest.Environment) *Daemon {
	return &Daemon{
		sock:    sock,
		testEnv: testEnv,
	}
}

func (d *Daemon) Run(ctx context.Context) error {
	config, err := d.testEnv.Start()
	if err != nil {
		return err
	}
	d.config = config

	engine := gin.Default()
	engine.GET("/kube-boat", d.kubeConfig)
	engine.DELETE("/kube-boat", d.shutdown)

	go func() {
		if err := engine.RunUnix(d.sock.Path()); err != nil {
			fmt.Println(err)
			d.clean()

			log.Fatalf("")
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the server...")
	d.clean()

	return nil
}

func (d *Daemon) kubeConfig(c *gin.Context) {
	c.JSON(http.StatusOK, common.KubeConfig{
		Server:     d.config.Host,
		ClientCert: base64.StdEncoding.EncodeToString(d.config.CertData),
		ClientKey:  base64.StdEncoding.EncodeToString(d.config.KeyData),
	})
}

func (d *Daemon) shutdown(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"message": "shutting down the server...",
	})

	d.ctx.Done()
}

func (d *Daemon) clean() {
	if err := d.testEnv.Stop(); err != nil {
		fmt.Println(err)
	}
	if err := d.sock.Close(); err != nil {
		fmt.Println(err)
	}
}
