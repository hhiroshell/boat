package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"

	"github.com/hhiroshell/boat/daemon"
)

var webhookConfigCmd = &cobra.Command{
	Use:   "webhook-config",
	Short: "Show API Server started by kube-boat.",
	Long:  `Update kubectl context to use API Server started by kube-boat.`,
	RunE:  webhookConfig,
}

func init() {
	rootCmd.AddCommand(webhookConfigCmd)
}

func webhookConfig(_ *cobra.Command, _ []string) error {
	client, err := daemon.NewClient()
	if err != nil {
		return err
	}

	webhookConfig, err := client.WebhookConfig()
	if err != nil {
		return err
	}

	fmt.Println("local serving host: " + webhookConfig.LocalServingHost)
	fmt.Println("local serving port: " + strconv.Itoa(webhookConfig.LocalServingPort))
	fmt.Println("local serving cert dir: " + webhookConfig.LocalServingCertDir)

	return nil
}
