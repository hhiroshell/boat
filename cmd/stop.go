package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var stopCmd = &cobra.Command{
	Use:          "stop",
	Short:        "Stop the running local Kubernetes API server",
	Long:         `Stop the running local Kubernetes API server`,
	SilenceUsage: true,

	RunE: stop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func stop(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create kube-boat daemon client: %w", err)
	}

	msg, err := client.StopDaemon()
	if err != nil {
		return fmt.Errorf("failed to stop kube-boat daemon or kube-apiserver: %w", err)
	}

	fmt.Println(msg)
	return nil
}
