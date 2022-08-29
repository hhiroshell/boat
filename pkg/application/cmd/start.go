package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:           "start",
	Short:         "",
	Long:          ``,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          start,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func start(_ *cobra.Command, _ []string) error {
	cmd := exec.Command("./boat", "serve")

	if err := cmd.Start(); err != nil {
		return err
	}
	fmt.Println("starting kube api server...")

	// TODO: wait for the daemon become running
	return nil
}
