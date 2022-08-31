package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hhiroshell/kube-boat/daemon"
	"github.com/hhiroshell/kube-boat/kubectl"
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

	kubeconfig, err := client.Kubeconfig()
	if err != nil {
		return err
	}

	if err := kubectl.SetContext(kubeconfig.Server, kubeconfig.ClientCert, kubeconfig.ClientKey, true); err != nil {
		return err
	}

	return nil
}
