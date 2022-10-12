package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hhiroshell/boat/daemon"
	"github.com/hhiroshell/boat/kubectl"
)

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

func setKubectlContext(client *daemon.Client) error {
	kubeconfig, err := client.Kubeconfig()
	if err != nil {
		return err
	}

	if err := kubectl.SetContext(kubeconfig.Server, kubeconfig.ClientCert, kubeconfig.ClientKey, true); err != nil {
		return err
	}

	return nil
}
