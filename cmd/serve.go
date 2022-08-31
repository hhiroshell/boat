package cmd

import (
	"context"
	"golang.org/x/sys/unix"
	"os/signal"

	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/hhiroshell/kube-boat/daemon"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start local Kubernetes API server.",
	Long: `Start local Kubernetes API server and also resides as a daemon
that receives requests from other kube-boat commands.
You should not execute this command directly. It is intended to
be called via the kube-boat start command.`,
	Hidden: true,
	RunE:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(_ *cobra.Command, _ []string) error {
	sock, err := daemon.NewSocket()
	if err != nil {
		return err
	}

	// TODO: add flags to specify envtest options
	testEnv := &envtest.Environment{}

	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	if err := daemon.NewDaemon(sock, testEnv).Run(ctx, cancel); err != nil {
		return err
	}

	return nil
}
