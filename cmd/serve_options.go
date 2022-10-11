package cmd

import "github.com/spf13/cobra"

var (
	crdPaths    []string
	crdPathFlag = "crd-path"

	webhookConfigurationPaths    []string
	webhookConfigurationPathFlag = "webhook-config-path"
)

func setServeFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(
		&crdPaths,
		crdPathFlag,
		nil,
		"paths to the directory or file containing CRDs",
	)
	cmd.Flags().StringArrayVar(
		&webhookConfigurationPaths,
		webhookConfigurationPathFlag,
		nil,
		"paths to the directory or file containing webhook configurations",
	)
}
