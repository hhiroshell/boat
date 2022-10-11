package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
)

var webhookConfigCmd = &cobra.Command{
	Use:   "webhook-config",
	Short: "Displays properties about admission / validation webhook server targeted by the kube-apiserver.",
	Long: `Displays properties about admission / validation webhook server targeted by the kube-apiserver.
Your local webhook server should be run with these properties.
`,
	RunE: webhookConfig,
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
