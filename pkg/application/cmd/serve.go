package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"os/signal"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/hhiroshell/kube-boat/pkg/application/daemon"
	"github.com/hhiroshell/kube-boat/pkg/infrastructure/socket"
)

var serveCmd = &cobra.Command{
	Use:           "serve",
	Short:         "Start a standalone, local Kubernetes api server.",
	Long:          `Start a standalone, local Kubernetes api server.`,
	Hidden:        true,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(_ *cobra.Command, _ []string) error {
	sock, err := socket.NewSocket()
	if err != nil {
		return err
	}

	testEnv := &envtest.Environment{}

	ctx, _ := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	daemon.NewDaemon(sock, testEnv).Run(ctx)

	return nil
}

//func serve(_ *cobra.Command, _ []string) error {
//	sock, err := socket.NewSocket()
//	if err != nil {
//		return err
//	}
//
//	testEnv := &envtest.Environment{}
//	config, err := testEnv.Start()
//	if err != nil {
//		return err
//	}
//
//	// gin.New()
//	engine := gin.Default()
//
//	engine.GET("/testenv", func(c *gin.Context) {
//		c.JSON(http.StatusOK, common.KubeConfig{
//			Server:     config.Host,
//			ClientCert: base64.StdEncoding.EncodeToString(config.CertData),
//			ClientKey:  base64.StdEncoding.EncodeToString(config.KeyData),
//		})
//	})
//
//	quit := make(chan os.Signal)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//
//	engine.DELETE("/testenv", func(c *gin.Context) {
//		c.JSON(http.StatusAccepted, gin.H{
//			"message": "shutting down the server...",
//		})
//
//		quit <- syscall.SIGTERM
//	})
//
//	go func() {
//		if err := engine.RunUnix(sock.Path()); err != nil {
//			fmt.Println(err)
//		}
//	}()
//
//	<-quit
//	log.Println("shutting down the server...")
//
//	if err := testEnv.Stop(); err != nil {
//		return err
//	}
//	if err := sock.Close(); err != nil {
//		return err
//	}
//
//	return nil
//}
