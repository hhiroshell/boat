package cmd

import "github.com/spf13/cobra"

var (
	crdPaths    []string
	crdPathFlag = "crd-path"
)

func setServeFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVar(&crdPaths, crdPathFlag, nil, "paths to the directory or file containing CRDs")
}
