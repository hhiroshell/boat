package cmd

import (
	"context"
	"golang.org/x/sys/unix"
	"os/signal"

	"github.com/spf13/cobra"
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

	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	daemon.NewDaemon(sock, testEnv).Run(ctx, cancel)

	return nil
}
