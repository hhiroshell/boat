package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/kube-boat/pkg/daemon"
)

var configCmd = &cobra.Command{
	Use:           "config",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          config,
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

	fmt.Println(kubeconfig)

	return nil
}
