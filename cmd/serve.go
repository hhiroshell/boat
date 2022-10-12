package cmd

import (
	"context"
	"fmt"
	"os/signal"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/hhiroshell/boat/daemon"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start local Kubernetes API server.",
	Long: `Start local Kubernetes API server and also resides as a daemon
that receives requests from other kube-boat commands.
You should not execute this command directly. It is intended to
be called via the kube-boat start command.`,
	Hidden:       true,
	SilenceUsage: true,

	RunE: serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	setServeFlags(serveCmd)
}

func serve(_ *cobra.Command, _ []string) error {
	sock, err := daemon.NewSocket()
	if err != nil {
		return fmt.Errorf("failed to create unix socket to communicate with kube-boat daemon: %w", err)
	}

	testEnv := &envtest.Environment{
		CRDInstallOptions: envtest.CRDInstallOptions{
			Paths: crdPaths,
		},
		WebhookInstallOptions: envtest.WebhookInstallOptions{
			Paths: webhookConfigurationPaths,
		},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), unix.SIGINT, unix.SIGTERM, unix.SIGHUP)

	if err := daemon.NewDaemon(sock, testEnv).Run(ctx, cancel); err != nil {
		return fmt.Errorf("failed to start kube-boat daemon or kube-apiserver: %w", err)
	}

	return nil
}
