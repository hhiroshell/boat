package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var configCmd = &cobra.Command{
	Use:   "kubeconfig",
	Short: "Update kubectl context to use API Server started by kube-boat.",
	Long:  `Update kubectl context to use API Server started by kube-boat.`,
	RunE:  config,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func config(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return err
	}

	if err := setKubectlContext(client); err != nil {
		return fmt.Errorf("failed to update kubeconfig: %w", err)
	}

	return nil
}
